package server

import (
	"api/internal/config"
	"fmt"

	"github.com/labstack/echo/v4"
)

func StartServer(e *echo.Echo, cfg *config.Config) error {
	if err := e.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		return err
	}

	return nil
}