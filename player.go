package riseberryd

import (
	"bytes"
	"fmt"
)

import "os/exec"

// Player defines an interface for playing sound.
type Player interface {
	// Play plays until it's done, or returns an error. Calling Play while
	// another Play is running, implicitly stops the other Play.
	Play() error
	// Stop stops the player and returns when it has stopped.
	Stop()
	// Close stops and closes the Player.
	Close()
}

// NewPlayer returns a new Player that plays the given sound by executing the
// given cmd with the file as the first argument.
func NewPlayer(cmd, file string) Player {
	p := &player{
		cmd:   cmd,
		file:  file,
		play:  make(chan chan error),
		stop:  make(chan chan struct{}),
		close: make(chan chan struct{}),
	}
	go p.loop()
	return p
}

// player implements the Player interface.
type player struct {
	cmd   string
	file  string
	play  chan chan error
	stop  chan chan struct{}
	close chan chan struct{}
}

// loop handles player requests.
func (p *player) loop() {
	var cmd *exec.Cmd
	for {
		select {
		case ch := <-p.play:
			syncKill(cmd)
			cmd = exec.Command(p.cmd, p.file)
			out := bytes.NewBuffer(nil)
			cmd.Stdout = out
			cmd.Stderr = out
			if err := cmd.Start(); err != nil {
				cmd = nil
				ch <- err
				continue
			}
			go func(cmd *exec.Cmd, out fmt.Stringer) {
				err := cmd.Wait()
				if err != nil {
					err = fmt.Errorf("%s: %s", err, out)
				}
				ch <- err
			}(cmd, out)
		case ch := <-p.stop:
			syncKill(cmd)
			ch <- struct{}{}
		case ch := <-p.close:
			syncKill(cmd)
			ch <- struct{}{}
			return
		}
	}
}

// syncKill kills the given cmd if it's running and waits until it exists.
func syncKill(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	cmd.Process.Kill()
	cmd.Wait()
}

// Play is part of the Player interface.
func (p *player) Play() error {
	ch := make(chan error, 1)
	p.play <- ch
	return <-ch
}

// Stop is part of the Player interface.
func (p *player) Stop() {
	ch := make(chan struct{}, 1)
	p.stop <- ch
	<-ch
}

// Close is part of the Player interface.
func (p *player) Close() {
	ch := make(chan struct{}, 1)
	p.close <- ch
	<-ch
}
