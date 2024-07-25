package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	// server
	defaultServerReadTimeout  = 180 * time.Second
	defaultServerWriteTimeout = 180 * time.Second
	defaultServerIdleTimeout  = 60 * time.Second

	defaultServerGraceTimeout = 60 * time.Second
)

var (
	// database
	defaultDBMaxOpenConns    = 100
	defaultDBMaxIdleConns    = 10
	defaultDBConnMaxLifetime = 600 * time.Second
)

type Config struct {
	Debug bool `mapstructure:"debug"`

	Server Server `mapstructure:"server"`

	Databases map[string]*Database `mapstructure:"databases"`
	Redis     map[string]*Redis    `mapstructure:"redis"`
	Sentry    Sentry               `mapstructure:"sentry"`

	Hosts map[string]*Host `mapstructure:"hosts"`

	Cache Cache `mapstructure:"cache"`
}

type Server struct {
	Addr string `mapstructure:"addr"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	ReadTimeout  int `mapstructure:"readTimeout"`
	WriteTimeout int `mapstructure:"writeTimeout"`
	IdleTimeout  int `mapstructure:"idleTimeout"`

	GraceTimeout int `mapstructure:"graceTimeout"`
}

func (s Server) GetAddr() string {
	if s.Addr != "" {
		return s.Addr
	}

	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s Server) GetReadTimeout() time.Duration {
	if s.ReadTimeout > 0 {
		return time.Duration(s.ReadTimeout) * time.Second
	}
	return defaultServerReadTimeout
}

func (s Server) GetWriteTimeout() time.Duration {
	if s.WriteTimeout > 0 {
		return time.Duration(s.WriteTimeout) * time.Second
	}

	return defaultServerWriteTimeout
}

func (s Server) GetIdleTimeout() time.Duration {
	if s.IdleTimeout > 0 {
		return time.Duration(s.IdleTimeout) * time.Second
	}

	return defaultServerIdleTimeout
}

func (s Server) GetGraceTimeout() time.Duration {
	if s.GraceTimeout > 0 {
		return time.Duration(s.GraceTimeout) * time.Second
	}

	return defaultServerGraceTimeout
}

type Cache struct {
	Enable bool `mapstructure:"enable"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`

	MaxOpenConns          int `mapstructure:"maxOpenConns"`
	MaxIdleConns          int `mapstructure:"maxIdleConns"`
	ConnMaxLifetimeSecond int `mapstructure:"connMaxLifetimeSecond"`
}

func (d Database) GetMaxOpenConns() int {
	if d.MaxOpenConns > 0 {
		return d.MaxOpenConns
	}

	return defaultDBMaxOpenConns
}

func (d Database) GetMaxIdleConns() int {
	if d.MaxIdleConns > 0 {
		return d.MaxIdleConns
	}

	return defaultDBMaxIdleConns
}

func (d Database) GetConnMaxLifetime() time.Duration {
	if d.ConnMaxLifetimeSecond > 0 {
		return time.Duration(d.ConnMaxLifetimeSecond) * time.Second
	}

	return defaultDBConnMaxLifetime
}

type Redis struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	DialTimeout  int    `mapstructure:"dialTimeout"`
	ReadTimeout  int    `mapstructure:"readTimeout"`
	WriteTimeout int    `mapstructure:"writeTimeout"`
	PoolSize     int    `mapstructure:"poolSize"`
	MinIdleConns int    `mapstructure:"minIdleConns"`
	IdleTimeout  int    `mapstructure:"idleTimeout"`
}

type Sentry struct {
	DSN string `mapstructure:"dsn"`
}

type Host struct {
	Addr string `mapstructure:"addr"`
}

func Load(cfgFile string) (*Config, error) {
	if cfgFile == "" {
		return nil, errors.New("config file missing")
	}

	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	// viper.SetEnvPrefix("MYGO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load config file")
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	return &cfg, nil
}
