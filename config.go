package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACKey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func DefaultConfig() Config {
	return Config{
		Port:     3000,
		Env:      "dev",
		Pepper:   "cC242xTzSG!6j!mWd2N3Vg3!!Q38wunu23a6YBUTm@e**GyP@!CyAzjW7JcR7*p!^524sNxs9H7RQkh3^xH3Q4eSFQtQNqnXqW!",
		HMACKey:  "cC242xTzSG!6j!mWd2N3Vh4!!Q38wunu23a6YBUTm@e**GyP@!CyAzjW7JcR7*p!^524sNxs9H7RQkh3^xH3Q4eSFQtQNqnXqW!",
		Database: DefaultPosgresConfig(),
	}
}

func DefaultPosgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		Name:     "lenslocked",
	}
}

func LoadConfig(isProd bool) Config {
	f, err := os.Open(".config")
	if err != nil {
		if isProd {
			panic(err)
		}
		fmt.Println("Using the default config...")
		return DefaultConfig()
	}
	var c Config
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded .config")
	return c
}
