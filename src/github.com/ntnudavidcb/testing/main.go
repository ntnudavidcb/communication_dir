package main

import (
	"github.com/ntnudavidcb/communication"
	"log"
)

func main() {
	log.Println(communication.GetMyIP())
	doneChannel := make(chan bool, 1)
	ipListChannel := make(chan []string, 1)
	doneChannel <- true
	port := ":20010"

	broadcastAddr := communication.GetBIP(communication.GetMyIP()) + port

	go communication.BroadcastUdp(broadcastAddr)

	log.Println(<-ipListChannel)
	<-doneChannel
}
