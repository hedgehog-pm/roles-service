package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v5"

	"github.com/hedgehog-pm/roles-service/internal/roles/domain"
)

func validationError(issues openapi3.MultiError) *echo.HTTPError {
	field := "request"
	code := domain.CodeInvalidInput
	message := "invalid request field"

	if request, ok := errors.AsType[*openapi3filter.RequestError](issues); ok && request.Parameter != nil {
		field = request.Parameter.Name
	}

	if schema, ok := errors.AsType[*openapi3.SchemaError](issues); ok {
		if path := schema.JSONPointer(); len(path) > 0 {
			field = strings.Join(path, ".")
		}

		if field == "permission_codes" && schema.SchemaField == "minItems" {
			code = domain.CodeEmptyPermissions
			message = "at least one permission is required"
		}

		if schema.SchemaField == "uniqueItems" {
			message = "duplicate values are not allowed"
		}
	}

	wrapped := echo.NewHTTPError(http.StatusBadRequest, message).Wrap(domain.FieldError(code, field, message))
	result, _ := errors.AsType[*echo.HTTPError](wrapped)

	return result
}
