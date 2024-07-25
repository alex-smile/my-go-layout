package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	systemSource = "mygo"
)

const (
	BadRequestError   = "BadRequest"
	UnauthorizedError = "Unauthorized"
	ForbiddenError    = "Forbidden"
	NotFoundError     = "NotFound"
	ConflictError     = "Conflict"
	TooManyRequests   = "TooManyRequests"
	SystemError       = "InternalServerError"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Source  string `json:"source"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

func SuccessJSONResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{Data: data})
}

func BaseErrorJSONResponse(c *gin.Context, code string, message string, statusCode int) {
	c.JSON(statusCode, ErrorResponse{Error: Error{
		Code:    code,
		Message: message,
		Source:  systemSource,
	}})
}

func NewErrorJSONResponse(code string, statusCode int) func(c *gin.Context, message string) {
	return func(c *gin.Context, message string) {
		BaseErrorJSONResponse(c, code, message, statusCode)
	}
}

var (
	BadRequestErrorJSONResponse = NewErrorJSONResponse(BadRequestError, http.StatusBadRequest)
	ForbiddenJSONResponse       = NewErrorJSONResponse(ForbiddenError, http.StatusForbidden)
	UnauthorizedJSONResponse    = NewErrorJSONResponse(UnauthorizedError, http.StatusUnauthorized)
	NotFoundJSONResponse        = NewErrorJSONResponse(NotFoundError, http.StatusNotFound)
	ConflictJSONResponse        = NewErrorJSONResponse(ConflictError, http.StatusConflict)
	TooManyRequestsJSONResponse = NewErrorJSONResponse(TooManyRequests, http.StatusTooManyRequests)
)

func SystemErrorJSONResponse(c *gin.Context, err error) {
	BaseErrorJSONResponse(c, SystemError, err.Error(), http.StatusInternalServerError)
}

func NewErrorJSONResponseWithError(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	if x, ok := err.(interface{ GetStatusCode() int }); ok {
		statusCode = x.GetStatusCode()
	}

	BaseErrorJSONResponse(c, SystemError, err.Error(), statusCode)
}
