package dto

import (
	"github.com/labstack/echo"
	"net/http"
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
