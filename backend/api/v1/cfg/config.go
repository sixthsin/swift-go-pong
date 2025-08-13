package cfg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db      DbConfig
	Storage StorageConfig
}

type DbConfig struct {
	Dsn string
}

type StorageConfig struct {
	Path string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error while loading .env file")
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Storage: StorageConfig{
			Path: os.Getenv("STORAGE_PATH"),
		},
	}
}
