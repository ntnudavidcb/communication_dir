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


	timeStampVar := make(map[string]time.Time)
	timeStampVar[GetMyIP()] = time.Now().Add(30000*time.Second)
	timeoutChannel := make(chan bool)
	go timeout(timeoutChannel)
	timed := false

	var buffer [1024]byte
	for {
		log.Println(port)
		_, ipAddr, err := udpListen.ReadFromUDP(buffer[:])
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(buffer[0:10]))

		updateIP(timeStampVar, ipAddr.String())
		TimeStampCheck(timeStampVar, ipAddr.String())

		log.Println("PC1")
		time.Sleep(1000 * time.Millisecond)

		select{
		case <- timeoutChannel:
			log.Println("Timed out")
			timed = true
		default:
			break
		}

		if timed{
			break
		}
	}
	var ipListAlive []string
	for key, _ := range timeStampVar {
		ipListAlive = append(ipListAlive, key)
		log.Println(key)
	}
	log.Println("Server ended")
	ipListChannel <- ipListAlive
	log.Println("SJEKKE DENNE")
}

func updateIP(list map[string]time.Time, IPAddrs string){
	list[IPAddrs] = time.Now().Add(2*time.Second)
}

func timeout(ch chan bool){
	time.Sleep(5*time.Second)
	ch <- true
}

func TimeStampCheck(list map[string]time.Time, MyIP string) {
		for key, val := range list {
			if val.Before(time.Now()) && key != MyIP {
				delete(list, key)
				break
			}
		}
}
