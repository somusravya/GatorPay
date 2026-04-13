package utils

import "github.com/gin-gonic/gin"

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": message,
	})
}

// SuccessResponse sends a standardized success JSON response
func SuccessResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"data": data,
	})
}
