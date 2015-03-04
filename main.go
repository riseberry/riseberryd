package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		addr   = flag.String("addr", ":80", "HTTP addr to listen on")
		assets = flag.String("assets", "./public", "Directory to serve assets from")
	)
	flag.Parse()
	handler := http.FileServer(http.Dir(*assets))
	handler = NewAlarmHandler(handler)
	return http.ListenAndServe(*addr, handler)
}

type Alarm struct {
	Hour    int  `json:"hour"`
	Minute  int  `json:"minute"`
	Enabled bool `json:"enabled"`
}

func NewAlarmHandler(next http.Handler) http.Handler {
	const prefix = "/alarm"
	alarm := Alarm{Enabled: true, Hour: 23, Minute: 15}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != prefix {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			response := interface{}(&alarm)
			if r.Method == "PUT" {
				if err := json.NewDecoder(r.Body).Decode(&alarm); err != nil {
					response = map[string]interface{}{"error": err.Error()}
				}
			}
			json.NewEncoder(w).Encode(response)
		}
	})

}

func NewSaySound(text string) Sound {
	return say(text)
}

type say string

func (a say) Play() error {
	return exec.Command("say", string(a)).Run()
}

type Sound interface {
	Play() error
}
