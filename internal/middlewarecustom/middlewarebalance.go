package middlewarecustom

import (
	"api/internal/servises/jwtservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func MiddlewareBalance(jwtService jwtservice.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error":"There is no cookie with the required header"})
			}

			token := cookie.Value

			err = jwtService.CheckJWT(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error":"Authorization error"})
			}

			id, err := jwtService.GetIdFromJWT(token)
			if err != nil {
				return err
			}

			c.Set("id", id)

			return next(c)
		}
	}
}