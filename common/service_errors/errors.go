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
	Field    string `json:"-"`
	module   string `json:"-"`
	function string `json:"-"`
	message  string `json:"-"`
	Message  string `json:"-"`
}

func New(code int, httpCode int, module string, function string, field string, message string, args ...any) *Error {
	e := new(Error)
	e.AppName = AppName
	e.Code = code
	e.HttpCode = httpCode
	e.module = module
	e.function = function
	e.message = fmt.Sprintf(message, args...)

	return e
}

func (e *Error) Error() string {
	e.Message = fmt.Sprintf("[%s::%s::%s] -> %s", e.AppName, e.module, e.function, e.message)
	return e.Message
}

func (e *Error) LogError() *Error {
	slog.Error(e.Error())
	return e
}
