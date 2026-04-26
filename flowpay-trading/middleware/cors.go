package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware sets up the CORS rules
func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"} // Frontend Angular default
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	return cors.New(config)
}
