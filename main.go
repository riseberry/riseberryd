package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"
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
		sound  = flag.String("sound", "say you did it", "Command for playing alarm sound")
	)
	flag.Parse()
	clock := NewClock(NewCmdSound(*sound))
	handler := http.FileServer(http.Dir(*assets))
	handler = NewAlarmHandler(clock, handler)
	return http.ListenAndServe(*addr, handler)
}

func NewClock(sound Sound) Clock {
	return &clock{sound: sound}
}

type clock struct {
	lock  sync.Mutex
	sound Sound
	alarm Alarm
	timer *time.Timer
}

func (c *clock) Get() Alarm {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.alarm
}

func (c *clock) Set(alarm Alarm) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	if alarm.Enabled {
		now := time.Now()
		when := time.Date(now.Year(), now.Month(), now.Day(), alarm.Hour, alarm.Minute, 0, 0, time.Local)
		if when.Before(now) {
			when = when.AddDate(0, 0, 1)
		}
		c.timer = time.AfterFunc(when.Sub(now), func() {
			c.sound.Play()
		})
	}
	c.alarm = alarm
}

type Clock interface {
	Get() Alarm
	Set(Alarm)
}

type Alarm struct {
	Hour    int  `json:"hour"`
	Minute  int  `json:"minute"`
	Enabled bool `json:"enabled"`
}

func NewAlarmHandler(clock Clock, next http.Handler) http.Handler {
	const prefix = "/alarm"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != prefix {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			var err error
			if r.Method == "PUT" {
				var alarm Alarm
				if err = json.NewDecoder(r.Body).Decode(&alarm); err == nil {
					clock.Set(alarm)
				}
			}
			var response interface{}
			if err == nil {
				response = clock.Get()
			} else {
				response = map[string]interface{}{"error": err.Error()}
			}
			json.NewEncoder(w).Encode(response)
		}
	})
}

func NewCmdSound(text string) Sound {
	return cmd(text)
}

type cmd string

func (c cmd) Play() error {
	return exec.Command("/bin/sh", "-c", string(c)).Run()
}

type Sound interface {
	Play() error
}
