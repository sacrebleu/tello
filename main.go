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
	"syscall"
	"time"

	ui "github.com/gizak/termui/v3"

	//"tello/route"
	"tello/video"
	"tello/util"

	//ui "github.com/gizak/termui/v3"
)

// device - a struct to represent the tello drone
type device struct {
	drone *tello.Driver
	robot *gobot.Robot
}

// application - a struct to represent the components of the application
type application struct {
	Ui *screen
	Dev *device
	live bool
	ctrl chan os.Signal
}

// container for UI elements
type screen struct {
	LogArea *widgets.List
	Speed *widgets.Paragraph
}

func (s * screen) displayAirspeed(ground float32, air float32) {
	s.Speed.Text = fmt.Sprintf("Air: %10f m/s\nGnd: %10f m/s", air, ground)
}

var connect = false

var buffer util.Buffer

func main() {
	buffer = util.Buffer.New(util.Buffer{}, 5)

	var dev device

	buffer.Append("Tello Driver 0.1")

	screen := buildUi()

	if connect {
		dev = initDrone()
	}

	// construct app state container
	var app = application{ Ui: screen, Dev: &dev, live: true, ctrl: make(chan os.Signal) }

	// ensure we catch events and handle them
	registerShutdownHook(app)

	//plan := route.NewPlan("ScareSomeCats.fpl")
	//
	//route.Append(plan, &route.Command{Action: "translate20", Description:"Translate 20cm", Offset:5} )
	//route.Append(plan, &route.Command{Action: "translate50", Description:"Translate 50cm", Offset:20} )

	go func() {
		var count = 0
		for app.live {
			//buffer.Append(fmt.Sprintf("%b", app.live))
			//buffer.Append(fmt.Sprintf("count: %d", count))
			time.Sleep(1*time.Second)
			render(app)
			count++
		}
	}()

	//route.Tabulate(plan)
	buffer.Append("Connecting to Drone...")
	if connect {
		err := dev.robot.Start()
		if err != nil {
			log.Fatal("Error", err)
		}
	} else {
		buffer.Append("dummy mode enabled")
		for app.live {
			time.Sleep(1*time.Second)
			//buffer.Append(fmt.Sprintf("%b", app.live))
			//buffer.Append("still alive")
		}
	}

}

// listen for ^C and attempt to issue a graceful shutdown
func registerShutdownHook(app application) {
	signal.Notify(app.ctrl, os.Interrupt, syscall.SIGINT)
	signal.Notify(app.ctrl, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer ui.Close() // only close the tty once we're done
		<-app.ctrl
		buffer.Append("Shutting down connection to drone")
		if connect {
			if app.Dev != nil && app.Dev.drone != nil {
				cleanup(app.Dev.drone)
			}
		}
		app.live = false
		buffer.Append("Waiting for shutdown.")
		time.Sleep(2*time.Second)
		os.Exit(0)
	}()
}

// Connect to the tello drone if it exists
func initDrone() device {

	drone := tello.NewDriver("8890")

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		func() { video.Grab(drone) },
	)

	dev := device{drone: drone, robot: robot}

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
func buildUi() * screen {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	//p := widgets.NewParagraph()
	p := widgets.NewList()
	p.Title = "Flightplan"
	p.Rows = buffer.Values
	p.SetRect(0, 0, 80, 7)
	p.BorderStyle.Fg = ui.ColorWhite

	spd := widgets.NewParagraph()
	spd.SetRect(81, 0, 120, 4)
	spd.BorderStyle.Fg = ui.ColorWhite
	spd.Title = "Drone Speed"

	s := &screen{LogArea: p, Speed: spd}

	return s
}

// render updates to the ui
func render(app application) {
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(200*time.Millisecond).C
	for app.live {
		select {
		case e := <-uiEvents:
			//buffer.Append(fmt.Sprintf("event: %s", e.ID))
			switch e.ID {
			case "q", "<C-c>":
				app.ctrl <- syscall.SIGTERM // invoke cleanup / shutdown handler
			}
		case <-ticker:
			app.Ui.displayAirspeed(rand.Float32() + 1, rand.Float32()+2)
			refreshElements(app.Ui)
		}
	}
}

func refreshElements(screen *screen) {
	updateParagraph(screen)
}

func updateParagraph(screen *screen) {
	screen.LogArea.Rows = buffer.Values
	ui.Render(screen.LogArea, screen.Speed)
}
