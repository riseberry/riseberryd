package main

import (
	"bytes"
	"os/exec"
)

import "text/template"

// NewCmdPlayer returns a new player that plays a sound by executing the given
// shell command which is expected to contain a {{.File}} placeholder.
func NewCmdPlayer(cmd string) (Player, error) {
	if tmp, err := template.New("cmd").Parse(cmd); err != nil {
		return nil, err
	} else {
		return &cmdPlayer{tmp: tmp}
	}
}

// Player defines an interface for playing sound files.
type Player interface {
	// Play plays the given sound file and returns when it is done.
	Play(file string) error
	// Stop stops the currently playing sound (if any) and causes Play to return.
	Stop() error
}

// cmdPlayer implements the layer interface.
type cmdPlayer struct {
	tmp *template.Template
	cmd *exec.Cmd
}

// Play plays the given
func (p *cmdPlayer) Play(file string) error {
	buf := bytes.NewBuffer(nil)
	if err := p.tmp.Execute(buf, struct {
		File string
	}{File: file}); err != nil {
		return err
	}
	p.cmd = exec.Command("/bin/sh", "-c", buf.String())
	return p.cmd.Run()
}

//NewCmdSound returns a sound that is played by invoking the given shell cmd.
//func NewCmdSound(cmd string) Sound {
//return cmdSound(cmd)
//}

//cmdSound implements the Sound interface.
//type cmdSound string

//Play is part of the Sound interface.
//func (c cmdSound) Play() error {
//return exec.Command("/bin/sh", "-c", string(c)).Run()
//}
