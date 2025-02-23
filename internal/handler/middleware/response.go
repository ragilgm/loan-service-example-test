package middleware

import (
	"github.com/labstack/echo"
	"github.com/test/loan-service/internal/dto"
	"strconv"
)

func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Menangani error yang tidak terduga di dalam handler
		err := next(c)
		if err != nil {
			return dto.SendError(c, 500, "Internal server error")
		}
		return nil
	}
}

func SuccessHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Panggil handler sebenarnya
		err := next(c)
		if err != nil {
			var msg string
			msg, err2 := GetErrorMessage(c, err.Error())
			if err2 != nil {
				return dto.SendError(c, 500, "Internal error")
			}

			var code int
			code, err = strconv.Atoi(err.Error())
			// Tangani error yang terjadi di handler
			return dto.SendError(c, code, msg)
		}
		// Menangani response sukses di sini jika perlu
		return nil
	}
}
