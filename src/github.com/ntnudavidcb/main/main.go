package main

import (
	"../config"
	"../costFunc"
	"../driver"
	"../io"
	"../queue"
	"../com"
	"../converter"
	"log"
	"time"
	"strings"
)

func eventButtonPushed(buttonPushed int) {
	queue.UpdateQueueWithButton(buttonPushed)
	if buttonPushed  < costFunc.CMD_1 {
		floor, direction, reserved := io.GetElevState()
		msg := com.Message{com.GetMyIP(), buttonPushed,floor, direction, reserved, time.Now()}
		com.SendMessage(msg)
	}
}

func eventMessageRecieved(messageStruct com.Message, timeStampMap map[string]time.Time){
	if messageStruct.Name == com.GetMyIP(){
		return
	}else if messageStruct.ButtonPushed != -1 {
		queue.UpdateQueueWithButton(messageStruct.ButtonPushed)
	}
	queue.UpdateElevStateMap(messageStruct.Name, messageStruct.Direction, messageStruct.Floor)
	com.CheckDisconnection(timeStampMap, messageStruct)
}

func eventFloorReached(sendAliveMessage chan com.Message, timer chan bool) {
	//log.Println("CheckOrder: ", queue.CheckOrder())
	queue.UpdateElevStateMap(com.GetMyIP(), io.GetElevStateDir(), io.GetElevStateFloor())
	if queue.CheckOrder() {
		io.StopAtFloorReached()
		queue.RemoveFromQueue(io.GetPressedButtons())
		queue.UpdateQueueFloorReached()
	}

	queue.SortQueue()
	button, outside_button := queue.GetNextOrder()
	log.Println("NextOrder: (button, outside_button)", button, outside_button)
	if outside_button == 2{
		io.SetElevStateReserved(button)
	}
	log.Println("Button to  GoToNextFloor: " , converter.ConvertButtonToFloor(button))
	io.GoToNextFloor(converter.ConvertButtonToFloor(button))
	if button == -1 {
		io.SetElevStateDir(config.DIR_UP)
	}

	select{
	case <-timer:
		floor, direction, reserved := io.GetElevState()
		m := com.Message{com.GetMyIP(), -1, floor, direction, reserved, time.Now()}
		sendAliveMessage <- m
		go timerCount(timer)
	default:
		break
	}
}

func timerCount(timerchan chan bool){
	time.Sleep(1000 * time.Millisecond)
	timerchan <- true
}

func initElevStateMap(floor int){
	IPAddr := com.GetMyIP()
	IPAddrWithoutPort := strings.Split(IPAddr, ":")
	queue.UpdateElevStateMap(IPAddrWithoutPort[0], floor, 0)
}

func initMyIP(){
	IPAddr := com.GetMyIP()
	IPAddrWithoutPort := strings.Split(IPAddr, ":")
	queue.SetMyIP(IPAddrWithoutPort[0])
}

func initElevator(buttonPressed chan int, floorReached chan bool) int { //return floor
	floor,_ := driver.Elev_init()
	initMyIP()
	initElevStateMap(floor)
	io.InitListeners(buttonPressed, floorReached)
	return floor
}

//EventManager
func main() {
	timeStampMap := make(map[string]time.Time) //Holde styr pa timestamps paa IP adressene som blir sendt inn
	timer := make(chan bool, 1)
	floorReached := make(chan bool, 1)
	buttonPressed := make(chan int, 1)
	ipListChannel := make(chan []string, 1)
	sendAliveMessage := make(chan com.Message, 1)
	messageRecieved := make(chan com.Message, 1)

	//Set up server
	go com.Server(ipListChannel, sendAliveMessage, messageRecieved)
	go timerCount(timer)

	floor := initElevator(buttonPressed, floorReached)
	io.SetElevState(floor, 0, -1)

	log.Println(config.ColC, "Elevator Initialized", config.ColN)
	for {
		select {
		case varButtonPressed := <-buttonPressed:
			eventButtonPushed(varButtonPressed)
		case msg := <- messageRecieved:
			eventMessageRecieved(msg, timeStampMap)
		case <-floorReached:
			eventFloorReached(sendAliveMessage, timer)
		default:
			break
		}
	}
	log.Println("Some shit got fucked")
	log.Println(<-ipListChannel)
}
