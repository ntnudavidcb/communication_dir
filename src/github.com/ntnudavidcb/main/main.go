package main

//Må gjøres: Ordne sånn at man sender alive messages hele tiden?
//Problemet er at den ikke ble oppdatert når heisen snudde
//i etasjen, 

import (
	"../config"
	//"../costFunc"
	"../driver"
	"../io"
	"../queue"
	"../com"
	"../converter"
	"log"
	"time"
)

//INIT FUNKSJONENE-------------------------------------------------------
func initElevator(buttonPressed chan int, floorReached chan bool) { 
	floor,_ := driver.Elev_init()
	io.SetElevState(floor, config.DIR_STOP, config.NOT_ANY_BUTTON)
	initMyIP()
	initElevStateMap(floor)
	io.InitElevState(floor)
	queue.InitQueue()
	//io.InitButtonAndFloorListeners(buttonPressed, floorReached)
}

func initMyIP(){
	IPAddr := com.GetMyIP()
	queue.SetMyIP(IPAddr)
}

//Godt mulig en funksjon som kan fjernes, finner det ut
func initElevStateMap(floor int){
	IPAddr := com.GetMyIP()
	queue.UpdateElevStateMap(IPAddr, floor, config.DIR_STOP)
}




//EVENT FUNKSJONER-------------------------------------------------------
func eventButtonPushed(buttonPushed int) {
	queue.AddButtonToQueue(buttonPushed)
	if buttonPushed  < config.CMD_1 { //Det vil si ingen av CMD knappene, da må den si ifra til de andre hvilken av de som ble trykket
		floor, direction, reserved := io.GetElevState()
		msg := com.Message{com.GetMyIP(), buttonPushed,floor, direction, reserved, config.NOT_ANY_BUTTON, time.Now()}
		com.SendMessage(msg)
	}
}

func eventMessageRecieved(messageStruct com.Message, timeStampMap map[string]time.Time){
	if messageStruct.ButtonPushed != config.NOT_ANY_BUTTON { //Hvis noen sier at en knapp ble trykket legges den til
		queue.AddButtonToQueue(messageStruct.ButtonPushed)
		io.SetPressedButton(messageStruct.ButtonPushed)
	}
	if messageStruct.Floor == config.NOT_ANY_FLOOR{
		return
	}
	queue.UpdateElevStateMap(messageStruct.Name, messageStruct.Floor, messageStruct.Direction)
	queue.RemoveButtonFromQueue(messageStruct.OrderTaken)
	queue.RemoveButtonFromQueue(messageStruct.Reserved)
	io.RemoveButtonFromPressedButtonList(messageStruct.OrderTaken)
	io.RemoveButtonFromPressedButtonList(messageStruct.Reserved)
	com.CheckDisconnection(timeStampMap, messageStruct)
}

