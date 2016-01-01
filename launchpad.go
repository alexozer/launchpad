package launchpad

import (
	"errors"
	"strings"
	"sync"

	"github.com/rakyll/portmidi"
)

const (
	bufSize = 1024
)

type Launchpad struct {
	in, out *portmidi.Stream
	events  chan Event
	lights  chan lightRequest

	lastErr      error
	lastErrMutex *sync.Mutex
}

type Event struct {
	X, Y  int
	Press bool
}

type lightRequest struct {
	X, Y int
	G, R int
}

func Open() (*Launchpad, error) {
	input, output, err := discover()
	if err != nil {
		return nil, err
	}

	pad := new(Launchpad)
	if pad.in, err = portmidi.NewInputStream(input, bufSize); err != nil {
		return nil, err
	}
	if pad.out, err = portmidi.NewOutputStream(output, bufSize, 0); err != nil {
		return nil, err
	}

	pad.events = make(chan Event)
	pad.lights = make(chan lightRequest)

	pad.lastErrMutex = new(sync.Mutex)

	go pad.readEvents()
	go pad.readLights()

	return pad, nil
}

func (this *Launchpad) Listen() <-chan Event {
	return this.events
}

func (this *Launchpad) Light(x, y, g, r int) {
	this.lights <- lightRequest{x, y, g, r}
}

func (this *Launchpad) LastError() error {
	this.lastErrMutex.Lock()
	defer this.lastErrMutex.Unlock()

	return this.lastErr
}

func (this *Launchpad) Close() {
	close(this.events)
	close(this.lights)
	this.in.Close()
	this.out.Close()
}

func (this *Launchpad) readEvents() {
	for me := range this.in.Listen() {
		var e Event

		if me.Status == 176 {
			// top row button
			e.X = int(me.Data1 - 104)
			e.Y = -1
		} else {
			e.X = int(me.Data1 % 16)
			e.Y = int(me.Data1 / 16)
		}

		if me.Data2 == 127 {
			e.Press = true
		} else {
			e.Press = false
		}

		this.events <- e
	}
}

func (this *Launchpad) readLights() {
	for l := range this.lights {
		var status, note, vel int64

		if l.Y == -1 {
			// top row button
			status = 176
			note = int64(l.X + 104)
		} else {
			status = 144
			note = int64(l.Y*16 + l.X)
		}
		vel = int64(l.G*16 + l.R + 8 + 4)

		this.logError(this.out.WriteShort(status, note, vel))
	}
}

func discover() (input portmidi.DeviceId, output portmidi.DeviceId, err error) {
	in := -1
	out := -1
	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.GetDeviceInfo(portmidi.DeviceId(i))
		if strings.Contains(info.Name, "Launchpad") {
			if info.IsInputAvailable {
				in = i
			}
			if info.IsOutputAvailable {
				out = i
			}
		}
	}
	if in == -1 || out == -1 {
		err = errors.New("launchpad: no launchpad is connected")
	} else {
		input = portmidi.DeviceId(in)
		output = portmidi.DeviceId(out)
	}
	return
}

func (this *Launchpad) logError(err error) {
	if err != nil {
		this.lastErrMutex.Lock()
		this.lastErr = err
		this.lastErrMutex.Unlock()
	}
}
