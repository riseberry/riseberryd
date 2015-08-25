package riseberryd

import "time"

// Alarm holds an alarm setting.
type Alarm struct {
	// Hour is the hour the alarm is supposed to go off, 0-23.
	Hour int `json:"hour"`
	// Minute is the minute the alarm is supposed to go off, 0-59.
	Minute int `json:"minute"`
	// Zone is the UTC offset in seconds the alarm time is given in.
	Zone int `json:"zone"`
	// Enabled determines if the alarm is enabled or not.
	Enabled bool `json:"enabled"`
}

// Time returns the next time after t when the alarm should go off.
func (a Alarm) Time(t time.Time) time.Time {
	zone := time.FixedZone("", a.Zone)
	when := time.Date(t.Year(), t.Month(), t.Day(), a.Hour, a.Minute, 0, 0, zone)
	if !when.After(t) {
		when = when.AddDate(0, 0, 1)
	}
	return when
}
