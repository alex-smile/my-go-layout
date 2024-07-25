package middleware

import (
	"github.com/gin-gonic/gin"

	"mygo/template/pkg/middleware/context"
	"mygo/template/pkg/util"
)

// RequestIDHeaderKey is a key to set the request id in header
const (
	RequestIDHeaderKey = "X-Request-Id"
)

// RequestID add the request_id for each api request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeaderKey)
		if requestID == "" || len(requestID) != 32 {
			requestID = util.GenUUID4()
		}
		context.SetRequestID(c, requestID)
		c.Writer.Header().Set(RequestIDHeaderKey, requestID)

		c.Next()
	}
}
