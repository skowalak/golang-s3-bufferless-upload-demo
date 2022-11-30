package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	Address string `yaml:"address"`
}

type AwsConfig struct {
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	KeyId     string `yaml:"aws-access-key-id"`
	KeySecret string `yaml:"aws-access-key-secret"`
}

type Config struct {
	ApplicationConfig `yaml:"application"`
	AwsConfig         `yaml:"aws"`
}

func ProvideConfig() *Config {
	conf := Config{}
	data, err := ioutil.ReadFile("config/base.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		// :(
		panic(err)
	}

	return &conf
}
