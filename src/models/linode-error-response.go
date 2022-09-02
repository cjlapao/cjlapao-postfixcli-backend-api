package models

type LinodeErrorResponse struct {
	Errors []LinodeError `json:"errors"`
}

type LinodeError struct {
	Reason string `json:"reason"`
	Field  string `json:"field"`
}
