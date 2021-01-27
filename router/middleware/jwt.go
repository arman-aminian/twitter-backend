package middleware

import (
	"errors"
	"fmt"
	"github.com/arman-aminian/twitter-backend/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type (
	JWTConfig struct {
		Skipper    Skipper
		SigningKey interface{}
	}
	Skipper      func(c echo.Context) bool
	jwtExtractor func(echo.Context) (string, error)
)

var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	ErrJWTInvalid = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt")
)

func JWT(key interface{}) echo.MiddlewareFunc {
	c := JWTConfig{}
	c.SigningKey = key
	return JWTWithConfig(c)
}

func JWTWithConfig(config JWTConfig) echo.MiddlewareFunc {
	extractor := jwtFromHeader("Authorization", "Token")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth, _ := extractor(c)
			if auth == "" {
				cookie, err := c.Cookie("Token")
				if err != nil {
					c.Set("username", "")
					return next(c)
				}
				auth = cookie.Value
			}
			if auth == "" {
				if config.Skipper != nil {
					if config.Skipper(c) {
						return next(c)
					}
				}
				return c.JSON(http.StatusUnauthorized, utils.NewError(errors.New("missing or malformed jwt")))
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
				//fmt.Println("jwt username")
				//fmt.Println(username)
				c.Set("username", username)
				return next(c)
			}
			return c.JSON(http.StatusForbidden, utils.NewError(ErrJWTInvalid))
		}
	}
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}
