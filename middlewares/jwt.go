package middlewares

import (
	"Backend-Recything/helper"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// jwtCustomClaims struct untuk klaim JWT
type jwtCustomClaims struct {
	Name   string `json:"name"`
	UserID uint   `json:"userID"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// CustomValidator struct untuk custom validator Echo
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate implementasi untuk validasi custom
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// AuthMiddleware middleware untuk memvalidasi token JWT
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token dari header Authorization
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, helper.APIResponse("Missing token", http.StatusUnauthorized, "error", nil))
		}

		// Memastikan token memiliki format "Bearer <token>"
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, helper.APIResponse("Invalid token format", http.StatusUnauthorized, "error", nil))
		}

		// Parsing token
		tokenString = parts[1]

		// Parsing klaim JWT
		claims := &jwtCustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		// Validasi token
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, helper.APIResponse("Invalid token", http.StatusUnauthorized, "error", nil))
		}

		// Menyimpan klaim di context untuk diakses di handler berikutnya
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		return next(c)
	}
}

// RoleMiddleware middleware untuk memvalidasi peran pengguna
func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Mendapatkan role pengguna dari context
			userRole, ok := c.Get("userRole").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, helper.APIResponse("Missing or invalid user role", http.StatusUnauthorized, "error", nil))
			}

			// Memeriksa apakah role pengguna sesuai dengan role yang diizinkan
			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			// Jika role tidak cocok, berikan respons "Access Denied"
			return c.JSON(http.StatusForbidden, helper.APIResponse("Access denied", http.StatusForbidden, "error", nil))
		}
	}
}
