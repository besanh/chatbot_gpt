package config

import (
	"strings"

	"github.com/besanh/chatbot_gpt/common/caching"
	"github.com/besanh/chatbot_gpt/pkg/redis"
	log "github.com/besanh/logger/logging/slog"
	"github.com/caarlos0/env"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port             string `yaml:"port"`
		Mode             string `yaml:"mode"`
		LogLevel         string `yaml:"log_level"`
		LogFile          string `yaml:"log_file"`
		GracefulShutdown struct {
			ShutdownTime int64 `yaml:"shutdown_time"`
			ReadTimeout  int64 `yaml:"read_timeout"`
			WriteTimeout int64 `yaml:"write_timeout"`
			IdleTimeout  int64 `yaml:"idle_timeout"`
		} `yaml:"graceful_shutdown"`
	}

	Api struct {
		ApiServiceName string `yaml:"api_service_name"`
		ApiVersion     string `yaml:"api_version"`
		TlsCertPath    string `yaml:"tls_cert_path"`
		TlsKeyPath     string `yaml:"tls_key_path"`
	} `yaml:"api"`

	Pkg struct {
		Openai struct {
			ApiKey string   `yaml:"api_key"`
			Models []string `yaml:"models"`
		} `yaml:"openai"`

		Redis struct {
			Dsn string `yaml:"dsn"`
		} `yaml:"redis"`
	}
}

func InitConfig(cfg *Config) {
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	initLogger(*cfg)
	initRedis(*cfg)

	registerMetrics()
}

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func registerMetrics() {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		HttpRequestsTotal,
		HttpRequestDuration,
		collectors.NewGoCollector(), // Go runtime metrics
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}), // Process metrics
	)
}

func initLogger(cfg Config) {
	logLevel := log.LEVEL_DEBUG
	switch cfg.Server.LogLevel {
	case "debug":
		logLevel = log.LEVEL_DEBUG
	case "info":
		logLevel = log.LEVEL_INFO
	case "error":
		logLevel = log.LEVEL_ERROR
	case "warn":
		logLevel = log.LEVEL_WARN
	}
	opts := []log.Option{}
	opts = append(opts, log.WithLevel(logLevel),
		log.WithRotateFile(cfg.Server.LogFile),
		log.WithFileSource(),
	)

	log.SetLogger(log.NewSLogger(opts...))
}

func initRedis(cfg Config) {
	redisClient, err := redis.NewRedis(redis.RedisConfig{Dsn: cfg.Pkg.Redis.Dsn})
	if err != nil {
		panic(err)
	}

	caching.RCache = caching.NewRedisCache(redisClient.GetClient())
}
