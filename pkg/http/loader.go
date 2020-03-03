package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

// NewMemStatsLoader create a loader
func NewMemStatsLoader(url string) *memStatsLoader {
	return &memStatsLoader{
		url: url,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// memStatsLoader is a struct implements pkg.MemStatsLoader
type memStatsLoader struct {
	url    string
	client *http.Client
}

func (p *memStatsLoader) Load() (*runtime.MemStats, error) {
	resp, err := p.client.Get(p.url)
	if err != nil {
		return nil, fmt.Errorf("load memstat connect err %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("load memstat bad code err, %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("load memstat read  err %w", err)
	}

	// define an anonymous struct for JSON parsing, then return *runtime.MemStats
	result := &struct {
		Stats *runtime.MemStats `json:"memstats"`
	}{}
	if err := json.Unmarshal(b, result); err != nil {
		return nil, fmt.Errorf("fetch memstat, json  err %w", err)
	}

	return result.Stats, nil
}
