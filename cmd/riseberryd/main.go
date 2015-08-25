package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/riseberry/riseberryd"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		httpAddr  = flag.String("http.addr", ":80", "HTTP addr to listen on")
		httpDir   = flag.String("http.dir", "./www", "Dir to serve http assets from")
		soundCmd  = flag.String("sound.cmd", "omxplayer", "Command to use for playing sound.")
		soundFile = flag.String("sound.file", "", "Path to sound file to play")
	)
	flag.Parse()
	log.Printf("[risberryd] started")
	player := riseberryd.NewPlayer(*soundCmd, *soundFile)
	player = LoggedPlayer(player)
	rb := riseberryd.NewRiseberry(player)
	rb = LoggedRiseberry(rb)
	defer rb.Close()
	handler := http.FileServer(http.Dir(*httpDir))
	handler = riseberryd.Handler(rb, handler)
	return http.ListenAndServe(*httpAddr, handler)
}
