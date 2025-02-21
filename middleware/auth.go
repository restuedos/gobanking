package middleware

import (
	"gobanking/config"
	"gobanking/model"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Remark: "Missing authorization header"})
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWT.Secret), nil
			})

			if err != nil || !token.Valid {
				cfg.Logger.Warn("invalid token", "error", err)
				return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Remark: "Invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				cfg.Logger.Warn("invalid token claims")
				return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Remark: "Invalid token"})
			}

			c.Set("user_id", claims["user_id"])
			c.Set("email", claims["email"])

			return next(c)
		}
	}
}
