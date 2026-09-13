package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   serverConfig
	Database databaseConfig
	MaxBot   maxBotConfig
	Uploads  uploadsConfig
}

type serverConfig struct {
	Port string
	Mode string
}

type databaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int32
	MinConns int32
}

type maxBotConfig struct {
	Token  string
	ChatID string
	Active bool
}

type uploadsConfig struct {
	Path       string
	MaxSizeMB  int
}

func Load() *Config {
	cfg := &Config{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Sprintf("config error: %v", err))
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("KONTURSVET")

	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "SERVER_MODE")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSLMODE")
	viper.BindEnv("database.max_conns", "DB_MAX_CONNS")
	viper.BindEnv("database.min_conns", "DB_MIN_CONNS")
	viper.BindEnv("max_bot.token", "MAX_BOT_TOKEN")
	viper.BindEnv("max_bot.chat_id", "MAX_BOT_CHAT_ID")
	viper.BindEnv("max_bot.active", "MAX_BOT_ACTIVE")
	viper.BindEnv("uploads.path", "UPLOADS_PATH")
	viper.BindEnv("uploads.max_size_mb", "UPLOADS_MAX_SIZE_MB")

	cfg.Server.Port = viper.GetString("server.port")
	cfg.Server.Mode = viper.GetString("server.mode")
	cfg.Database.Host = viper.GetString("database.host")
	cfg.Database.Port = viper.GetString("database.port")
	cfg.Database.User = viper.GetString("database.user")
	cfg.Database.Password = viper.GetString("database.password")
	cfg.Database.DBName = viper.GetString("database.dbname")
	cfg.Database.SSLMode = viper.GetString("database.sslmode")
	cfg.Database.MaxConns = viper.GetInt32("database.max_conns")
	cfg.Database.MinConns = viper.GetInt32("database.min_conns")
	cfg.MaxBot.Token = viper.GetString("max_bot.token")
	cfg.MaxBot.ChatID = viper.GetString("max_bot.chat_id")
	cfg.MaxBot.Active = viper.GetBool("max_bot.active")
	cfg.Uploads.Path = viper.GetString("uploads.path")
	cfg.Uploads.MaxSizeMB = viper.GetInt("uploads.max_size_mb")

	return cfg
}