func eventFloorReached(sendAliveMessage chan com.Message, timer chan bool) {
	select{
	case <-timer:
		floor, direction, reserved := io.GetElevState()
		m := com.Message{com.GetMyIP(), config.NOT_ANY_BUTTON, floor, direction, reserved,config.NOT_ANY_BUTTON, time.Now()}
		sendAliveMessage <- m
		go timerCount(timer)
	default:
		break
	}
	floor, direction, reserved := io.GetElevState()

	if queue.CheckOrder(floor, direction){
		log.Println("Direction: ", direction)
		orderTaken := config.NOT_ANY_BUTTON
		if queue.CheckUpOrDownButton() == config.BTN_UP{
			orderTaken, _, _ = converter.ConvertDirAndFloorToMapIndex(floor, direction)
			m := com.Message{com.GetMyIP(), config.NOT_ANY_BUTTON, floor, direction, reserved, orderTaken , time.Now()}
			sendAliveMessage <- m
		} else if queue.CheckUpOrDownButton() == config.BTN_DOWN{
			_, orderTaken, _ = converter.ConvertDirAndFloorToMapIndex(floor, direction)
			m := com.Message{com.GetMyIP(), config.NOT_ANY_BUTTON, floor, direction, reserved, orderTaken , time.Now()}
			sendAliveMessage <- m
		} 
		
		if converter.ConvertButtonToFloor(io.GetElevStateReserved()) == floor{
			if direction == config.DIR_UP && io.GetElevStateReserved() < config.DOWN_4{
				io.SetElevStateReserved(config.NOT_ANY_BUTTON)	
			} else if direction == config.DIR_DOWN && io.GetElevStateReserved() < config.CMD_1 && io.GetElevStateReserved() > config.UP_3{
				io.SetElevStateReserved(config.NOT_ANY_BUTTON)	
			} else if direction == config.DIR_DOWN && io.GetElevStateReserved() < config.DOWN_4{
				io.SetElevStateDir(config.DIR_UP)
			} else if direction == config.DIR_UP && io.GetElevStateReserved() > config.UP_3{
				io.SetElevStateDir(config.DIR_DOWN)
			} 
			io.SetElevStateReserved(config.NOT_ANY_BUTTON)	
		}
		
		io.HandleWantedFloorReached()
		queue.SynchronizeQueueWithIO(io.GetPressedButtons()) //Denne er rar, men må være her
	}else if io.GetElevStateReserved() != config.NOT_ANY_BUTTON{
		log.Println("GetElevStateReserved: ", io.GetElevStateReserved())
		io.GoToNextFloor(converter.ConvertButtonToFloor(io.GetElevStateReserved()))
	} else{
		button, outside_button := queue.GetNextOrder()
		log.Println("outside_button: ", outside_button)
		if io.GetElevStateReserved() != config.NOT_ANY_BUTTON{
		} else {
			if outside_button == 2{
				io.SetElevStateReserved(button)
				m := com.Message{com.GetMyIP(), config.NOT_ANY_BUTTON, floor, direction, button, config.NOT_ANY_BUTTON , time.Now()}
				sendAliveMessage <- m
				time.Sleep(500 * time.Millisecond)
			} 
			io.GoToNextFloor(converter.ConvertButtonToFloor(button))
		}
		log.Println("Button to  GoToNextFloor: " , converter.ConvertButtonToFloor(button))
		log.Println("Inside eventFloorReached: NextOrder: (button, outside_button)", button, outside_button)
	}
	
	if queue.EmptyQueue(){
		io.SetElevStateDir(config.DIR_STOP)
		driver.Elev_set_motor_direction(config.DIR_STOP)
	} 
	button, outside_button := queue.GetNextOrder()
	log.Println("eventFloorReached: GetNextOrder:",button, outside_button)
	if io.GetElevStateFloor() == -1{
		return 
	}
	queue.UpdateElevStateMap(com.GetMyIP(), io.GetElevStateFloor(), io.GetElevStateDir())
	queue.UpdateQueue()
	
	//log.Println("Inside eventFloorReached: NextOrder: (button, outside_button)", button, outside_button)

}





//TRÅDENE------------------------------------------------------------------
func timerCount(timerchan chan bool){
	time.Sleep(200 * time.Millisecond)
	timerchan <- true
}

func buttonPushedHandler(buttonPressed chan int){
	for {
		varButtonPressed := <-buttonPressed
		eventButtonPushed(varButtonPressed)

	}
}

func msgRecievedHandler(messageRecieved chan com.Message, timeStampMap map[string]time.Time){
	for {
		msg := <- messageRecieved
		log.Println("Message recieved")
		eventMessageRecieved(msg, timeStampMap)
		time.Sleep(10*time.Millisecond)
	}
}

func floorReachedHandler(floorReached chan bool, timer chan bool, sendAliveMessage chan com.Message){
	for{
		<-floorReached
		eventFloorReached(sendAliveMessage, timer)
	}
}




//EventManager
func main() {
	ipListChannel := make(chan []string)
	timeStampMap := make(map[string]time.Time) //Holde styr pa timestamps paa IP adressene som blir sendt inn
	timer := make(chan bool, 1)
	floorReached := make(chan bool, 1)
	buttonPushed := make(chan int, 1)
	sendAliveMessage := make(chan com.Message, 1)
	messageRecieved := make(chan com.Message, 1)

	initElevator(buttonPushed, floorReached)

	log.Println(config.ColC, "Elevator Initialized", config.ColN)

	go com.Server(ipListChannel, sendAliveMessage, messageRecieved)
	go io.ReadAllButtons(buttonPushed)
	go io.FloorSignalListener(floorReached)
	go timerCount(timer)
	go buttonPushedHandler(buttonPushed)
	go msgRecievedHandler(messageRecieved, timeStampMap)
	go floorReachedHandler(floorReached, timer, sendAliveMessage)

	snorre := make(chan bool)
	<-snorre
}