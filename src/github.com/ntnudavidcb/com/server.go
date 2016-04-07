package com

import (
	"log"
	"net"
	"strings"
	"time"
)

//Kjorer hele tiden, holder hele nettverkssystemet oppe og styrer hele showet
func Server(addr string, port string, ipListChannel chan []string) {
	timeStampVar := make(map[string]time.Time) //Holde styr pa timestamps paa IP adressene som blir sendt inn
	timeoutChannel := make(chan bool)          //En enkel timeout for a teste systemet, fjernes etterhvert
	connIPAddrs := make(chan string)           //Brukes til a kunne sende IP adresser imellom connListener og selve serveren
	msgRecieved := make(chan string)           //Meldinger som man har fatt
	connected := make(chan bool)               //om man er koblet til eller ikke
	timedOut := false                          //Brukes for a avslutte for lokken i server

	//Kjorer et par trader som kjorer ved siden av server, alle kjorer uendelig med unntak av timeout
	go timeout(timeoutChannel)
	go connListener(connIPAddrs, msgRecieved, connected, port)
	go statusUpdater(addr, port)

	//for-loop som holder serveren oppe
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

//Lytter etter koblinger, faar tak i meldingen og IP for saa aa sende dette videre
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
		//DecodeJSON returnerer en struct message av meldingen
		_ = DecodeJSON(buffer[:lenOfBuffer]) //&m)
		connected <- true
		IPAddrs <- addressFromReciever.String()
		msgContainer <- string(buffer)
	}
}

//Oppdaterer IP-listen med den IP'en som sendes inn
func updateIP(list map[string]time.Time, IPAddrWithPort string) {
	IPAddrWithoutPort := strings.Split(IPAddrWithPort, ":")
	list[IPAddrWithoutPort[0]] = time.Now().Add(2 * time.Second)
}

func timeout(ch chan bool) {
	time.Sleep(60 * time.Second)
	ch <- true
}

//Sjekker om timeStampen er good, ellers fjerner den IP adressen, bor fikse denne funksjonen etterhvert
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
