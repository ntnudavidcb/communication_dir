package main

import (
	"../config"
	"../driver"
	"../io"
	"../queue"
	"log"
	//"time"
)

func eventButtonPushed(buttonPushed int) {
	queue.UpdateQueueWithButton(buttonPushed)
}

func eventFloorReached() {
	log.Println("CheckOrder: ", queue.CheckOrder())
	if queue.CheckOrder() {
		log.Println("It should have stopped here")
		io.WantedFloorReached()
		queue.RemoveFromQueue(io.GetPressedButtons())
		queue.UpdateQueueFloorReached()
	}
	queue.SortQueue()
	log.Println("NextOrder: ", queue.GetNextOrder())
	io.GoToNextFloor(queue.GetNextOrder())
	if queue.GetNextOrder() == -1 {
		io.SetElevStateDir(config.DIR_STOP)
	} 
}

//EventManager
func main() {
	asd := make(chan int, 1)
	floorReached := make(chan bool, 1)
	buttonPressed := make(chan int, 1)
	var varButtonPressed int
	driver.Elev_init()
	io.SetElevState(driver.Elev_get_floor_sensor_signal(), 0, -1)
	log.Println("Hei")

	io.InitListeners(buttonPressed, floorReached)

	log.Println(config.ColC, "Test Run Initialized", config.ColN)
	queue.InitQueue()
	for {
		select {
		case varButtonPressed = <-buttonPressed:
			eventButtonPushed(varButtonPressed)
		case <-floorReached:
			eventFloorReached()
		default:
			break
		}
	}
	log.Println("Some shit got fucked")
	asd <- 1
}
