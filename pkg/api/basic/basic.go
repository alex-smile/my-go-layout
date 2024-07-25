package basic

import (
	"mygo/template/pkg/version"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Pong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func Version(c *gin.Context) {
	now := time.Now()

	c.JSON(http.StatusOK, gin.H{
		"version":   version.Version,
		"commit":    version.Commit,
		"buildTime": version.BuildTime,
		"goVersion": version.GoVersion,
		// add the date and timestamp
		"date":      now,
		"timestamp": now.Unix(),
	})
}
