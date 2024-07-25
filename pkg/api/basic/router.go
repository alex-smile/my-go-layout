package basic

import (
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AddRouter(r *gin.Engine) {
	r.GET("/_ping", Pong)
	r.GET("/_healthz", Healthz)
	r.GET("/_version", Version)

	// metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// pprof
	pprofGroup := r.Group("/debug/pprof")
	{
		pprofGroup.GET("", gin.WrapF(pprof.Index))
		pprofGroup.GET("/profile/", gin.WrapF(pprof.Profile))
		pprofGroup.GET("/cmdline/", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/symbol/", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace/", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs/", gin.WrapH(pprof.Handler("allocs")))
		pprofGroup.GET("/block/", gin.WrapH(pprof.Handler("block")))
		pprofGroup.GET("/goroutine/", gin.WrapH(pprof.Handler("goroutine")))
		pprofGroup.GET("/heap/", gin.WrapH(pprof.Handler("heap")))
		pprofGroup.GET("/mutex/", gin.WrapH(pprof.Handler("mutex")))
		pprofGroup.GET("/threadcreate/", gin.WrapH(pprof.Handler("threadcreate")))
	}
}