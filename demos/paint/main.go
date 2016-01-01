package main

import (
	"log"
	"sync"
	"time"

	"github.com/alexozer/launchpad"
)

type Canvas struct {
	pad *launchpad.Launchpad

	selectedColor int
	mutex         sync.Mutex
}

type Color struct {
	G, R int
}

var palette = []Color{{0, 3}, {1, 3}, {2, 3}, {3, 3}, {3, 2}, {3, 1}, {3, 0}, {0, 0}}

const flashDelay = time.Millisecond * 250

func NewCanvas() (canvas *Canvas, err error) {
	canvas = new(Canvas)
	if canvas.pad, err = launchpad.Open(); err != nil {
		return nil, err
	}

	for x, color := range palette {
		canvas.pad.Light(x, -1, color.G, color.R)
	}

	go canvas.flashSelectedColor()
	canvas.paint()

	return
}

func (this *Canvas) Close() {
	this.pad.Close()
}

func (this *Canvas) getSelectedColor() int {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.selectedColor
}

func (this *Canvas) setSelectedColor(color int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.selectedColor = color
}

func (this *Canvas) flashSelectedColor() {
	for {
		x := this.getSelectedColor()
		color := palette[x]

		this.pad.Light(x, -1, 0, 0)
		time.Sleep(flashDelay)
		this.pad.Light(x, -1, color.G, color.R)
		time.Sleep(flashDelay)
	}
}

func (this *Canvas) paint() {
	for e := range this.pad.Listen() {
		if e.Y == -1 {
			this.setSelectedColor(e.X)
		} else {
			color := palette[this.getSelectedColor()]
			this.pad.Light(e.X, e.Y, color.G, color.R)
		}
	}
}

func main() {
	_, err := NewCanvas()
	if err != nil {
		log.Fatal(err)
	}
}
