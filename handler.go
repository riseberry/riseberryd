package riseberryd

import (
	"encoding/json"
	"errors"
)
import "net/http"

// Handler returns a handler that handles /alarm requests, or calls the next
// handler.
func Handler(rb Riseberry, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/alarm" {
			next.ServeHTTP(w, r)
			return
		}
		var (
			err   error
			alarm Alarm
			code  = http.StatusInternalServerError
		)
		switch r.Method {
		case "GET":
			alarm, err = rb.Get()
		case "PUT":
			if err = json.NewDecoder(r.Body).Decode(&alarm); err == nil {
				err = rb.Set(alarm)
			} else {
				code = http.StatusBadRequest
			}
		default:
			code = http.StatusMethodNotAllowed
			err = errors.New(http.StatusText(code))
			w.Header().Set("Allow", "GET, PUT")
		}
		if err != nil {
			writeJSONError(w, code, err)
		} else {
			writeJSON(w, alarm)
		}
	})
}

// writeJSON writes the given data as JSON to the given w or returns an error.
func writeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(w)
	return e.Encode(data)
}

// writeJSONError writes err with the given code to w or returns an error.
func writeJSONError(w http.ResponseWriter, code int, err error) error {
	w.WriteHeader(code)
	return writeJSON(w, struct {
		Error string `json:"error"`
	}{Error: err.Error()})
}
