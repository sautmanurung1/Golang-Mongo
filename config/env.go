package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)
func EnvMongoURI() string{
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error Loading .env files")
	}

	return os.Getenv("MONGOURI")
}