package main

import (
	//"../config"
	"../driver"
	"../io"
	"../queue"
	"log"
	//"time"
)

func eventButtonPushed(buttonPushed int) {
	//io.UpdateLights()
	queue.UpdateQueue(buttonPushed)
}

func eventFloorReached() {
	if queue.CheckOrder() {
		log.Println("It should have stopped here")
		io.WantedFloorReached()
		queue.RemoveFromQueue() //ta inn elevState eller no?)
	}

}

func getNextInQueue() {

}

//EventManager
func main() {
	asd := make(chan int, 1)
	floorReached := make(chan bool, 1)
	buttonPressed := make(chan int, 1)
	nextFloor := make(chan int, 1)
	//var floor int
	//var direction int
	//var NextFloor int

	driver.Elev_init()
	io.SetElevState(driver.Elev_get_floor_sensor_signal(), 0)
	log.Println("Hei")
	go io.Testrun2(floorReached, buttonPressed, nextFloor)

	for {
		select {
		//Reaction when a button is pressed
		case varButtonPressed <- buttonPressed:
			eventButtonPushed(varButtonPressed)
			//Reaction when a floor is reached
		case <-floorReached:
			eventFloorReached()
		default:
			break
		}
		nextFloor <- queue.GetNextOrder()
	}
	log.Println("FUCK")
	asd <- 1
}
