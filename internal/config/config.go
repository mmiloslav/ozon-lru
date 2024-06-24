// пакет конфигурации сервиса
package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v8"
	log "github.com/sirupsen/logrus"
)

const defaultLevel log.Level = log.WarnLevel // Уровень логированя по умолчанию

// тип с конфигурацией настроек сервиса
type Conf struct {
	ServerHostPort  string        `env:"SERVER_HOST_PORT" envDefault:"localhost:8080"`
	CacheSize       int           `env:"CACHE_SIZE" envDefault:"10"`
	DefaultCacheTTL time.Duration `env:"DEFAULT_CACHE_TTL" envDefault:"1m"`
	LogLevel        string        `env:"LOG_LEVEL" envDefault:"WARN"`
}

// инициализация конфигурации
func InitConf() (Conf, error) {
	var conf Conf
	if err := env.Parse(&conf); err != nil {
		log.Errorf("failed to parse config env vars with error [%s]", err.Error())
		return Conf{}, err
	}

	serverHostPort := flag.String("server-host-port", conf.ServerHostPort, "Server host and port")
	cacheSize := flag.Int("cache-size", conf.CacheSize, "Cache size")
	defaultCacheTTL := flag.Duration("default-cache-ttl", conf.DefaultCacheTTL, "Default cache TTL")
	logLevel := flag.String("log-level", conf.LogLevel, "Log level")

	flag.Parse()

	conf.ServerHostPort = *serverHostPort
	conf.CacheSize = *cacheSize
	conf.DefaultCacheTTL = *defaultCacheTTL
	conf.LogLevel = *logLevel

	log.Infof("config [%+v]", conf)

	return conf, nil
}

// устанавливает уровень логирования
func SetLogLevel(lvl string) {
	level, err := log.ParseLevel(lvl)
	if err != nil {
		log.Errorf("failed to parse log level [%s]. Using default value [%d]", lvl, defaultLevel)
		level = defaultLevel
	}

	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
