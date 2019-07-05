package util

import (
	"fmt"
	"github.com/gizak/termui/v3/widgets"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"os"
)

// device - a struct to represent the tello drone
type Device struct {
	Drone *tello.Driver
	Robot *gobot.Robot
}

// container for UI elements
type Screen struct {
	LogArea *widgets.List
	Speed *widgets.Paragraph
	FlightPlan *widgets.List
	Help *widgets.Paragraph
}

func (s * Screen) DisplayAirspeed(ground float32, air float32) {
	s.Speed.Text = fmt.Sprintf("Air: %10f m/s\nGnd: %10f m/s", air, ground)
}

// application - a struct to represent the components of the application
type Application struct {
	Ui *Screen
	Dev *Device
	FlightPlan [] string
	Live bool
	Ctrl chan os.Signal
}
