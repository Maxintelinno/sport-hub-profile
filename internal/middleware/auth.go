package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid authorization header format",
				})
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid or expired token: " + err.Error(),
				})
			}

			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid token",
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid token claims",
				})
			}

			// Extract user ID from claims (flexible handling for user_id or userid)
			var userID string
			if id, exists := claims["user_id"]; exists {
				userID = fmt.Sprintf("%v", id)
			} else if id, exists := claims["userid"]; exists {
				userID = fmt.Sprintf("%v", id)
			}

			if userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid user ID in token (missing user_id or userid claim)",
				})
			}

			c.Set("user_id", userID)
			return next(c)
		}
	}
}
