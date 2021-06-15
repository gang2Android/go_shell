package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Name   string       `yaml:"name"`
	Port   string       `yaml:"port"`
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Ip   string `yaml:"ip"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pwd  string `yaml:"pwd"`
}

func LoadConfig(configPath string) *Config {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	config := &Config{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}
	return config
}
