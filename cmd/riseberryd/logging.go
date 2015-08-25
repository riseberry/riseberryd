package main

import (
	"log"

	"github.com/riseberry/riseberryd"
)

func LoggedRiseberry(rb riseberryd.Riseberry) riseberryd.Riseberry {
	return &loggedRiseberry{Riseberry: rb}
}

type loggedRiseberry struct {
	riseberryd.Riseberry
}

func (l *loggedRiseberry) Get() (riseberryd.Alarm, error) {
	alarm, err := l.Riseberry.Get()
	log.Printf("[riseberry] get alarm=%#v err=%s", alarm, err)
	return alarm, err
}

func (l *loggedRiseberry) Set(a riseberryd.Alarm) error {
	err := l.Riseberry.Set(a)
	log.Printf("[riseberry] set alarm=%#v err=%s", a, err)
	return err
}

func (l *loggedRiseberry) Stop() {
	l.Riseberry.Stop()
	log.Printf("[riseberry] stop")
}

func LoggedPlayer(p riseberryd.Player) riseberryd.Player {
	return &loggedPlayer{Player: p}
}

type loggedPlayer struct {
	riseberryd.Player
}

func (l *loggedPlayer) Play() error {
	err := l.Player.Play()
	log.Printf("[player] play err=%s", err)
	return err
}

func (l *loggedPlayer) Stop() {
	l.Player.Stop()
	log.Printf("[player] stop")
}
