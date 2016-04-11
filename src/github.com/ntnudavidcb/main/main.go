package main

import (
	"../config"
	"../driver"
	"../io"
	"../queue"
	"../com"
	"log"
	"time"
	"strings"
)

func eventButtonPushed(buttonPushed int) {
	queue.UpdateQueueWithButton(buttonPushed)
}

func eventMessageRecieved(messageStruct com.Message){
	queue.UpdateElevStateMap(messageStruct.Name, messageStruct.Direction, messageStruct.Floor)

}

func eventFloorReached(sendMessage chan com.Message, timer chan bool) {
	//log.Println("CheckOrder: ", queue.CheckOrder())
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

	select{
	case <-timer:
		m := com.Message{"Alive", 0, 0, 0, 0, time.Now()}
		sendMessage <- m
		go timerCount(timer)
	default:
		break
	}
}

func timerCount(timerchan chan bool){
	time.Sleep(1000 * time.Millisecond)
	timerchan <- true
}

func InitElevStateMap(floor int){
	IPAddr := com.GetMyIP()
	IPAddrWithoutPort := strings.Split(IPAddr, ":")
	queue.UpdateElevStateMap(IPAddrWithoutPort[0], floor, 0)
}

func InitMyIP(){
	IPAddr := com.GetMyIP()
	IPAddrWithoutPort := strings.Split(IPAddr, ":")
	queue.SetMyIP(IPAddrWithoutPort[0])
}

//EventManager
func main() {
	timer := make(chan bool, 1)
	asd := make(chan int, 1)
	floorReached := make(chan bool, 1)
	buttonPressed := make(chan int, 1)
	var varButtonPressed int

	//Set up server
	ipListChannel := make(chan []string, 1)
	sendMessage := make(chan com.Message, 1)

	port := ":20010"

	go com.Server(com.GetBIP(com.GetMyIP()), port, ipListChannel, sendMessage)

	
	//Setup of server done

	floor,_ := driver.Elev_init()
	InitMyIP()
	InitElevStateMap(floor)
	go timerCount(timer)
	io.SetElevState(floor, 0, -1)
	log.Println("Hei")

	io.InitListeners(buttonPressed, floorReached)

	log.Println(config.ColC, "Test Run Initialized", config.ColN)
	queue.InitQueue()
	for {
		select {
		case varButtonPressed = <-buttonPressed:
			eventButtonPushed(varButtonPressed)
		case <-floorReached:
			eventFloorReached(sendMessage, timer)
		default:
			break
		}
		//log.Println(<-ipListChannel)
	}
	log.Println("Some shit got fucked")
	log.Println(<-ipListChannel)
	asd <- 1
}
