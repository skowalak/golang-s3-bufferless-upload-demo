package config

import (
	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	Address string `mapstructure:"address"`
}

type AwsConfig struct {
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	KeyId     string `mapstructure:"aws-access-key-id"`
	KeySecret string `mapstructure:"aws-access-key-secret"`
}

type Config struct {
	ApplicationConfig `mapstructure:"application"`
	AwsConfig         `mapstructure:"aws"`
}

func ProvideConfig() *Config {
	viper.SetConfigName("base")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	conf := Config{}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
