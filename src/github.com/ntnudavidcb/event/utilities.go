package event

import (
	"../com"
	"../config"
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
	driver.Elev_set_motor_direction(config.DIR_STOP)
	log.Fatal(config.ColR, "User terminated program.", config.ColN)
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

func ButtonPushedHandler(buttonPressed chan int, disconnected chan bool) {
	for {
		varButtonPressed := <-buttonPressed
		eventButtonPushed(varButtonPressed, disconnected)
		time.Sleep(10 * time.Millisecond)
	}
}

func MsgRecievedHandler(messageRecieved chan com.Message, timeStampMap map[string]time.Time) {
	for {
		msg := <-messageRecieved
		log.Println("Message recieved")
		eventMessageRecieved(msg, timeStampMap)
		time.Sleep(10 * time.Millisecond)
	}
}

func FloorReachedHandler(floorReached chan bool, timer chan bool, sendAliveMessage chan com.Message) {
	timedOutChan := make(chan bool, 1)
	engineActive := make(chan bool)
	activeDirection := config.DIR_STOP
	timeStamp := time.Now()
	for {
		<-floorReached
		activeDirection, timeStamp = eventFloorReached(engineActive, sendAliveMessage, timer, timedOutChan, activeDirection, timeStamp)
		if timeStamp.Before(time.Now()) {
			config.Restart.Run()
			log.Fatal("Not any action taken for too long, engine might not work.")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func HandleDisconnnect(disconnected chan bool, ipListChannel chan []string, sendAliveMessage chan com.Message, messageRecieved chan com.Message) {
	for {
		<-disconnected
		eventDisconnect(disconnected, ipListChannel, sendAliveMessage, messageRecieved)
		time.Sleep(10 * time.Millisecond)
	}
}
