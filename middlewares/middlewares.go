package middlewares

import (
	"automotive/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var API_SECRET = utils.GetEnv("API_SECRET", "rahasiasekali")

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := utils.APIResponse("Unauthorized", http.StatusUnauthorized, "error", "Authorization header is missing")
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")

		if len(tokenString) != 2 {
			response := utils.APIResponse("Unauthorized", http.StatusUnauthorized, "error", "Authorization header is invalid")
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			// Make sure that the token's signing method is valid
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(API_SECRET), nil
		})

		if err != nil {
			response := utils.APIResponse("Unauthorized", http.StatusUnauthorized, "error", "Invalid Token")
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// You can access the claims here
			user_id := claims["user_id"]
			username := claims["username"]
			// Pass the user_id to the next handler
			c.Set("user_id", user_id)
			c.Set("username", username)
			c.Next()
		} else {
			response := utils.APIResponse("Unauthorized", http.StatusUnauthorized, "error", "Invalid Token")
			c.JSON(http.StatusUnauthorized, response)
			return
		}
	}
}
