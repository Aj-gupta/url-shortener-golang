package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set up a context with the custom timeout duration
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace the original context with the new context
		c.Request = c.Request.WithContext(ctx)

		// Wait for the request to finish
		done := make(chan struct{})
		go func() {
			c.Next()
			done <- struct{}{}
		}()

		// Wait for either the request to finish or the timeout to occur
		select {
		case <-done:
		case <-ctx.Done():
			// If the timeout occurs, set the response code to 408 (Request Timeout)
			// and return a custom error message
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"error": fmt.Sprintf("Request timed out after %v", timeout),
			})
		}
	}
}
