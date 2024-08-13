package middleware

import "github.com/gin-gonic/gin"

// NoCache is a middleware to set headers to prevent caching of HTTP responses
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set the Cache-Control header to prevent caching
		// "no-store" ensures that no part of the response is cached
		// "no-cache" indicates that the response must be validated with the origin server before using it
		// "must-revalidate" forces the cache to revalidate the response with the origin server
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")

		// Set the Pragma header to "no-cache" for compatibility with HTTP/1.0 caches
		// This header is used by older HTTP/1.0 clients to indicate that the response should not be cached
		c.Header("Pragma", "no-cache")

		// Set the Expires header to "0" to indicate that the response is already expired
		// This prevents the client from storing and reusing the response
		c.Header("Expires", "0")

		// Call the next middleware/handler in the chain
		c.Next()
	}
}
