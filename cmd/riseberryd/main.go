package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/riseberry/riseberryd"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		httpAddr   = flag.String("http.addr", ":80", "HTTP addr to listen on")
		httpDir    = flag.String("http.dir", "./www", "Dir to serve http assets from")
		soundCmd   = flag.String("sound.cmd", "omxplayer", "Command to use for playing sound.")
		soundFile  = flag.String("sound.file", "", "Path to sound file to play")
		buttonPin  = flag.Int("button.pin", 17, "Pin the button is connected to")
		buttonRate = flag.Duration("button.rate", 5*time.Millisecond, "Rate for sampling the button")
	)
	flag.Parse()
	log.Printf("[risberryd] started")
	player := riseberryd.NewPlayer(*soundCmd, *soundFile)
	player = LoggedPlayer(player)
	if err := rpio.Open(); err != nil {
		return err
	}
	defer rpio.Close()
	rb := riseberryd.NewRiseberry(player)
	rb = LoggedRiseberry(rb)
	defer rb.Close()
	button := riseberryd.NewButton(rpio.Pin(*buttonPin), *buttonRate, rb.Stop)
	defer button.Stop()
	handler := http.FileServer(http.Dir(*httpDir))
	handler = riseberryd.Handler(rb, handler)
	return http.ListenAndServe(*httpAddr, handler)
}
