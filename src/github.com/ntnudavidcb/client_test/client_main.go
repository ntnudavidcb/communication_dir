package main

import (
	"github.com/ntnudavidcb/com"
	"log"
)

func main() {
	log.Println(com.GetMyIP())
	doneChannel := make(chan bool, 1)
	ipListChannel := make(chan []string, 1)
	doneChannel <- true
	port := ":20010"

	broadcastAddr := com.GetBIP(com.GetMyIP()) + port

	go com.BroadcastUdp(broadcastAddr)

	log.Println(<-ipListChannel)
	<-doneChannel
}
