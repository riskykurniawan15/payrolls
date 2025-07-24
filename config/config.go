package config

import (
	"fmt"
	"log"

	"github.com/riskykurniawan15/payrolls/utils/env"
)

type (
	Config struct {
		CompanyName string
		Http        HttpServer
		PostgressDB PostgressDB
		JWT         JWTConfig
		Logger      LoggerConfig
	}

	HttpServer struct {
		Server string
		Port   int
		URL    string
	}

	PostgressDB struct {
		DBUser        string
		DBPass        string
		DBServer      string
		DBPort        int
		DBName        string
		DBMaxIdleCon  int
		DBMaxOpenCon  int
		DBMaxLifeTime int
		DBTimeZone    string
		SSLMode       string
		DBDebug       bool
	}

	JWTConfig struct {
		SecretKey string
		Expired   int
	}

	LoggerConfig struct {
		OutputMode string
		LogLevel   string
		LogDir     string
	}
)

func Configuration() Config {
	if err := env.LoadEnv(".env"); err != nil {
		log.Println("error read .env file %w", err.Error())
	}

	cfg := Config{
		CompanyName: env.GetEnv("COMPANY_NAME", "Blank Company"),
		Http:        loadHttpServer(),
		PostgressDB: loadDBServer(),
		JWT:         loadJWTConfig(),
		Logger:      loadLoggerConfig(),
	}

	log.Println("Success for load all configuration")

	return cfg
}

func loadHttpServer() HttpServer {
	var cfg HttpServer

	cfg.Server = env.GetEnv("SERVER", "localhost")
	cfg.Port = env.GetEnv("PORT", 9000)
	if env.GetEnv("USING_SECURE", true) {
		cfg.URL = "https://" + cfg.Server
	} else {
		cfg.URL = "http://" + cfg.Server
	}

	if cfg.Port != 0 {
		cfg.URL += fmt.Sprintf(":%d", cfg.Port)
	}
	cfg.URL += "/"

	return cfg
}

func loadDBServer() PostgressDB {
	return PostgressDB{
		DBUser:        env.GetEnv("DB_USER", "root"),
		DBPass:        env.GetEnv("DB_PASS", ""),
		DBServer:      env.GetEnv("DB_SERVER", "localhost"),
		DBPort:        env.GetEnv("DB_PORT", 5432),
		DBName:        env.GetEnv("DB_NAME", "public"),
		DBMaxIdleCon:  env.GetEnv("DB_MAX_IDLE_CON", 10),
		DBMaxOpenCon:  env.GetEnv("DB_MAX_OPEN_CON", 100),
		DBMaxLifeTime: env.GetEnv("DB_MAX_LIFE_TIME", 10),
		DBTimeZone:    env.GetEnv("DB_TIME_ZONE", "Asia/Jakarta"),
		SSLMode:       env.GetEnv("DB_SSL_MODE", "disable"),
		DBDebug:       env.GetEnv("DB_DEBUG", false),
	}
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey: env.GetEnv("JWT_SECRET_KEY", ""),
		Expired:   env.GetEnv("JWT_EXPIRED", 24),
	}
}

func loadLoggerConfig() LoggerConfig {
	return LoggerConfig{
		OutputMode: env.GetEnv("LOG_OUTPUT_MODE", "both"), // "terminal", "file", "both"
		LogLevel:   env.GetEnv("LOG_LEVEL", "debug"),      // "debug", "info", "warn", "error"
		LogDir:     env.GetEnv("LOG_DIR", "logger"),       // directory for log files
	}
}
