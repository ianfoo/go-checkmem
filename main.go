package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// maxErrCount is the maximum number of errors before the program quits.
const maxErrCount = 10

var (
	addr         = "localhost:6060"
	tickInterval = 30 * time.Second
	errCount     = 0
)

// memstats is subset of a runtime.MemStats
type memstats struct {
	Timestamp     string
	TimestampUnix int64
	Alloc         uint64
	TotalAlloc    uint64
	Sys           uint64
	HeapInUse     uint64
	HeapAlloc     uint64
	HeapIdle      uint64
	HeapReleased  uint64
	HeapObjects   uint64
	HeapSys       uint64
}

func NewFromMemStats(ms map[string]interface{}) memstats {
	ts := time.Now()
	return memstats{
		Timestamp:     ts.Format(time.RFC3339),
		TimestampUnix: ts.Unix(),
		Alloc:         uint64(ms["Alloc"].(float64)),
		TotalAlloc:    uint64(ms["TotalAlloc"].(float64)),
		Sys:           uint64(ms["Sys"].(float64)),
		HeapInUse:     uint64(ms["HeapInuse"].(float64)),
		HeapAlloc:     uint64(ms["HeapAlloc"].(float64)),
		HeapIdle:      uint64(ms["HeapIdle"].(float64)),
		HeapReleased:  uint64(ms["HeapReleased"].(float64)),
		HeapObjects:   uint64(ms["HeapObjects"].(float64)),
		HeapSys:       uint64(ms["HeapSys"].(float64)),
	}
}

func init() {
	if envAddr := os.Getenv("ADDR"); envAddr != "" {
		addr = envAddr
	}
	if envIntv := os.Getenv("INTERVAL"); envIntv != "" {
		intv, err := strconv.Atoi(envIntv)
		if err == nil {
			tickInterval = time.Duration(intv) * time.Second
		}
	}
}

func main() {
	url := "http://" + addr + "/debug/vars"
	log.Printf("checking expvars at %s every %s", url, tickInterval)
	tick := time.NewTicker(tickInterval)
	for range tick.C {
		if errCount >= maxErrCount {
			log.Fatalf("maximum error threshold of %d breached; exiting", maxErrCount)
		}
		resp, err := http.Get(url)
		if err != nil {
			errCount++
			log.Printf("error fetching expvars: %v", err)
			continue
		}
		var body map[string]interface{}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&body); err != nil {
			errCount++
			log.Printf("error decoding response from expvar check: %v", err)
			continue
		}
		if _, ok := body["memstats"]; !ok {
			errCount++
			log.Println("error: expected \"memstats\" in response, but missing")
			continue
		}
		mss := NewFromMemStats(body["memstats"].(map[string]interface{}))
		formatted, err := json.Marshal(mss)
		if err != nil {
			errCount++
			log.Printf("error encoding memory stats as JSON: %v", err)
			continue
		}
		fmt.Println(string(formatted))
	}
}
