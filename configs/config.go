package configs

import (
	"time"
)

// TODO: fix config structure for better modularity

type Config struct {
	Addr   string
	Port   string
	Server *ServerConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	// Environment
	// Debug
}

func (c *ServerConfig) GetAddr() string {
	return c.Host + c.Port
}

func (c *ServerConfig) GetPort() string {
	return ":" + c.Port
}

func LoadConfig() *Config {
	server := NewDefaultServerConfig()

	return &Config{
		Addr:   server.GetAddr(),
		Port:   server.GetPort(),
		Server: server,
	}
}

func NewDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:         "localhost",
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
