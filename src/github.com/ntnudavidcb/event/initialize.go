package event

import (
	"../com"
	. "../config"
	"../driver"
	"../io"
	"../queue"
)

func InitElevator(buttonPressed chan int, floorReached chan bool) {
	floor, _ := driver.Elev_init()
	io.SetElevState(floor, DIR_STOP, NOT_ANY_BUTTON)
	initMyIP()
	initElevStateMap(floor)
	io.InitElevState(floor)
	queue.InitQueue()
	io.RunBackup()
	queue.SynchronizeQueueWithIO(io.GetPressedButtons())
}

func initMyIP() {
	IPAddr := com.GetMyIP()
	queue.SetMyIP(IPAddr)
}

func initElevStateMap(floor int) {
	IPAddr := com.GetMyIP()
	queue.UpdateElevStateMap(IPAddr, floor, DIR_STOP, NOT_ANY_BUTTON)
}
