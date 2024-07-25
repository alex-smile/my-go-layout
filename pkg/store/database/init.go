package database

import (
	"fmt"
	"sync"

	"mygo/template/pkg/config"
)

var (
	defaultDBClient *DBClient

	defaultDBClientOnce sync.Once
)

func InitDBClients(dbConfig map[string]*config.Database) {
	defaultDBConfig, ok := dbConfig["default"]
	if !ok {
		panic("database default should be configured")
	}
	initDefaultDBClient(defaultDBConfig)
}

func initDefaultDBClient(dbConfig *config.Database) {
	if defaultDBClient != nil {
		return
	}

	defaultDBClientOnce.Do(func() {
		defaultDBClient = NewDBClient(dbConfig)
		if err := defaultDBClient.Connect(); err != nil {
			panic(fmt.Errorf("failed to connect to database, err: %v", err))
		}
	})
}

func GetDefaultDBClient() *DBClient {
	return defaultDBClient
}
