package config

//Инструкция с кодом https://golangwiki.ru/article/golang-viper-prakticheskoe-rukovodstvo
import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		User     string
		Password string
		Host     string
		Port     int
		DBName   string
		SSLMode  string
	}
	Server struct {
		Port int
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
