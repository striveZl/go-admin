package middleware

import "github.com/gin-gonic/gin"

func Empty() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
