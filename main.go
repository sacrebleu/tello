package main

import (
	"flag"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"log"
	"math"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"tello/util"
	//"tello/route"
	"tello/video"

	//ui "github.com/gizak/termui/v3"
)

var dryrun = false

var buffer util.Buffer

var gp = os.Getenv("GOPATH")
var ap = path.Join(gp, "src/tello/cfg")

var flightPath = ap

// this is the main entry point into the application.
func main() {
	parseCli()

	buffer = util.Buffer.New(util.Buffer{}, 5)

	var dev util.Device

	buffer.Append("Tello Driver 0.1")

	screen := buildUi()

	if !dryrun {
		dev = initDrone(*screen)
	}

	// construct app state container
	var app = util.Application{ Ui: screen, Dev: &dev, Live: true, Ctrl: make(chan os.Signal) }

	registerShutdownHook(app)

	app.ReadPlan(flightPath)

	// ensure we catch events and handle them
	go func() {
		var count = 0
		for app.Live {
			time.Sleep(1*time.Second)
			render(&app)
			count++
		}
	}()

	//route.Tabulate(plan)
	buffer.Append("Connecting to Drone...")
	if dryrun {
		buffer.Append("dummy mode enabled")
		ticker := time.NewTicker(time.Second).C
		for app.Live {
			select {
			case <-ticker:
				updateTelemetry(util.MockFlightData(), screen.Telemetry)
			}
		}
	} else {
		err := dev.Robot.Start()

		if err != nil {
			log.Fatal("Error", err)
		}
	}

}

// consume any command line arguments
func parseCli() {
	flag.StringVar(&flightPath, "path", "./cfg/scare_the_cats.fp", "Specify the path to the flightplan you want to load")
	flag.BoolVar(&dryrun, "dryrun", false, "Enable/Disable drone connection for development [default enabled]")
	flag.Parse()
}

// listen for SIGINT + SIGTERM and attempt to issue a graceful shutdown
func registerShutdownHook(app util.Application) {
	signal.Notify(app.Ctrl, os.Interrupt, syscall.SIGINT)
	signal.Notify(app.Ctrl, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer ui.Close() // only close the tty once we're done
		<-app.Ctrl
		buffer.Append("Shutting down connection to drone")
		if !dryrun {
			if app.Dev != nil && app.Dev.Drone != nil {
				buffer.Append("Cleaning up drone connections...")
				cleanup(app.Dev.Drone)
			}
		}
		app.Live = false
		buffer.Append("Waiting for shutdown.")
		time.Sleep(2*time.Second)
		os.Exit(0)
	}()
}

// Connect to the tello drone if it exists
func initDrone(screen util.Screen) util.Device {

	drone := tello.NewDriver("8890")

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		func() { video.Grab(drone) },
	)

	drone.On(tello.FlightDataEvent, func(data interface{}) {
		df := data.(*tello.FlightData)
		updateTelemetry(df, screen.Telemetry)
	})

	return util.Device{ Drone: drone, Robot: robot}
}

// clean shutdown of the Tello drone
func cleanup(robot *tello.Driver) {
	buffer.Append("interrupt received, cleanup")

	err := robot.Halt()
	if err != nil {
		fmt.Println("Error", err)
	}
}

// construct and lay out the ui
func buildUi() * util.Screen {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	// border
	hp := widgets.NewParagraph()
	hp.Title = "DJI Tello Telemetry"
	//hp.Text
	hp.SetRect(0, 0, 140, 32)
	hp.BorderStyle.Fg = ui.ColorWhite

	hlp := widgets.NewParagraph()
	hlp.Title = " q - quit   j - scroll down   k - scroll up "
	hlp.SetRect(0, 32, 140, 32)

	//p := widgets.NewParagraph()
	p := widgets.NewList()
	p.Title = "Telemetry"
	p.Rows = buffer.Values
	p.SetRect(1,  1, 110, 8)
	p.BorderStyle.Fg = ui.ColorWhite

	fp := widgets.NewList()
	fp.Title = "Flight Plan"
	fp.TextStyle = ui.NewStyle(ui.ColorCyan)
	fp.WrapText = false
	p.BorderStyle.Fg = ui.ColorWhite
	fp.SetRect(1, 8, 110, 16)

	s := &util.Screen{
		LogArea: p,
		FlightPlan: fp,
		Help : hp,
	}

	spd := widgets.NewParagraph()
	spd.SetRect(111, 1, 139, 8)
	spd.BorderStyle.Fg = ui.ColorWhite
	spd.Title = "Speed"

	disp := util.NewCompass(111, 8, 139, 24)

	table1 := widgets.NewTable()
	table1.TextStyle = ui.NewStyle(ui.ColorWhite)
	table1.SetRect(111, 24, 139, 31)
	table1.Rows = [][]string{
		{ "", "", ""},
		{ "", "", ""},
		{ "", "", ""},
	}
	//ui.Render(table1)

	s.Telemetry = util.Telemetry{Speed: spd, Direction: disp, Warnings: table1}

	ui.Render(hp, hlp)

	return s
}

// render updates to the ui
func render(app * util.Application) {
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(200*time.Millisecond).C
	for app.Live {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				app.Ctrl <- syscall.SIGTERM // invoke cleanup / shutdown handler
			case "j":
				app.Ui.FlightPlan.ScrollDown()
			case "k":
				app.Ui.FlightPlan.ScrollUp()
			}
		case <-ticker:
			app.FlightPlan.Next()
			refreshElements(app)
		}
	}
}

func refreshElements(app * util.Application) {
	app.Ui.LogArea.Rows = buffer.Values
	app.Ui.FlightPlan.Rows = app.FlightPlan.Render()
	ui.Render(app.Ui.LogArea, app.Ui.FlightPlan)
	ui.Render(app.Ui.Telemetry.Direction, app.Ui.Telemetry.Speed, app.Ui.Telemetry.Warnings)
}

func updateTelemetry(data * tello.FlightData, telemetry util.Telemetry) {
	angle := math.Atan2(float64(data.NorthSpeed), float64(data.EastSpeed))
	bearing := angle * 180/math.Pi

	var degs = int16(90 - bearing + 360) % 360 // radians

	telemetry.Speed.Text = fmt.Sprintf("Airspeed: %f\nGroundSpeed: %f\nVertical: %d", data.AirSpeed(), data.GroundSpeed(), data.VerticalSpeed)
	telemetry.Direction.Title = fmt.Sprintf("Bearing %d" , degs) // convert to degrees
	telemetry.Direction.AngleOffset = -angle - 0.1*math.Pi

	telemetry.Warnings.Rows = [][]string{
		{		format(data.TemperatureHigh, "Temp", "red", "white"),
			    format(data.BatteryLow, "Batt", "red", "white"),
			    format(data.BatteryLower, "!BAT", "red", "white")},
		{		format(data.DroneHover, "Hvr", "green", "white"),
				format(data.Flying, "Fly", "green", "white"),
				format(data.FrontIn, "FrI", "green", "white")},
		{		format(data.PressureState, "Prs", "green", "white"),
				format(data.PowerState, "Pwr", "green", "white"),
				format(data.ImuState, "Imu", "green", "white")},
		}
}

func format(val bool, title string, onColour string, offColour string) (string) {
	var ret string
	if val {
		ret = fmt.Sprintf("[%s](fg:%s)", title, onColour)
	} else {
		ret = fmt.Sprintf("[%s](fg:%s)", title, offColour)
	}
	return ret
}
