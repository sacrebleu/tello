package util

import (
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

// capture meaningful telemetry constructs here
type Telemetry struct {
	Speed *widgets.Paragraph
	Direction *widgets.PieChart
	Yaw *widgets.Paragraph
	Warnings *widgets.Table
}

// container for UI elements
type Screen struct {
	LogArea *widgets.List
	FlightPlan *widgets.List
	Telemetry Telemetry
	Help *widgets.Paragraph
}

// application - a struct to represent the components of the application
type Application struct {
	Ui *Screen
	Dev *Device
	FlightPlan Plan
	Live bool
	Ctrl chan os.Signal
}
