package main

import (
	"github.com/ntnudavidcb/communication"
	"log"
)

func main() {
	doneChannel := make(chan bool, 1)
	ipListChannel := make(chan []string, 1)

	port := ":20010"

	go communication.ListenUdp(port, ipListChannel)

	log.Println(<-ipListChannel)
	<-doneChannel
}
