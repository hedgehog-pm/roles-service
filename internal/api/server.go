package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v5"
	middleware "github.com/oapi-codegen/echo-v5-middleware"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
)

// Authenticate must verify the request credentials before returning a positive user ID.
type Authenticate func(context.Context, *http.Request) (int, error)

type actorKey struct{}

func actorID(ctx context.Context) (int, error) {
	id, ok := ctx.Value(actorKey{}).(int)
	if !ok || id <= 0 {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	return id, nil
}

func NewServer(handler generated.StrictServerInterface, authenticate Authenticate) (*echo.Echo, error) {
	spec, err := generated.GetSpec()
	if err != nil {
		return nil, err
	}

	server := echo.New()
	server.HTTPErrorHandler = HTTPErrorHandler
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if authenticate == nil {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "authentication is not configured")
			}

			id, err := authenticate(c.Request().Context(), c.Request())
			if err != nil {
				return err
			}

			if id <= 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), actorKey{}, id)))

			return next(c)
		}
	})
	generated.RegisterHandlers(server.Group("", middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{Options: openapi3filter.Options{MultiError: true}, MultiErrorHandler: validationError})), generated.NewStrictHandler(handler, nil))

	return server, nil
}

func HTTPErrorHandler(c *echo.Context, err error) {
	if r, e := echo.UnwrapResponse(c.Response()); e == nil && r.Committed {
		return
	}

	status, body := errorResponse(err)
	if status >= 500 {
		slog.ErrorContext(c.Request().Context(), "request failed", "error", err)
	}

	if writeErr := c.JSON(status, body); writeErr != nil {
		slog.ErrorContext(c.Request().Context(), "write error response", "error", writeErr)
	}
}

func errorResponse(err error) (int, generated.Error) {
	status := http.StatusInternalServerError
	body := generated.Error{Code: "internal_error", Message: "internal error"}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
	default:
		if problem, ok := errors.AsType[*domain.Error](err); ok {
			return http.StatusBadRequest, generated.Error{Code: string(problem.Code), Message: problem.Message, Details: &[]generated.ErrorDetail{{Field: problem.Field, Code: string(problem.Code), Message: problem.Message}}}
		} else if code := echo.StatusCode(err); code >= 400 && code <= 599 {
			status = code
		}
	}

	switch status {
	case http.StatusBadRequest:
		body = generated.Error{Code: "invalid_input", Message: "invalid input"}
	case http.StatusUnauthorized:
		body = generated.Error{Code: "unauthorized", Message: "unauthorized"}
	case http.StatusNotFound:
		body = generated.Error{Code: "not_found", Message: "not found"}
	case http.StatusConflict:
		body = generated.Error{Code: "version_conflict", Message: "conflict"}
	case http.StatusServiceUnavailable:
		body = generated.Error{Code: "unavailable", Message: "service unavailable"}
	default:
		if status < 500 {
			body = generated.Error{Code: "invalid_input", Message: http.StatusText(status)}
		}
	}

	return status, body
}
