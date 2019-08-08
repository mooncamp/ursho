package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config contains the configuration of the url shortener.
type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	Redis struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		DB       string `json:"db"`
	} `json:"redis"`
	Postgres struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DB       string `json:"db"`
	} `json:"postgres"`
	Dgraph struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
	Crypto struct {
		Key   string `json:"key"`
		Nonce string `json:"nonce"`
	}
	Options struct {
		Prefix   string `json:"prefix"`
		Database string `json:"database"`
		Encoding string `json:"encoding"`
	} `json:"options"`
}

// FromFile returns a configuration parsed from the given file.
func FromFile(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
