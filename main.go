package main

import (
	"fmt"
	"github.com/gizak/termui/v3/widgets"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"log"
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

type device struct {
	drone *tello.Driver
	robot *gobot.Robot
}

type application struct {
	Ui *screen
	Dev *device
	live bool
	ctrl chan os.Signal
}

var connect = false

var buffer util.Buffer

func main() {
	buffer = util.Buffer.New(util.Buffer{}, 5)

	var dev device

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
			buffer.Append(fmt.Sprintf("%b", app.live))
			buffer.Append(fmt.Sprintf("count: %d", count))
			time.Sleep(1*time.Second)
			render(app)
			count++
		}
	}()

	//route.Tabulate(plan)

	if connect {
		err := dev.robot.Start()
		if err != nil {
			log.Fatal("Error", err)
		}
	} else {
		for app.live {
			time.Sleep(1*time.Second)
			buffer.Append(fmt.Sprintf("%b", app.live))
			buffer.Append("still alive")
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
		if connect {

			if app.Dev != nil && app.Dev.drone != nil {
				buffer.Append("Shutting down connection to drone")
				cleanup(app.Dev.drone)
			}
		}
		app.live = false
		os.Exit(0)
	}()
}

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

func cleanup(robot *tello.Driver) {
	fmt.Println("interrupt received, cleanup")
	err := robot.Halt()
	if err != nil {
		fmt.Println("Error", err)
	}
}

type screen struct {
	LogArea *widgets.Paragraph
}

func buildUi() * screen {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	p := widgets.NewParagraph()
	p.SetRect(0, 0, 80, 7)

	s := &screen{LogArea: p}

	return s
}

func render(app application) {
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(200*time.Millisecond).C
	for app.live {
		select {
		case e := <-uiEvents:
			buffer.Append(fmt.Sprintf("event: %s", e.ID))
			switch e.ID {
			case "q", "<C-c>":
				app.ctrl <- syscall.SIGTERM // invoke cleanup / shutdown handler
			}
		case <-ticker:
			refreshElements(app.Ui)
		}
	}
}

func refreshElements(screen *screen) {
	updateParagraph(screen)
}

func updateParagraph(screen *screen) {
	screen.LogArea.Text = buffer.Join()
	ui.Render(screen.LogArea)
}
