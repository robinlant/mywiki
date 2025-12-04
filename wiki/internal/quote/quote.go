package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Quote struct {
	Text   string   `json:"text"`
	Author string   `json:"author"`
	Tags   []string `json:"tags"`
}

type Service struct {
	BaseUrl string
}

func (qs Service) GetRandomQuote() (*Quote, error) {
	url := qs.BaseUrl + "/random"
	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving random quote from '%s': %s", url, err.Error())

	}
	rData, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading random quote from '%s': %s", url, err.Error())
	}
	var q Quote
	if err := json.Unmarshal(rData, &q); err != nil {
		return nil, fmt.Errorf("error while unmarshaling random quote from '%s': %s", url, err.Error())
	}
	return &q, nil
}
