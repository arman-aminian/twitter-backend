package middleware

import (
	"fmt"
	"github.com/arman-aminian/twitter-backend/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type (
	GlobalConfig struct {
		Skipper    Skipper
		SigningKey interface{}
	}
	GlobalSkipper   func(c echo.Context) bool
	GlobalExtractor func(echo.Context) (string, error)
)

func Global(key interface{}) echo.MiddlewareFunc {
	c := GlobalConfig{}
	c.SigningKey = key
	return GlobalWithConfig(c)
}

func GlobalWithConfig(config GlobalConfig) echo.MiddlewareFunc {
	extractor := globalFromHeader("Authorization", "Token")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth, _ := extractor(c)

			if auth == "" {
				return next(c)
			}
			token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return config.SigningKey, nil
			})
			if err != nil {
				return c.JSON(http.StatusForbidden, utils.NewError(ErrJWTInvalid))
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				username := claims["id"]
				c.Set("username", username)
				return next(c)
			}
			return next(c)
		}
	}
}

// globalFromHeader returns a `jwtExtractor` that extracts token from the request header.
func globalFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}
