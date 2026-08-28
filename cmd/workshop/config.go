package main

import (
	"fmt"
	"strings"
)

type Config struct {
	DatabasePath string
	Address      string
	ReadOnly     bool
}

func DefaultConfig() Config { return Config{DatabasePath: "workshopnotice.db", Address: ":8080"} }

func ParseConfig(values map[string]string) (Config, error) {
	config := DefaultConfig()
	for key, value := range values {
		switch key {
		case "db":
			if strings.TrimSpace(value) == "" {
				return Config{}, fmt.Errorf("db cannot be empty")
			}
			config.DatabasePath = strings.TrimSpace(value)
		case "addr":
			if strings.TrimSpace(value) == "" {
				return Config{}, fmt.Errorf("addr cannot be empty")
			}
			config.Address = strings.TrimSpace(value)
		case "readonly":
			config.ReadOnly = value == "true" || value == "1"
		default:
			return Config{}, fmt.Errorf("unknown option %s", key)
		}
	}
	return config, nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.DatabasePath) == "" || strings.TrimSpace(c.Address) == "" {
		return fmt.Errorf("database and address are required")
	}
	return nil
}

func (c Config) String() string {
	return fmt.Sprintf("db=%s addr=%s readonly=%t", c.DatabasePath, c.Address, c.ReadOnly)
}
