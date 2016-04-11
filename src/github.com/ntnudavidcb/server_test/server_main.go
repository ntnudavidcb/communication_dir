package main

import (
	"github.com/ntnudavidcb/com"
	"log"
	"time"
)

func main() {
	doneChannel := make(chan bool, 1)
	ipListChannel := make(chan []string, 1)
	sendMessage := make(chan com.Message, 1)

	port := ":20010"

	go com.Server(com.GetBIP(com.GetMyIP()), port, ipListChannel, sendMessage)

	m := com.Message{"Alive", 0, 0, 0, 0, time.Now()}
	sendMessage <- m
	log.Println(<-ipListChannel)
	<-doneChannel
}
