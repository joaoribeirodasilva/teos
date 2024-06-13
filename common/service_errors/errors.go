package service_errors

import (
	"fmt"
	"log/slog"
)

var AppName string

type Error struct {
	AppName  string `json:"-"`
	Code     int    `json:"code"`
	HttpCode int    `json:"-"`
	Field    string `json:"field"`
	Module   string `json:"module"`
	Function string `json:"function"`
	message  string `json:"-"`
	Message  string `json:"Message"`
}

func New(code int, httpCode int, module string, function string, field string, message string, args ...any) *Error {
	e := new(Error)
	e.AppName = AppName
	e.Code = code
	e.HttpCode = httpCode
	e.Module = module
	e.Function = function
	e.Message = fmt.Sprintf(message, args...)

	return e
}

func (e *Error) Error() string {
	e.message = fmt.Sprintf("[%s::%s::%s] -> %s", e.AppName, e.Module, e.Function, e.message)
	return e.message
}

func (e *Error) LogError() *Error {
	slog.Error(e.Error())
	return e
}
