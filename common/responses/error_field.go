package responses

type ErrField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrFields struct {
	Fields  []string `json:"fields"`
	Message string   `json:"message"`
}

type ResponseErrorField struct {
	Error ErrField `json:"error"`
}

type ErrMessage struct {
	Message string `json:"message"`
}

type ResponseErrorMessage struct {
	Error ErrMessage `json:"error"`
}
