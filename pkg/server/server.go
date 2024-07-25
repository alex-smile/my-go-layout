package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mygo/template/pkg/config"
	"mygo/template/pkg/infra/logging"
)

func Run(cfg *config.Config) {
	router := NewRouter(cfg)

	s := &http.Server{
		Addr:         cfg.Server.GetAddr(),
		Handler:      router,
		ReadTimeout:  cfg.Server.GetReadTimeout(),
		WriteTimeout: cfg.Server.GetWriteTimeout(),
		IdleTimeout:  cfg.Server.GetIdleTimeout(),
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.GetLogger().Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.GetLogger().Info("Shutting down server...")

	// The context is used to inform the server it has grace timeout seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GetGraceTimeout())
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logging.GetLogger().Fatal("Server forced to shutdown: ", err)
	}

	// nolint:gosimple
	select {
	case <-ctx.Done():
		logging.GetLogger().Infof("Server exiting by timeout of %d seconds", cfg.Server.GetGraceTimeout()/time.Second)
	}

	logging.GetLogger().Println("Server exiting")
}
