package gpt

// APIError provides error information returned by the OpenAI API.
type APIError struct {
	Code       *string `json:"code,omitempty"`
	Message    string  `json:"message"`
	Param      *string `json:"param,omitempty"`
	Type       string  `json:"type"`
	StatusCode int     `json:"-"`
}
type ErrorResponse struct {
	Error *APIError `json:"error,omitempty"`
}
