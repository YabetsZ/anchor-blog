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

	GenAI struct {
		GeminiAPIKey string `mapstructure:"gemini_api_key"`
		GeminiModel  string `mapstructure:"gemini_model"`
		MinWords     int    `mapstructure:"min_words"`
		MaxWords     int    `mapstructure:"max_words"`
	} `mapstructure:"genai"`

	Redis struct {
		Host            string `mapstructure:"host"`
		Port            string `mapstructure:"port"`
		Password        string `mapstructure:"password"`
		DB              int    `mapstructure:"db"`
		ViewTrackingTTL int    `mapstructure:"view_tracking_ttl"`
	} `mapstructure:"redis"`

	OAuth struct {
		Google struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
			RedirectURI  string `mapstructure:"redirect_uri"`
		} `mapstructure:"google"`
	} `mapstructure:"oauth"`
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
