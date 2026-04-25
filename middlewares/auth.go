package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// The bearer token sent in the request
		autHeader := c.GetHeader("Authorization")

		if autHeader == "" || !strings.HasPrefix(autHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error:": "Access is not authorized"})
			return
		}
		// Extract the token from the header
		tokenString := strings.TrimPrefix(autHeader, "Bearer ")

		// Parse the token to ensure it is valid
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}
			return []byte(os.Getenv("JWT_SECREt")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid or expired"})
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token cannot be read"})
			return
		}
		userID := int(claim["UserID"].(float64))

		c.Set("userID", userID)

		c.Next()
	}
}
