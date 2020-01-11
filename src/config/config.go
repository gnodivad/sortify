package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type SpotifyConfig struct {
	ClientID  string
	SecretKey string
}

type Config struct {
	Spotify SpotifyConfig
}

func Init() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file found")
	}

	return &Config{
		Spotify: SpotifyConfig{
			ClientID:  getEnv("SPOTIFY_ID", ""),
			SecretKey: getEnv("SPOTIFY_SECRET", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}