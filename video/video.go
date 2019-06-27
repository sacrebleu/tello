package video

import (
	"fmt"
	"os/exec"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

// need some sort of system state mutex to track when video needs to be shut down
// to avoid broken pipe bollocks

func Grab(drone *tello.Driver) {
	{
		if drone == nil {
			fmt.Println("Error - drone is not initialised")
			return
		}

		mplayer := exec.Command("mplayer", "-fps", "25", "-")
		mplayerIn, _ := mplayer.StdinPipe()
		if err := mplayer.Start(); err != nil {
			fmt.Println(err)
			return
		}

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Video capture connected")
			drone.StartVideo()
			drone.SetVideoEncoderRate(3)
			gobot.Every(100*time.Millisecond, func() {
				drone.StartVideo()
				//fmt.Println(drone.Rate())
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := mplayerIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})
	}
}
