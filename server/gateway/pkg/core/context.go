package core

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

var UserCtxKey = &contextKey{"user"}
var HttpWriterKey = &contextKey{"httpWriter"}

type contextKey struct {
	name string
}

func UserMiddleware(db *sqlx.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("Authorization")
			if err != nil {
				return next(c)
			}

			tokenData, err := DecodeToken(cookie.Value)
			if err != nil {
				return next(c)
			}

			if claims, ok := tokenData.Claims.(jwt.MapClaims); ok && tokenData.Valid {
				sub, ok := claims["sub"].(string)
				if !ok {
					return next(c)
				}

				usr, err := GetUser(db, sub)
				if err != nil {
					return next(c)
				}

				ctx := context.WithValue(c.Request().Context(), UserCtxKey, usr)
				c.SetRequest(c.Request().WithContext(ctx))
			}

			return next(c)
		}
	}
}

// // CookieMiddleWare
func CookieMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), HttpWriterKey, c.Response().Writer)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func GetUserFromContext(ctx context.Context) *User {
	usr, ok := ctx.Value(UserCtxKey).(*User)

	if !ok {
		return nil
	}

	return usr
}
