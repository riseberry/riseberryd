package riseberryd

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

// Button defines an button interface.
type Button interface {
	// Stop stops the button.
	Stop()
}

// NewButton returns a new Button that samples pin at the given rate and calls
// handler everytime the button is pressed and released.
func NewButton(pin rpio.Pin, rate time.Duration, handler func()) Button {
	b := &button{
		pin:     pin,
		rate:    rate,
		handler: handler,
		stop:    make(chan struct{}),
	}
	go b.loop()
	return b
}

// button implements the Button interface.
type button struct {
	pin     rpio.Pin
	rate    time.Duration
	handler func()
	stop    chan struct{}
}

// loop samples the button and fires the handler.
func (b *button) loop() {
	var (
		prev   = b.pin.Read()
		ticker = time.NewTicker(b.rate)
		n      int
	)
	defer ticker.Stop()
	for {
		state := b.pin.Read()
		if state != prev {
			n++
		}
		prev = state
		if n >= 2 {
			b.handler()
			n = 0
		}
		select {
		case <-ticker.C:
		case <-b.stop:
			return
		}
	}
}

// Stop is part of the Button interface.
func (b *button) Stop() {
	b.stop <- struct{}{}
}
