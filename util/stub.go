package util

import (
	"gobot.io/x/gobot/platforms/dji/tello"
	"math/rand"
)

func MockFlightData() tello.FlightData {
	return tello.FlightData{
		BatteryLow: true,
		BatteryPercentage: 10,
		CameraState: 1,
		DroneFlyTimeLeft: 45,
		BatteryLower: true,
		BatteryState: true,
		DownVisualState: false,
		DroneHover: false,
		EastSpeed: int16(-10+ rand.Intn(20)),
		ElectricalMachineryState: 1,
		DroneBatteryLeft: 10,
		EmOpen: false,
		FactoryMode: false,
		Flying: true,
		FlyMode: 1,
		FlyTime: 130,
		FrontIn:false,
		FrontLSC:false,
		FrontOut:false,
		GravityState:true,
		Height: 10,
		ImuCalibrationState: 1,
		ImuState: false,
		LightStrength:1,
		NorthSpeed: int16(-10+ rand.Intn(20)),
		OnGround: false,
		OutageRecording: true,
		PowerState: true,
		PressureState: true,
		SmartVideoExitMode: 1,
		TemperatureHigh: true,
		ThrowFlyTimer: 0,
		VerticalSpeed: int16(-10+ rand.Intn(20)),
		WindState: false }
}