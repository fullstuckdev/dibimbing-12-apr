package middleware

import (
	"webroutes/utils"
	"strings"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// BAKAL BACA HEADER => AUTHORIZATION
		authHeader := c.GetHeader("Authorization")


		// CHECK AUTH
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization is required"})
			c.Abort()
			return
		}

		// MEMBACA AUTH HEADER
		parts := strings.Split(authHeader, " ")

		// MEMBACA INDEX PERTAMA
		userId, err := utils.ValidateToken(parts[1])

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		// SETTING key userId dan value userID
		c.Set("userId", userId)
		c.Next()
	}
}