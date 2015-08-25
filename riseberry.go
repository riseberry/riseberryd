// Package riseberryd implements an alarm clock that can be controlled via an
// http api.
package riseberryd

import "time"

// Riseberry defines an interfaces which mirrors that of a simple alarm clock.
type Riseberry interface {
	// Get returns the currently set alarm, or an error.
	Get() (Alarm, error)
	// Set sets the given alarm, or returns an error. Calling Set with a Disabled
	// alarm does not stop the alarm if it is already playing.
	Set(Alarm) error
	// Stop stops the alarm if it's playing. Otherwise it does nothing.
	Stop()
	// Close stops all internal goroutines.
	Close()
}

// NewRiseberry returns a new Riseberry instance.
func NewRiseberry(p Player) Riseberry {
	r := &riseberry{
		player: p,
		get:    make(chan chan<- Alarm),
		set:    make(chan Alarm),
		stop:   make(chan struct{}),
		close:  make(chan struct{}),
	}
	go r.loop()
	return r
}

// riseberry implements the Riseberry interface.
type riseberry struct {
	player Player
	get    chan chan<- Alarm
	set    chan Alarm
	stop   chan struct{}
	close  chan struct{}
}

// loop handles riseberry requests.
func (r *riseberry) loop() {
	var (
		alarm Alarm
		play  = make(<-chan time.Time)
	)
	for {
		select {
		case ch := <-r.get:
			ch <- alarm
		case alarm = <-r.set:
			if alarm.Enabled {
				now := time.Now()
				t := alarm.Time(now)
				play = time.After(t.Sub(now))
			} else {
				play = make(<-chan time.Time)
			}
		case <-play:
			alarm.Enabled = false
			go r.player.Play()
		case <-r.stop:
			r.player.Stop()
		case <-r.close:
			r.player.Stop()
			return
		}
	}
}

// Get is part of the Riseberry interface.
func (r *riseberry) Get() (Alarm, error) {
	ch := make(chan Alarm)
	r.get <- ch
	return <-ch, nil
}

// Set is part of the Riseberry interface.
func (r *riseberry) Set(a Alarm) error {
	r.set <- a
	return nil
}

// Stop is part of the Riseberry interface.
func (r *riseberry) Stop() {
	r.stop <- struct{}{}
	return
}

// Close is part of the Riseberry interface.
func (r *riseberry) Close() {
	r.close <- struct{}{}
	return
}

//func Machine(button <-chan struct{}, set <-chan Alarm, get <-chan chan<- Alarm) {
//var (
//playing bool
//alarm  Alarm
//ticker = time.NewTicker(time.Second)
//)
//for {
//select {
//case alarm = <-set:
//case ch := <-get:
//ch <- alarm
//case <-button:
//alarm.Enabled = false
//if playing {
//// stop
//}
//case tick = <-ticker.C:
//if !playing && alarm.Enabled {
//// play
//}
//}
////switch state {
////case Unset:
////select {
////case alarm = <-set:
////case ch := <-get:
////ch <- alarm
////}
////if alarm.Enabled {
////state = Set
////}
////case Set:
////select {
////case <-button:
////alarm.Enabled = false
////case alarm = <-set:
////case tick = <-ticker.C:
////case ch := <-get:
////ch <- alarm
////}
////if alarm.Enabled {
////state = Unset
////}
////case Playing:
//}
//}
//}
