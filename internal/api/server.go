package api

import (
	"fmt"

	"github.com/labstack/echo/v5"
	middleware "github.com/oapi-codegen/echo-v5-middleware"

	generated "github.com/hedgehog-pm/roles-service/internal/generated/api"
)

func NewServer(handler generated.StrictServerInterface) (*echo.Echo, error) {
	spec, err := generated.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load roles OpenAPI contract: %w", err)
	}

	server := echo.New()
	server.HTTPErrorHandler = HTTPErrorHandler
	generated.RegisterHandlers(
		server.Group("", middleware.OapiRequestValidator(spec)),
		generated.NewStrictHandler(handler, nil),
	)

	return server, nil
}
