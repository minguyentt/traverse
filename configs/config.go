package configs

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

const projectDir = "traverse" // change to project name directory

type Config struct {
	MIGRATIONS *MigrationConfig
	SERVER     *ServerConfig
	DB         *DBConfig
	DEV_DB     *Local_DBConfig
	AUTH       *AuthConfig
	REDIS      *RedisConfig
	RATELIMIT  *RateLimitConfig
}

type AuthConfig struct {
	Token apiToken
	Admin adminConfig
}

type adminConfig struct {
	Username string
	Password string
}

type apiToken struct {
	Secret string
	Exp    time.Duration
	Iss    string
	Aud    string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Debug        bool
}

type DBConfig struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

type Local_DBConfig struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

type MigrationConfig struct {
	DIR string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Enabled  bool
}

type RateLimitConfig struct {
	Standard *RateLimitOpts
	Test     *RateLimitOpts
	Enabled  bool
}

type RateLimitOpts struct {
	Buckets uint
	Depth   uint

	Limit  int
	Window time.Duration
	NumWin int
}

var Env = LoadEnvs()

func LoadEnvs() *Config {
	return &Config{
		MIGRATIONS: &MigrationConfig{
			DIR: getEnv("MIGRATION_DIR"),
		},
		SERVER: &ServerConfig{
			Port:         getEnv("SERVER_PORT"),
			ReadTimeout:  getEnvAsTime("SERVER_TIMEOUT_READ", 5*time.Second),
			WriteTimeout: getEnvAsTime("SERVER_TIMEOUT_WRITE", 10*time.Second),
			IdleTimeout:  getEnvAsTime("SERVER_TIMEOUT_IDLE", 15*time.Second),
		},
		DB: &DBConfig{
			Name:     getEnv("DB_NAME"),
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
		},
		DEV_DB: &Local_DBConfig{
			Name:     getEnv("LOCAL_DB_NAME"),
			Host:     getEnv("LOCAL_DB_HOST"),
			Port:     getEnv("LOCAL_DB_PORT"),
			User:     getEnv("LOCAL_DB_USER"),
			Password: getEnv("LOCAL_DB_PASSWORD"),
		},
		AUTH: &AuthConfig{
			Token: apiToken{
				Secret: getEnv("AUTH_API_KEY"),
				Exp:    getEnvAsTime("AUTH_EXP_TIME", 24*time.Hour),
				Iss:    getEnv("AUTH_ISS"),
				Aud:    getEnv("AUTH_AUD"),
			},
			Admin: adminConfig{
				Username: getEnv("AUTH_ADMIN_USER"),
				Password: getEnv("AUTH_ADMIN_PASS"),
			},
		},
		REDIS: &RedisConfig{
			Addr:     getEnv("REDIS_CLIENT_ADDR"),
			Password: getEnv("REDIS_CLIENT_PASSWORD"),
			DB:       0,
			Enabled:  true,
		},
		RATELIMIT: &RateLimitConfig{
			Standard: &RateLimitOpts{
				Buckets: 10000,
				Depth:   4,
				Limit:   100,
				Window:  time.Minute,
				NumWin:  5,
			},
			Test: &RateLimitOpts{
				Buckets: 1000,
				Depth:   3,
				Limit:   10,
				Window:  time.Second,
				NumWin:  3,
			},
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

// local dev testing
func (c *Local_DBConfig) String() string {
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
