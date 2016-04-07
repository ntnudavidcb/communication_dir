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
	msg := make(chan string)
	connected := make(chan bool)
	timed := false

	go timeout(timeoutChannel)
	go connListener(connIPAddrs, msg, connected, port)
	go statusUpdater(addr)

	for {
		//Sjekke om noen sier at de er koblet til, oppdaterer IP-listen
		select {
		case <-connected:
			ip := <-connIPAddrs
			updateIP(timeStampVar, ip)
			log.Println(<-msg)
		default:
			break
		}

		TimeStampCheck(timeStampVar) //GetMyIP()

		log.Println("PC1")
		time.Sleep(500 * time.Millisecond)

		select {
		case <-timeoutChannel:
			log.Println("Timed out")
			timed = true
		default:
			break
		}

		if timed {
			break
		}
	}
	var ipListAlive []string
	for key, _ := range timeStampVar {
		ipListAlive = append(ipListAlive, key)
	}
	ipListChannel <- ipListAlive
}

func connListener(IPAddrs chan string, msg chan string, connected chan bool, port string) {
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
		lenOfBuffer, ip, err := udpListen.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		var m Message
		m = DecodeJSON(buffer[:lenOfBuffer])
		log.Println(m.Name)
		connected <- true
		IPAddrs <- ip.String()
		msg <- string(buffer)
	}
}

func updateIP(list map[string]time.Time, IPAddrs string) {
	IPWithoutPort := strings.Split(IPAddrs, ":")
	list[IPWithoutPort[0]] = time.Now().Add(2 * time.Second)
}

func timeout(ch chan bool) {
	time.Sleep(150 * time.Second)
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
