package basic

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"mygo/template/pkg/config"
)

func Healthz(c *gin.Context) {
	if err := checkDatabase(nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"database": fmt.Sprintf("failed to connect to db, err %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"database": "ok"})
}

func checkDatabase(dbConfig *config.Database) error {
	// TODO: add database check
	return nil
}
