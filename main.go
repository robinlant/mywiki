package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/robinlant/mywiki/internal/store"
	"github.com/robinlant/mywiki/internal/web"
)

type Config struct {
	Addr     string
	MongoURI string
	DB       string
}

func LoadConfig() Config {
	return Config{
		Addr:     mustGetenv("ADDR"),
		MongoURI: mustGetenv("MONGO_CON"),
		DB:       mustGetenv("MONGO_DB"),
	}
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Errorf("failed to get env var '%s'", key))
	}
	return v
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] No .env file found, using system env")
	}
	conf := LoadConfig()

	st, disc := store.NewMongoStore(conf.MongoURI, conf.DB)
	defer disc()
	log.Printf("[INFO] Starting a web server at %s", conf.Addr)
	web.Run(st, conf.Addr)
}
