package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Backend error types
var (
	ErrInvalidRequest     = errors.New("invalid_request")
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrUnauthorizedClient = errors.New("unauthorized_client")
	ErrAccessDenied       = errors.New("access_denied")
	ErrServerError        = errors.New("server_error")
	ErrNotFound           = errors.New("not_found")
)

type APIError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func respondError(c *gin.Context, err error) {
	var apiErr APIError
	var status int

	switch err {
	case ErrInvalidCredentials:
		status = http.StatusUnauthorized
		apiErr = APIError{Error: "invalid_credentials", ErrorDescription: "incorrect email or password"}
	case ErrUnauthorizedClient:
		status = http.StatusUnauthorized
		apiErr = APIError{Error: "unauthorized_client"}
	case ErrInvalidRequest:
		status = http.StatusBadRequest
		apiErr = APIError{Error: "invalid_request"}
	case ErrNotFound:
		status = http.StatusNotFound
		apiErr = APIError{Error: "not_found"}
	default:
		status = http.StatusInternalServerError
		apiErr = APIError{Error: "server_error"}
	}

	c.JSON(status, apiErr)
}
