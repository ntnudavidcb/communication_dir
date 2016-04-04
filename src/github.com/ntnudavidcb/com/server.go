package com

import (
	"log"
	"net"
	"time"
)

func ListenUdp(port string, ipListChannel chan []string) {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatal(err)
	}
	udpListen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer udpListen.Close()

	ipList := make([]string, 0)
	var buffer [1024]byte

	timer := make(chan bool, 1)
	timeout := false

	go timerout(timer)

	for {

		_, ipAddr, err := udpListen.ReadFromUDP(buffer[:])
		if err != nil {
			log.Fatal(err)
		}

		if !ipInList(ipAddr.String(), ipList) {
			ipList = append(ipList, ipAddr.String())
		}
		log.Println(string(buffer[0:10]))

		select {
		case <-timer:
			timeout = true
		default:
			break
		}
		if timeout {
			break
		}
		log.Println("PC1")
		time.Sleep(1000 * time.Millisecond)
	}
	log.Println("Server ended")
	ipListChannel <- ipList
}

func ipInList(ipAddr string, ipList []string) bool {
	for _, b := range ipList {
		if b == ipAddr {
			return true
		}
	}
	return false
}

func timerout(timer chan bool) {
	time.Sleep(10 * time.Second)
	timer <- true
}
