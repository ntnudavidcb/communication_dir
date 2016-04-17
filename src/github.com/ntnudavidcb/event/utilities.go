package event

import (
	"../com"
	. "../config"
	"../driver"
	"log"
	"os"
	"os/signal"
	"time"
)

func SafeKill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	driver.Elev_set_motor_direction(DIR_STOP)
	log.Fatal(ColR, "User terminated program.", ColN)
}

//TRÅDENE------------------------------------------------------------------
func TimerCount(timerchan chan bool) {
	time.Sleep(200 * time.Millisecond)
	timerchan <- true
}

func TimedOut(timerChan chan bool) {
	time.Sleep(10 * time.Second)
	timerChan <- true
}

func timeFloor(timedOutFloor chan bool) {
	time.Sleep(6 * time.Second)
	timedOutFloor <- true
}

func ButtonPushedHandler(buttonPressed chan int, disconnected chan bool) {
	for {
		varButtonPressed := <-buttonPressed
		eventButtonPushed(varButtonPressed, disconnected)
		time.Sleep(10 * time.Millisecond)
	}
}

func MsgRecievedHandler(messageRecieved chan com.Message) {
	for {
		msg := <-messageRecieved
		log.Println("Message recieved")
		eventMessageRecieved(msg)
		time.Sleep(10 * time.Millisecond)
	}
}

func FloorReachedHandler(floorReached chan bool, timer chan bool, sendAliveMessage chan com.Message) {
	timedOutChan := make(chan bool, 1)
	engineActive := make(chan bool)
	timedOutFloor := make(chan bool, 1)
	activeDirection := DIR_STOP
	timeStamp := time.Now()
	for {
		select {
		case <-floorReached:
			break
		case <-timedOutFloor:
			break
		default:
			break
		}
		activeDirection, timeStamp = eventFloorReached(engineActive, sendAliveMessage, timer, timedOutChan, activeDirection, timeStamp)
		if timeStamp.Before(time.Now()) {
			Restart.Run()
			log.Fatal("Not any action taken for too long, engine might not work.")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func HandleDisconnnect(disconnected chan bool, ipListChannel chan []string, sendAliveMessage chan com.Message, messageRecieved chan com.Message) {
	for {
		<-disconnected
		log.Println(ColR, "Disconnected from the network.", ColN)
		eventDisconnect(disconnected, ipListChannel, sendAliveMessage, messageRecieved)
		time.Sleep(10 * time.Millisecond)
	}
}
