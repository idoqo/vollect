package response

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrorResponse struct {
	Err error `json:"-"`
	StatusCode int `json:"-"`
	StatusText string `json:"status_text"`
	Message string `json:"message"`
}

var (
	ErrMethodNotAllowed = &ErrorResponse{StatusCode: 405, StatusText: "Not Allowed", Message: "Method now allowed"}
	ErrNotFound = &ErrorResponse{StatusCode: 404, StatusText: "Not Found", Message: "Resource not found"}
)

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func ErrServerError(err error) *ErrorResponse {
	return &ErrorResponse{
		Err: err,
		StatusCode: 500,
		StatusText: "Internal server error",
		Message: err.Error(),
	}
}

func ErrBadRequest(err error) *ErrorResponse {
	return &ErrorResponse{
		Err: err,
		StatusCode: 400,
		StatusText: "Bad request",
		Message: err.Error(),
	}
}
