package middleware

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type SuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error message"`
}

func SendSuccess(c echo.Context, data interface{}) error {
	response := SuccessResponse{
		Code: 0,
		Data: data,
	}
	return c.JSON(http.StatusOK, response)
}

func SendError(c echo.Context, code int, message string) error {
	response := ErrorResponse{
		Code:         code,
		ErrorMessage: message,
	}
	return c.JSON(http.StatusOK, response)
}

func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Menangani error yang tidak terduga di dalam handler
		err := next(c)
		if err != nil {
			return SendError(c, 500, "Internal server error")
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
				return SendError(c, 500, "Internal error")
			}

			var code int
			code, err = strconv.Atoi(err.Error())
			// Tangani error yang terjadi di handler
			return SendError(c, code, msg)
		}
		// Menangani response sukses di sini jika perlu
		return nil
	}
}
