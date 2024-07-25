package context

import (
	"github.com/gin-gonic/gin"
)

const (
	requestIDKey = "x_request_id"
)

func SetRequestID(c *gin.Context, requestID string) {
	c.Set(requestIDKey, requestID)
}

func GetRequestID(c *gin.Context) string {
	return c.GetString(requestIDKey)
}
