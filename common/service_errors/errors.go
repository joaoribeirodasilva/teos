package service_errors

var AppName string

type Error struct {
	Code     int    `json:"code"`
	HttpCode int    `json:"-"`
	Location string `json:"location"`
	Fields   string `json:"fields,omitempty"`
	Message  string `json:"Message"`
}

func New(code int, httpCode int, location string, fields string, message string) *Error {
	e := new(Error)
	e.Code = code
	e.Location = location
	e.Fields = fields
	e.Message = message
	e.HttpCode = httpCode

	return e
}

func (e *Error) Error() string {
	return e.Message
}
