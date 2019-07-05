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
		dev = initDrone()
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
		events := dev.Drone.Subscribe()

		for app.Live {
			ev, _ := <-events
			fmt.Println(ev)
			if ev.Name == "FlightDataEvent" {
				fd, ok := ev.Data.(tello.FlightData)
				if ok {
					updateTelemetry(fd, screen.Telemetry)
				}
			}
		}
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
func initDrone() util.Device {

	drone := tello.NewDriver("8890")

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		func() { video.Grab(drone) },
	)

	return util.Device{ Drone: drone, Robot: robot}
}

// clean shutdown of the Tello drone
func cleanup(robot *tello.Driver) {
	fmt.Println("interrupt received, cleanup")
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
	hp.SetRect(0, 0, 140, 30)
	hp.BorderStyle.Fg = ui.ColorWhite

	hlp := widgets.NewParagraph()
	hlp.Title = " q - quit   j - scroll down   k - scroll up "
	hlp.SetRect(0, 29, 140, 30)

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

	s.Telemetry = util.Telemetry{Speed: spd, Direction: disp}

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
	ui.Render(app.Ui.Telemetry.Direction, app.Ui.Telemetry.Speed)
}

func updateTelemetry(data tello.FlightData, telemetry util.Telemetry) {
	angle := math.Atan2(float64(data.NorthSpeed), float64(data.EastSpeed))
	bearing := angle * 180/math.Pi

	var degs = int16(90 - bearing + 360) % 360 // radians

	telemetry.Speed.Text = fmt.Sprintf("Airspeed: %f\nGroundSpeed: %f\nVertical: %d", data.AirSpeed(), data.GroundSpeed(), data.VerticalSpeed)
	telemetry.Direction.Title = fmt.Sprintf("Bearing %d" , degs) // convert to degrees
	telemetry.Direction.AngleOffset = -angle
}
