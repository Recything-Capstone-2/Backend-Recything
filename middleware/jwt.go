package middleware

import (
	"Backend-Recything/helper"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type jwtCustomClaims struct {
	Name   string `json:"name"`
	UserID uint   `json:"userID"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, helper.APIResponse("Missing token", http.StatusUnauthorized, "error", nil))
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, helper.APIResponse("Invalid token", http.StatusUnauthorized, "error", nil))
		}

		// Set user ID dari token ke context
		c.Set("userID", (*claims)["userID"])

		return next(c)
	}
}
func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("user").(*jwtCustomClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, helper.APIResponse("Invalid token claims", http.StatusUnauthorized, "error", nil))
			}

			userRole := claims.Role
			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, helper.APIResponse("Forbidden: Access denied", http.StatusForbidden, "error", nil))
		}
	}
}
