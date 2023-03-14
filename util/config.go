package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	SecretKey            string        `mapstructure:"SECRET_KEY"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	ServerAddr           string        `mapstructure:"SEVER_ADDR"`
	DBHOST               string        `mapstructure:"DB_HOST"`
	DBPort               string        `mapstructure:"DB_PORT"`
	DBUser               string        `mapstructure:"DB_USER"`
	DBPassword           string        `mapstructure:"DB_PASSWORD"`
	DBName               string        `mapstructure:"DB_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`

	// DBSource            string        `mapstructure:"DB_SOURCE"`
}

func (config *Config) DBSource() string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", config.DBUser, config.DBPassword, config.DBHOST, config.DBPort, config.DBName)
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return config, err

}
