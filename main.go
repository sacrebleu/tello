package main

import (
	"fmt"
	"github.com/gizak/termui/v3/widgets"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
	"flag"

	ui "github.com/gizak/termui/v3"

	//"tello/route"
	"tello/video"
	"tello/util"

	//ui "github.com/gizak/termui/v3"
)

var connect = true

var buffer util.Buffer

var gp = os.Getenv("GOPATH")
var ap = path.Join(gp, "src/tello/cfg")

var flightPath = ap

func main() {
	flag.StringVar(&flightPath, "path", "path/to/flightplan.fp", "Specify the path to the flightplan you want to load")
	flag.Parse()

	buffer = util.Buffer.New(util.Buffer{}, 5)

	var dev util.Device

	buffer.Append("Tello Driver 0.1")

	screen := buildUi()

	if connect {
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
	if connect {
		err := dev.Robot.Start()
		if err != nil {
			log.Fatal("Error", err)
		}
	} else {
		buffer.Append("dummy mode enabled")
		for app.Live {
			time.Sleep(1*time.Second)
			//buffer.Append(fmt.Sprintf("%b", app.Live))
			//buffer.Append("still aLive")
		}
	}

}

// listen for ^C and attempt to issue a graceful shutdown
func registerShutdownHook(app util.Application) {
	signal.Notify(app.Ctrl, os.Interrupt, syscall.SIGINT)
	signal.Notify(app.Ctrl, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer ui.Close() // only close the tty once we're done
		<-app.Ctrl
		buffer.Append("Shutting down connection to drone")
		if connect {
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

	dev := util.Device{ Drone: drone, Robot: robot}

	return dev
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
	hp.SetRect(0, 0, 120, 30)
	hp.BorderStyle.Fg = ui.ColorWhite

	hlp := widgets.NewParagraph()
	hlp.Title = " q - quit   j - scroll down   k - scroll up "
	hlp.SetRect(0, 29, 120, 30)

	//p := widgets.NewParagraph()
	p := widgets.NewList()
	p.Title = "Telemetry"
	p.Rows = buffer.Values
	p.SetRect(1,  1, 80, 8)
	p.BorderStyle.Fg = ui.ColorWhite

	spd := widgets.NewParagraph()
	spd.SetRect(81, 1, 119, 5)
	spd.BorderStyle.Fg = ui.ColorWhite
	spd.Title = "Drone Speed"

	fp := widgets.NewList()
	fp.Title = "Flight Plan"
	fp.SetRect(1, 8, 80, 16)
	p.BorderStyle.Fg = ui.ColorWhite

	s := &util.Screen{
		LogArea: p,
		Speed: spd,
		FlightPlan: fp,
		Help : hp,
	}

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
			//buffer.Append(fmt.Sprintf("event: %s", e.ID))
			switch e.ID {
			case "q", "<C-c>":
				app.Ctrl <- syscall.SIGTERM // invoke cleanup / shutdown handler
			case "j":
				app.Ui.FlightPlan.ScrollDown()
			case "k":
				app.Ui.FlightPlan.ScrollUp()
			}
		case <-ticker:
			app.Ui.DisplayAirspeed(rand.Float32() + 1, rand.Float32()+2)
			refreshElements(app)
		}
	}
}

func refreshElements(app * util.Application) {
	app.Ui.LogArea.Rows = buffer.Values
	app.Ui.FlightPlan.Rows = app.FlightPlan
	ui.Render(app.Ui.LogArea, app.Ui.Speed, app.Ui.FlightPlan)
}
