package database

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"

	"mygo/template/pkg/config"
)

type DBClient struct {
	DB *sqlx.DB

	dataSource string

	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
}

func NewDBClient(cfg *config.Database) *DBClient {
	dataSource := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&interpolateParams=true&loc=%s&time_zone=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		"utf8",
		"UTC",
		url.QueryEscape("'+00:00'"),
	)

	return &DBClient{
		dataSource:      dataSource,
		maxOpenConns:    cfg.GetMaxOpenConns(),
		maxIdleConns:    cfg.GetMaxIdleConns(),
		connMaxLifetime: cfg.GetConnMaxLifetime(),
	}
}

func (db *DBClient) Connect() error {
	var err error
	db.DB, err = sqlx.Connect("mysql", db.dataSource)
	if err != nil {
		return err
	}

	db.DB.SetMaxOpenConns(db.maxOpenConns)
	db.DB.SetMaxIdleConns(db.maxIdleConns)
	db.DB.SetConnMaxLifetime(db.connMaxLifetime)

	_, err = db.DB.Exec(`SET time_zone = "+00:00";`) // set session time zon to utc
	if err != nil {
		return err
	}

	return nil
}

func (db *DBClient) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}
