package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	// "time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"

	"sacrebleu/tello/video"
)

func main() {
	drone := tello.NewDriver("8890")

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		func() { video.Grab(drone) },
	)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup(drone)
		os.Exit(1)
	}()

	err := robot.Start()
	if err != nil {
		fmt.Println("Error", err)
	}
}

func cleanup(robot *tello.Driver) {
	fmt.Println("interrupt received, cleanup")
	err := robot.Halt()
	if err != nil {
		fmt.Println("Error", err)
	}
}
