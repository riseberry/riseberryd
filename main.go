package main

import (
	"time"
)

// State holds the states the process can be in.
type State int

const (
	// Unset is the active state when no alarm is enabled.
	Unset State = iota
	// Set is the active state when an alarm is enabled, but is not playing yet.
	Set
	// Playing is the active state when an alarm is playing.
	Playing
)

func Machine(button <-chan struct{}, set <-chan Alarm, get <-chan chan<- Alarm) {
	var (
		state  State
		ticker = time.NewTicker(time.Second)
		alarm  Alarm
	)
	for {
		switch state {
		case Unset:
			select {
			case alarm = <-set:
			case ch := <-get:
				ch <- alarm
			}
			if alarm.Enabled {
				state = Set
			}
		case Set:
			select {
			case <-button:
				alarm.Enabled = false
			case alarm = <-set:
			case tick = <-ticker.C:
			case ch := <-get:
				ch <- alarm
			}
			if alarm.Enabled {
				state = Unset
			}
		case Playing:
		}
	}
}

// Alarm holds the settings that make up an alarm.
type Alarm struct {
	Hour    int  `json:"hour"`
	Minute  int  `json:"minute"`
	Zone    int  `json:"zone"`
	Enabled bool `json:"enabled"`
}
