package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RetrieveProject(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []map[string]interface{}{
			{
				"id":   1,
				"name": "demo",
			},
		},
	})
}
