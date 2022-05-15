package serialize

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type ErrorResponse struct {
	Errors *map[string]string `json:"errors,omitempty"`
	Response
}

type DataResponse[T interface{}] struct {
	Data T `json:"result"`
	Response
}
