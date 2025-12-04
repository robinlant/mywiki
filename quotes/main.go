package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const defaultDataPath = "quotes.csv"
const addrVar = "ADDR"
const dataVAr = "DATA"

type Conf struct {
	Addr string
	Data string
}

type Quote struct {
	Text   string
	Author string
	Tags   []string
}

func QuoteFromRow(r []string) *Quote {
	if len(r) != 3 {
		panic(fmt.Errorf("expected an slice of length 3 but got length %d", len(r)))
	}
	rawTags := strings.Split(r[2], ",")
	tagsSlice := make([]string, 0, len(rawTags))
	seen := make(map[string]bool)
	for _, v := range rawTags {
		v = strings.TrimSpace(v)
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		tagsSlice = append(tagsSlice, v)
	}
	return &Quote{
		Text:   r[0],
		Author: r[1],
		Tags:   tagsSlice,
	}
}

func mustLoadConf() Conf {
	a := os.Getenv(addrVar)
	if a == "" {
		log.Panicf("unable to find environmantal variable '%s'", addrVar)
	}
	d := os.Getenv(dataVAr)
	if d == "" {
		d = defaultDataPath
	}
	return Conf{Data: d, Addr: a}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, []*Quote), q []*Quote) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, q)
	}
}

func randomHandler(w http.ResponseWriter, r *http.Request, q []*Quote) {
	randomQuote := q[rand.Intn(len(q))]
	b, err := json.Marshal(randomQuote)
	if err != nil {
		err := fmt.Errorf("error while marshaling into json: %s", err.Error())
		log.Printf("[ERROR] %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Print("[INFO] Serving a random quote")
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	godotenv.Load()
	conf := mustLoadConf()

	f, err := os.Open(conf.Data)
	if err != nil {
		log.Panicf("error while opening file '%v': %s", conf.Data, err.Error())
	}
	defer f.Close()
	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Panicf("eror while parsing csv file '%v': %s", conf.Data, err.Error())
	}
	quotes := make([]*Quote, len(data)-1)
	for i := 1; i < len(data); i++ {
		quotes[i-1] = QuoteFromRow(data[i])
	}

	http.HandleFunc("/random", makeHandler(randomHandler, quotes))

	log.Printf("[INFO] Starting a web server at %s", conf.Addr)
	log.Fatal(http.ListenAndServe(conf.Addr, nil))
}
