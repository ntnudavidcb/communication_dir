package main

//Må gjøres: Ordne sånn at man sender alive messages hele tiden?
//Problemet er at den ikke ble oppdatert når heisen snudde
//i etasjen,

import (
	"../com"
	"../config"
	. "../event"
	"../io"
	"log"
	"time"
)

func main() {
	ipListChannel := make(chan []string)
	timeStampMap := make(map[string]time.Time) //Holde styr pa timestamps paa IP adressene som blir sendt inn
	timer := make(chan bool, 1)
	floorReached := make(chan bool, 1)
	buttonPushed := make(chan int, 1)
	sendAliveMessage := make(chan com.Message, 1)
	messageRecieved := make(chan com.Message, 1)
	disconnected := make(chan bool, 1)

	InitElevator(buttonPushed, floorReached)

	log.Println(config.ColC, "Elevator Initialized", config.ColN)

	go com.Server(ipListChannel, sendAliveMessage, messageRecieved, disconnected)
	go io.ReadAllButtons(buttonPushed)
	go io.FloorSignalListener(floorReached)
	go TimerCount(timer)
	go ButtonPushedHandler(buttonPushed, disconnected)
	go MsgRecievedHandler(messageRecieved, timeStampMap)
	go FloorReachedHandler(floorReached, timer, sendAliveMessage)
	go HandleDisconnnect(disconnected, ipListChannel, sendAliveMessage, messageRecieved)
	go SafeKill() //If user ends the program ( CTRL + C )

	done := make(chan bool)
	<-done
}
