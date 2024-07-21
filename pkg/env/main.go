package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const PORT = 0

var keys = []string {
	"PORT",
}

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func Read(key int) string {
	if key < 0 || key >= len(keys) {
		log.Fatalf("Invalid env key")
	}
	return os.Getenv(keys[key])
}
