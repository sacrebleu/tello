package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	// "time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"

	"tello/route"
	"tello/video"
)

type device struct {
	drone *tello.Driver
	robot *gobot.Robot
}

func main() {

	dev := initDrone()

	registerShutdownHook(dev)

	plan := route.Plan{ Initial: "hello", Name: "ScareSomeCats.fpl" }

	route.Make(&plan)

	err := dev.robot.Start()
	if err != nil {
		fmt.Println("Error", err)
	}
}

// listen for ^C and attempt to issue a graceful shutdown
func registerShutdownHook(dev device) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup(dev.drone)
		os.Exit(1)
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
