package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

var ErrUnauthenticated = errors.New("unauthenticated")

type VerifyToken func(ctx context.Context, token string) (int, error)

func Authentication(verify VerifyToken) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if verify == nil {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "authentication is not configured")
			}

			headers := c.Request().Header.Values("Authorization")
			if len(headers) != 1 {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			parts := strings.Fields(headers[0])
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			id, err := verify(c.Request().Context(), parts[1])
			switch {
			case errors.Is(err, ErrUnauthenticated):
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			case err != nil:
				return echo.NewHTTPError(http.StatusServiceUnavailable, "authentication unavailable").Wrap(err)
			case id <= 0 || id > 2147483647:
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			ctx := context.WithValue(c.Request().Context(), actorKey{}, id)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
