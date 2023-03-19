package main

import (
	"net/http"
	"urlshortner/config/db"
	dotenv "urlshortner/config/dotenv"
	middleware "urlshortner/middlewares"
	routes "urlshortner/routes"
	Logger "urlshortner/utils/logger"
	validator "urlshortner/utils/validator"

	helmet "github.com/danielkov/gin-helmet"
	ginzap "github.com/gin-contrib/zap"

	"time"

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
	log.Fatal(router.Run(":" + dotenv.Global.GoPORT))
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
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	router.Use(ginzap.Ginzap(Logger.Log, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	router.Use(ginzap.RecoveryWithZap(Logger.Log, true))
	// use timout middleware with time config.timeout seconds if config.timeout not
	// defined then upto 10 seconds
	if dotenv.Global.RequestTimeout == 0 {
		dotenv.Global.RequestTimeout = 10
	}
	router.Use(middleware.TimeoutMiddleware(time.Duration(dotenv.Global.RequestTimeout) * time.Second))
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
	/**
	@description Init All Route
	*/
	routes.InitShortUrlRoute(DB, router)
	routes.InitRedirectRoute(DB, router)

	return router
}
