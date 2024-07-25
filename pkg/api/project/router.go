package project

import "github.com/gin-gonic/gin"

func AddRouter(r *gin.RouterGroup) {
	r.GET("/projects/:id", RetrieveProject)
}
