package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Mongo struct {
		URI             string `mapstructure:"uri"`
		Database        string `mapstructure:"database"`
		UserCollection  string `mapstructure:"user_collection"`
		TokenCollection string `mapstructure:"token_collection"`
		PostCollection  string `mapstructure:"post_collection"`
	} `mapstructure:"mongo"`

	JWT struct {
		AccessTokenSecret  string `mapstructure:"access_token_secret"`
		RefreshTokenSecret string `mapstructure:"refresh_token_secret"`
	} `mapstructure:"jwt"`

	HMAC struct {
		Secret string `mapstructure:"hmac_secret"`
	} `mapstructure:"hmac"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config.dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
