package utils

import (
	"net/http"
)

type HTTPResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
}

func NewHttpResponse(statusCode int, data interface{}) HTTPResponse {
	switch statusCode {
	case http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusRequestTimeout:

		if e, ok := data.(error); ok {
			return HTTPResponse{
				StatusCode: statusCode,
				Success:    false,
				Message:    e.Error(),
			}
		}
		return HTTPResponse{
			StatusCode: statusCode,
			Success:    false,
			Message:    data.(string),
		}
	default:
		return HTTPResponse{
			StatusCode: statusCode,
			Success:    true,
			Data:       data,
		}
	}
}
