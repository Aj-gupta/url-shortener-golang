package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"urlshortner/config/db"
	dotenv "urlshortner/config/dotenv"
	middleware "urlshortner/middlewares"
	routes "urlshortner/routes"
	Logger "urlshortner/utils/logger"
	validator "urlshortner/utils/validator"

	"time"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	/**
	@description Setup Server
	*/
	dotenv.LoadConfig(".")
	router := SetupRouter()
	/**
	@description Run Server
	*/
	srv := &http.Server{
		Addr:    ":" + dotenv.Global.GoPORT,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func SetupRouter() *gin.Engine {
	/*
		@description	Init Logger
	*/
	Logger.Init()
	/**
	@description Setup Database Connection
	*/
	DB := db.ConnectDB()
	/**
	@description Setup Mode Application
	*/
	if dotenv.Global.GoEnv != "production" && dotenv.Global.GoEnv != "test" {
		gin.SetMode(gin.DebugMode)
	} else if dotenv.Global.GoEnv == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	/**
	@description Init Router
	*/
	router := gin.Default()
	/*
		@description	Init Validator
	*/
	validator.Init()
	/*
		@description add requestID middleware
	*/
	router.Use(requestid.New())
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	router.Use(ginzap.GinzapWithConfig(Logger.Log, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			// log request ID
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request-id", requestID))
			}

			// log trace and span ID
			if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
				fields = append(fields, zap.String("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()))
				fields = append(fields, zap.String("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()))
			}

			if gin.Mode() == gin.DebugMode {
				// log request body
				var body []byte
				var buf bytes.Buffer
				tee := io.TeeReader(c.Request.Body, &buf)
				body, _ = io.ReadAll(tee)
				c.Request.Body = io.NopCloser(&buf)
				fields = append(fields, zap.String("body", string(body)))
			}

			return fields
		}),
	}))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	router.Use(ginzap.RecoveryWithZap(Logger.Log, true))
	/**
	@description Setup Middleware
	*/
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	router.Use(helmet.Default())
	router.Use(gzip.Gzip(gzip.BestCompression))
	router.HandleMethodNotAllowed = true
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "method not allowed"})
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "url not found"})
	})
	// use timout middleware with time config.timeout seconds if config.timeout not
	// defined then upto 10 seconds
	if dotenv.Global.RequestTimeout == 0 {
		dotenv.Global.RequestTimeout = 10
	}
	router.Use(middleware.TimeoutMiddleware(time.Duration(dotenv.Global.RequestTimeout) * time.Second))
	/**
	@description Init All Route
	*/
	routes.InitShortUrlRoute(DB, router)
	routes.InitRedirectRoute(DB, router)

	return router
}
