package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/robinlant/mywiki/wiki/internal/quote"
	"github.com/robinlant/mywiki/wiki/internal/store"
	"github.com/robinlant/mywiki/wiki/internal/web"
)

type Config struct {
	Addr      string
	MongoURI  string
	DB        string
	QuotesUrl string
}

func LoadConfig() Config {
	return Config{
		Addr:      getenvOrDefault("ADDR", ":8000"),
		MongoURI:  mustGetenv("MONGO_CON"),
		DB:        getenvOrDefault("MONGO_DB", "wiki"),
		QuotesUrl: getenvOrWarning("QUOTES_URL"),
	}
}

func getenvOrDefault(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getenvOrWarning(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Printf("[WARN] failed to get env var '%s'", key)
	}
	return v
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Errorf("failed to get env var '%s'", key))
	}
	return v
}

func main() {
	web.SetDevMode(true)

	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] No .env file found, using system env")
	}
	conf := LoadConfig()

	qs := quote.Service{BaseUrl: conf.QuotesUrl}
	st, disc := store.NewMongoStore(conf.MongoURI, conf.DB)
	defer disc()
	log.Printf("[INFO] Starting a web server at %s", conf.Addr)
	web.Run(st, conf.Addr, &qs)
}
