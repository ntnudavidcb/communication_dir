package com

import (
	"log"
	"net"
	"time"
)


var myBIP string
var myIP string
var port string

//Kjorer hele tiden, holder hele nettverkssystemet oppe og styrer hele showet
func Server(ipListChannel chan []string, sendAliveMessage chan Message, msgRecievedChan chan Message) {
	//connIPAddrsChan := make(chan string)           //Brukes til a kunne sende IP adresser imellom connListener og selve serveren
	//connectedChan := make(chan bool)               //om man er koblet til eller ikke
	port = "20010"
	myIP = GetMyIP()
	myBIP = GetBIP(myIP)

	udpBroadcast := getUDPcon()
	
	defer udpBroadcast.Close()

	//Kjorer et par trader som kjorer ved siden av server, alle kjorer uendelig med unntak av timeout
	go connListener(msgRecievedChan) //connIPAddrsChan, , connectedChan, port

	//for-loop som holder serveren oppe
	for {	
		msg := <-sendAliveMessage
		udpBroadcast.Write(CreateJSON(msg))
	}
}

//Lytter etter koblinger, faar tak i meldingen og IP for saa aa sende dette videre
func connListener(msgRecievedChan chan Message){ //connectedChan chan bool) { IPAddrs chan string,
	udpListen := getUDPconListener()

	defer udpListen.Close()

	buffer := make([]byte, 1024)
	for {
		lenOfBuffer, _, err := udpListen.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		//DecodeJSON returnerer en struct message av meldingen
		msg := DecodeJSON(buffer[:lenOfBuffer]) //&m)

		if msg.Name != GetMyIP(){
			msgRecievedChan <- msg
		}
	}
}

//Oppdaterer IP-listen med den IP'en som sendes inn
func updateIP(timeStampMap map[string]time.Time, IPAddr string) {
	timeStampMap[IPAddr] = time.Now().Add(6 * time.Second)
}

//Sjekker om timeStampen er good, ellers fjerner den IP adressen, bor fikse denne funksjonen etterhvert
func timeStampCheck(timeStampMap map[string]time.Time) { //MyIP string

	for key, val := range timeStampMap {
		if val.Before(time.Now()) { //&& key != MyIP {
			log.Println("Found disconnect")
			delete(timeStampMap, key)
		}
	}
}

func printAliveList(timeStampChan map[string]time.Time) {
	var ipListAlive []string
	for key, _ := range timeStampChan {
		ipListAlive = append(ipListAlive, key)
	}
	log.Println(ipListAlive)
}

func getUDPconListener() *net.UDPConn{
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		log.Println("connListener failed")
		log.Fatal(err)
	}
	log.Println(udpAddr)
	udpListen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println("connListener failed")
		log.Fatal(err)
	}
	return udpListen
}

func getUDPcon() *net.UDPConn{
    serverAddr, err := net.ResolveUDPAddr("udp",myBIP+":"+port)
    if err != nil {
            log.Println("getUDPcon failed")
            log.Fatal(err)
    }
    con, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
            log.Println("getUDPcon failed")
            log.Fatal(err)
    }
    return con
}

func CheckDisconnection(timeStampMap map[string]time.Time, messageStruct Message){
	updateIP(timeStampMap, messageStruct.Name)
	timeStampCheck(timeStampMap)
	//printAliveList(timeStampMap)
}