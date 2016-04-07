package com

import (
	"log"
	"net"
	"strings"
	"time"
)

func Server(addr string, port string, ipListChannel chan []string) {
	timeStampVar := make(map[string]time.Time)
	timeStampVar[GetMyIP()] = time.Now().Add(30000 * time.Second)
	timeoutChannel := make(chan bool)
	connIPAddrs := make(chan string)
	msgRecieved := make(chan string)
	connected := make(chan bool)
	timedOut := false

	go timeout(timeoutChannel)
	go connListener(connIPAddrs, msgRecieved, connected, port)
	go statusUpdater(addr, port)

	for {
		//Sjekke om noen sier at de er koblet til, oppdaterer IP-listen, faa melding/msg
		select {
		case <-connected:
			ip := <-connIPAddrs
			updateIP(timeStampVar, ip)
			<-msgRecieved
		default:
			break
		}

		TimeStampCheck(timeStampVar) //GetMyIP()

		printAliveList(timeStampVar)
		time.Sleep(500 * time.Millisecond)

		select {
		case <-timeoutChannel:
			log.Println("Timed out")
			timedOut = true
		default:
			break
		}

		if timedOut {
			break
		}
	}
	var ipListAlive []string
	for key, _ := range timeStampVar {
		ipListAlive = append(ipListAlive, key)
	}
	ipListChannel <- ipListAlive
}

func connListener(IPAddrs chan string, msgContainer chan string, connected chan bool, port string) {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatal(err)
	}
	udpListen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpListen.Close()

	buffer := make([]byte, 1024)
	for {
		lenOfBuffer, addressFromReciever, err := udpListen.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		//var m Message
		_ = DecodeJSON(buffer[:lenOfBuffer])
		connected <- true
		IPAddrs <- addressFromReciever.String()
		msgContainer <- string(buffer)
	}
}

func updateIP(list map[string]time.Time, IPAddrWithPort string) {
	IPAddrWithoutPort := strings.Split(IPAddrWithPort, ":")
	list[IPAddrWithoutPort[0]] = time.Now().Add(2 * time.Second)
}

func timeout(ch chan bool) {
	time.Sleep(60 * time.Second)
	ch <- true
}

func TimeStampCheck(list map[string]time.Time) { //MyIP string
	for key, val := range list {
		if val.Before(time.Now()) { //&& key != MyIP {
			log.Println("Found disconnect")
			delete(list, key)
			break
		}
	}
}

func printAliveList(timeStampVar map[string]time.Time) {
	var ipListAlive []string
	for key, _ := range timeStampVar {
		ipListAlive = append(ipListAlive, key)
	}
	log.Println(ipListAlive)
}
