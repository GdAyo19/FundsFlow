package middleware

import (
	"net/http"
	"strings"

	"github.com/GdAyo19/FundsFlow/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	// return a middleware function that checks for a valid JWT token in the Authorization header
	return func(c *gin.Context) {

		// get the Authorization header from the request
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}
		// check if the Authorization header starts with "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// check if the tokenString is empty after trimming the prefix
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization format",
			})
			c.Abort()
			return
		}

		// validate the token using the ValidateToken function from the utils package
		claims, err := utils.ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// set the user ID from the claims in the context for further use in the request
		c.Set("userID", claims.UserID)

		// call the next handler in the chain
		c.Next()
	}
}
