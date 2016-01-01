package main

import (
	"log"
	"time"

	"github.com/alexozer/launchpad"
)

var pad *launchpad.Launchpad

const fadeDelay = time.Millisecond * 60

func main() {
	var err error
	if pad, err = launchpad.Open(); err != nil {
		log.Fatal(err)
	}
	defer pad.Close()

	for e := range pad.Listen() {
		if e.Press {
			go fadeOn(e.X, e.Y)
		} else {
			go fadeOff(e.X, e.Y)
		}
	}
}

func fadeOn(x, y int) {
	for i := 1; i <= 3; i++ {
		pad.Light(x, y, i, 0)
		time.Sleep(fadeDelay)
	}
}

func fadeOff(x, y int) {
	for i := 2; i >= 0; i-- {
		pad.Light(x, y, i, 0)
		time.Sleep(fadeDelay)
	}
}
