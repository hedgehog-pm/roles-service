package api

import (
	"context"

	"github.com/labstack/echo/v5"
)

type actorKey struct{}

func actorID(ctx context.Context) int {
	id, _ := ctx.Value(actorKey{}).(int)
	return id
}

func DevelopmentActor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), actorKey{}, 1)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
