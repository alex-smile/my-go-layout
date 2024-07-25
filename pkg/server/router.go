package server

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/sentry"
	"github.com/gin-gonic/gin"

	"mygo/template/pkg/api/basic"
	"mygo/template/pkg/api/project"
	"mygo/template/pkg/config"
	"mygo/template/pkg/middleware"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.RequestID())
	// metrics
	router.Use(middleware.Metrics())
	// recovery sentry
	router.Use(sentry.Recovery(raven.DefaultClient, false))

	basic.AddRouter(router)

	v1 := router.Group("/v1")
	{
		project.AddRouter(v1)
	}

	return router
}
