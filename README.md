# launchpad

A Go library for communicating with the Novation Launchpad originally based on [this](https://github.com/rakyll/launchpad). Improvements include complete thread safety, support for button release events, and support for the top button row.

```sh
go get github.com/alexozer/launchpad
```

## Usage

### Initialization

```go
if pad, err = launchpad.Open(); err != nil {
    log.Error("error while initializing launchpad")
}
defer pad.Close()
```
If there are no currently connected Launchpads, initialization will fail with an error. You can fake a device by creating an input and output MIDI device with the name `Launchpad`.

### Coordinate system

```
+--------- arrow keys -----------+  +--- mode keys ---+
{0,-1} {1,-1} {2,-1} {3,-1} {4,-1} {5,-1} {6,-1} {7,-1} | ableton
----------------------------------------------------------------
{0, 0} {1, 0} {2, 0} {3, 0} {4, 0} {5, 0} {6, 0} {7, 0} | {8, 0} vol
----------------------------------------------------------------
{0, 1} {1, 1} {2, 1} {3, 1} {4, 1} {5, 1} {6, 1} {7, 1} | {8, 1} pan
----------------------------------------------------------------
{0, 2} {1, 2} {2, 2} {3, 2} {4, 2} {5, 2} {6, 2} {7, 2} | {8, 2} sndA
----------------------------------------------------------------
{0, 3} {1, 3} {2, 3} {3, 3} {4, 3} {5, 3} {6, 3} {7, 3} | {8, 3} sndB
----------------------------------------------------------------
{0, 4} {1, 4} {2, 4} {3, 4} {4, 4} {5, 4} {6, 4} {7, 4} | {8, 4} stop
----------------------------------------------------------------
{0, 5} {1, 5} {2, 5} {3, 5} {4, 5} {5, 5} {6, 5} {7, 5} | {8, 5} trk on
----------------------------------------------------------------
{0, 6} {1, 6} {2, 6} {3, 6} {4, 6} {5, 6} {6, 6} {7, 6} | {8, 6} solo
----------------------------------------------------------------
{0, 7} {1, 7} {2, 7} {3, 7} {4, 7} {5, 7} {6, 7} {7, 7} | {8, 7} arm
----------------------------------------------------------------
```

### Light buttons

```go
pad.Light(0, 0, 3, 0) // lights the bottom left button with bright green
```

The g and r components of the color are from 0-3 inclusive.

### Read events
```go
for e := range pad.Listen() {
	fmt.Println(e.X, e.Y, e.Pressed)
}
```

## Demos

Some demos are available in the [demos](/demos) directory.
