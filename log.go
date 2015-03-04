package main

import (
	"log"
	"net/http"
)

func LogStart(args []string, level int) {
	if level > 0 {
		log.Printf("started with args: %#v", args)
	}
}

func LoggedHandler(handler http.Handler, level int) http.Handler {
	if level <= 0 {
		return handler
	}
	return &loggedHandler{handler: handler}
}

type loggedHandler struct {
	handler http.Handler
}

func (h *loggedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := LoggedResponseWriter(w)
	h.handler.ServeHTTP(lw, r)
	log.Printf("http: %d %s %s", lw.Code, r.Method, r.URL.Path)
}

func LoggedResponseWriter(w http.ResponseWriter) *loggedResponseWriter {
	return &loggedResponseWriter{ResponseWriter: w, Code: http.StatusOK}
}

type loggedResponseWriter struct {
	http.ResponseWriter
	Code int
}

func (l *loggedResponseWriter) WriteHead(code int) {
	l.Code = code
	l.ResponseWriter.WriteHeader(code)
}

func LoggedClock(clock Clock, level int) Clock {
	if level <= 0 {
		return clock
	}
	return &loggedClock{clock: clock, level: level}
}

type loggedClock struct {
	clock Clock
	level int
}

func (c *loggedClock) Get() Alarm {
	alarm := c.clock.Get()
	if c.level > 1 {
		log.Printf("clock: get alarm=%#v", alarm)
	}
	return alarm
}

func (c *loggedClock) Set(alarm Alarm) {
	c.clock.Set(alarm)
	log.Printf("clock: set alarm=%#v", alarm)
}
