package configs

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
)

const projectDir = "traverse" // change to project name directory

type Config struct {
	MIGRATIONS *MigrationConfig
	SERVER     *ServerConfig
	DB         *DBConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	// Debug        bool
}

type DBConfig struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

type MigrationConfig struct {
	DIR string
}

var ENVS = initEnvs()

func initEnvs() *Config {
	godotenv.Load()

	return &Config{
		MIGRATIONS: &MigrationConfig{
			DIR: getEnv("MIGRATION_DIR", "./cmd/migrates/migrations/"),
		},
		SERVER: &ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsTime("SERVER_TIMEOUT_READ", 5*time.Second),
			WriteTimeout: getEnvAsTime("SERVER_TIMEOUT_WRITE", 10*time.Second),
			IdleTimeout:  getEnvAsTime("SERVER_TIMEOUT_IDLE", 15*time.Second),
		},
		DB: &DBConfig{
			Name:     getEnv("DB_NAME", "traverse_db"),
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "8080"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "sendhelp"),
		},
	}
}

func (c *DBConfig) String() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
	)
}

func getRootProjectName() string {
	projName := regexp.MustCompile("^(.*" + projectDir + ")")
	CWD, _ := os.Getwd()
	rootPath := projName.Find([]byte(CWD))

	return string(rootPath) + "/.env"
}
