package cmd

import (
	"fmt"

	"github.com/getsentry/raven-go"
	"github.com/getsentry/sentry-go"

	"mygo/template/pkg/config"
	"mygo/template/pkg/infra/logging"
	"mygo/template/pkg/infra/metric"
	"mygo/template/pkg/infra/redis"
	"mygo/template/pkg/store/database"
)

var globalConfig *config.Config

func initConfig() {
	var err error
	globalConfig, err = config.Load(cfgFile)
	if err != nil {
		panic(fmt.Errorf("failed to load config, err: %v", err))
	}
}

func initDatabase() {
	database.InitDBClients(globalConfig.Databases)
}

func initRedis() {
	redis.InitRedisClient(globalConfig.Redis)
}

func initSentry() {
	if globalConfig.Sentry.DSN == "" {
		logging.GetLogger().Info("sentry is not enabled, will not init it")
		return
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: globalConfig.Sentry.DSN,
	})
	if err != nil {
		logging.GetLogger().Errorf("init Sentry fail: %v", err)
		return
	}
	logging.GetLogger().Info("init Sentry success")

	// init gin sentry
	err = raven.SetDSN(globalConfig.Sentry.DSN)
	if err != nil {
		logging.GetLogger().Errorf("init gin Sentry fail: %s", err)
		return
	}
	logging.GetLogger().Info("init gin Sentry success")
}

func initMetrics() {
	metric.InitMetrics()
}
