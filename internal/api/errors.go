package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/hedgehog-pm/roles-service/internal/domain"
	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
)

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
	case http.StatusServiceUnavailable:
		body = generated.Error{Code: "unavailable", Message: "service unavailable"}
	case http.StatusNotFound:
		body = generated.Error{Code: "not_found", Message: "not found"}
	case http.StatusConflict:
		body = generated.Error{Code: "version_conflict", Message: "conflict"}
	default:
		if status < 500 {
			body = generated.Error{Code: "invalid_input", Message: http.StatusText(status)}
		}
	}

	return status, body
}
