package event

import (
	"../com"
	. "../config"
	"../converter"
	"../costFunc"
	"../driver"
	"../io"
	"../queue"
	"log"
	"time"
)

func eventButtonPushed(buttonPushed int, disconnected chan bool) {
	queue.AddButtonToQueue(buttonPushed)
	if buttonPushed < CMD_1 { //Det vil si ingen av CMD knappene, da må den si ifra til de andre hvilken av de som ble trykket
		floor, direction, reserved := io.GetElevState()
		msg := com.Message{costFunc.MyIP, buttonPushed, floor, direction, reserved, NOT_ANY_BUTTON, time.Now()}
		com.SendMessage(msg, disconnected)
	}
}

func eventMessageRecieved(messageStruct com.Message, timeStampMap map[string]time.Time) {
	if messageStruct.ButtonPushed != NOT_ANY_BUTTON { //Hvis noen sier at en knapp ble trykket legges den til
		queue.AddButtonToQueue(messageStruct.ButtonPushed)
		io.SetPressedButton(messageStruct.ButtonPushed, true)
	}
	if messageStruct.Floor == NOT_ANY_FLOOR {
		Restart.Run()
		log.Fatal("eventMessageRecieved failed, messageStruct.Floor not valid.")
	}
	log.Println(messageStruct)
	queue.UpdateElevStateMap(messageStruct.Name, messageStruct.Floor, messageStruct.Direction, messageStruct.Reserved)
	queue.RemoveButtonFromQueue(messageStruct.OrderTaken)
	queue.RemoveButtonFromQueue(messageStruct.Reserved)
	io.RemoveButtonFromPressedButtonList(messageStruct.OrderTaken)
	io.RemoveButtonFromPressedButtonList(messageStruct.Reserved)
	io.UpdateLightsWithRecievedMessage(messageStruct.Floor, messageStruct.Direction)
	IPAddrs := com.CheckDisconnection(timeStampMap, messageStruct)
	if IPAddrs != "" {
		reserved := costFunc.ElevStateMap[IPAddrs].Reserved
		queue.AddButtonToQueue(reserved)
		io.SetPressedButton(reserved, true)
	}
}

func eventFloorReached(engineActive chan bool, sendAliveMessage chan com.Message, timer chan bool, timedOutChan chan bool, activeDirection int, timeStamp time.Time) (int, time.Time) {
	select {
	case <-timer:
		floor, direction, reserved := io.GetElevState()
		m := com.Message{com.GetMyIP(), NOT_ANY_BUTTON, floor, direction, reserved, NOT_ANY_BUTTON, time.Now()}
		sendAliveMessage <- m
		go TimerCount(timer)
	default:
		break
	}
	var newActiveDirection int
	floor, direction, reserved := io.GetElevState()
	driver.Elev_set_floor_indicator(floor)
	log.Println("eventFloorReached", floor, direction, reserved, queue.CheckOrder(floor, direction))
	log.Println(ColY, "GetPressedButtons: ", io.GetPressedButtons(), ColN)
	if queue.CheckOrder(floor, direction) {
		log.Println("Direction: ", direction)
		orderTaken := NOT_ANY_BUTTON
		if queue.CheckUpOrDownButton() == BTN_UP {
			orderTaken, _, _ = converter.ConvertDirAndFloorToMapIndex(floor, direction)
			m := com.Message{com.GetMyIP(), NOT_ANY_BUTTON, floor, direction, reserved, orderTaken, time.Now()}
			sendAliveMessage <- m
		} else if queue.CheckUpOrDownButton() == BTN_DOWN {
			_, orderTaken, _ = converter.ConvertDirAndFloorToMapIndex(floor, direction)
			m := com.Message{com.GetMyIP(), NOT_ANY_BUTTON, floor, direction, reserved, orderTaken, time.Now()}
			sendAliveMessage <- m
		}

		if converter.ConvertButtonToFloor(io.GetElevStateReserved()) == floor {
			if direction == DIR_UP && io.GetElevStateReserved() < DOWN_4 {
				io.SetElevStateReserved(NOT_ANY_BUTTON)
			} else if direction == DIR_DOWN && io.GetElevStateReserved() < CMD_1 && io.GetElevStateReserved() > UP_3 {
				io.SetElevStateReserved(NOT_ANY_BUTTON)
			} else if direction == DIR_DOWN && io.GetElevStateReserved() < DOWN_4 {
				io.SetElevStateDir(DIR_UP)
			} else if direction == DIR_UP && io.GetElevStateReserved() > UP_3 {
				io.SetElevStateDir(DIR_DOWN)
			}
			io.SetElevStateReserved(NOT_ANY_BUTTON)
		}

		io.HandleWantedFloorReached()
		queue.SynchronizeQueueWithIO(io.GetPressedButtons()) //To do: Denne er rar, men må være her
	} else if io.GetElevStateReserved() != NOT_ANY_BUTTON {
		log.Println("GetElevStateReserved: ", io.GetElevStateReserved())
		newActiveDirection = io.GoToNextFloor(converter.ConvertButtonToFloor(io.GetElevStateReserved()))
	} else {
		button, outside_button := queue.GetNextOrder()
		log.Println("outside_button: ", outside_button)

		if outside_button == OUTSIDE_BTN {
			io.SetElevStateReserved(button)
			m := com.Message{com.GetMyIP(), NOT_ANY_BUTTON, floor, direction, button, NOT_ANY_BUTTON, time.Now()}
			sendAliveMessage <- m
			time.Sleep(10 * time.Millisecond)
		}
		newActiveDirection = io.GoToNextFloor(converter.ConvertButtonToFloor(button))

		log.Println("Button to  GoToNextFloor: ", converter.ConvertButtonToFloor(button))
		log.Println("Inside eventFloorReached: NextOrder: (button, outside_button)", button, outside_button)
	}

	if activeDirection != newActiveDirection || newActiveDirection == DIR_STOP {
		timeStamp = time.Now().Add(10 * time.Second)
	}

	if queue.EmptyQueue() {
		io.SetElevStateDir(DIR_STOP)
		driver.Elev_set_motor_direction(DIR_STOP)
	}
	button, outside_button := queue.GetNextOrder()
	log.Println("eventFloorReached: GetNextOrder:", button, outside_button)
	queue.UpdateElevStateMap(costFunc.MyIP, io.GetElevStateFloor(), io.GetElevStateDir(), io.GetElevStateReserved())
	queue.UpdateQueue()

	return newActiveDirection, timeStamp
	//log.Println("Inside eventFloorReached: NextOrder: (button, outside_button)", button, outside_button)
}

func eventDisconnect(disconnected chan bool, ipListChannel chan []string, sendAliveMessage chan com.Message, messageRecieved chan com.Message) {
	for i := 0; i < CMD_1; i++ {
		queue.RemoveButtonFromQueue(i)
		io.RemoveButtonFromPressedButtonList(i)
	}
Loop:
	for {
		time.Sleep(100 * time.Millisecond)
		if driver.Elev_get_floor_sensor_signal() != -1 {
			io.InitElevState(driver.Elev_get_floor_sensor_signal())
			costFunc.DelElevStateMap()
			queue.UpdateElevStateMap(costFunc.MyIP, io.GetElevStateFloor(), io.GetElevStateDir(), io.GetElevStateReserved())
			break Loop
		}
	}
	go com.Server(ipListChannel, sendAliveMessage, messageRecieved, disconnected)
}
