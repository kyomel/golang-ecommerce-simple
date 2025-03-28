package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Todo: ambil header authorization
		key := os.Getenv("ADMIN_SECRET")
		// Validasi header sesuai dengan kata sandi admin
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if auth != key {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Lanjutkan request ke handler
		c.Next()
	}

}
