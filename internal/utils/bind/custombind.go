package bind

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

func Bind(c echo.Context, v any) error {
	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&v); err != nil {
		return err
	}

	return nil
}
