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
func Server(ipListChannel chan []string, sendAliveMessage chan Message, msgRecievedChan chan Message, disconnected chan bool) {
	//connIPAddrsChan := make(chan string)           //Brukes til a kunne sende IP adresser imellom connListener og selve serveren
	//connectedChan := make(chan bool)               //om man er koblet til eller ikke
	port = "20010"
	myIP = GetMyIP()
	myBIP = GetBIP(myIP)

	if len(myIP) < 7{
		disconnected<- true
		time.Sleep(10*time.Millisecond)
		return
	}

	udpBroadcast := getUDPcon(disconnected)
	
	defer udpBroadcast.Close()

	//Kjorer et par trader som kjorer ved siden av server, alle kjorer uendelig med unntak av timeout
	go connListener(msgRecievedChan, disconnected) //connIPAddrsChan, , connectedChan, port

	//for-loop som holder serveren oppe
	for {	
		msg := <-sendAliveMessage
		_, err := udpBroadcast.Write(CreateJSON(msg))
		if err != nil{
			log.Println("sendAliveMessage failed")
			disconnected <- true
			break
		}
		time.Sleep(100*time.Millisecond)
	}
}

//Lytter etter koblinger, faar tak i meldingen og IP for saa aa sende dette videre
func connListener(msgRecievedChan chan Message, disconnected chan bool){ //connectedChan chan bool) { IPAddrs chan string,
	udpListen := getUDPconListener(disconnected)
	if udpListen == nil{
		return
	}
	defer udpListen.Close()

	buffer := make([]byte, 1024)
	for {
		lenOfBuffer, _, err := udpListen.ReadFromUDP(buffer)
		if err != nil {
			log.Println("connListenersd failed") //Bør nok ha error handling her
			disconnected <- true
			udpListen.Close()
			break
		}
		//DecodeJSON returnerer en struct message av meldingen
		msg := DecodeJSON(buffer[:lenOfBuffer]) //&m)

		if msg.Name != GetMyIP(){
			msgRecievedChan <- msg
		}
		select{
		case <-disconnected:
			udpListen.Close()
		default:
			break
		}
		time.Sleep(10*time.Millisecond)
	}
}

//Oppdaterer IP-listen med den IP'en som sendes inn
func updateIP(timeStampMap map[string]time.Time, IPAddr string) {
	timeStampMap[IPAddr] = time.Now().Add(6 * time.Second)
}

//Sjekker om timeStampen er good, ellers fjerner den IP adressen, bor fikse denne funksjonen etterhvert
func timeStampCheck(timeStampMap map[string]time.Time) string { //MyIP string

	for key, val := range timeStampMap {
		if val.Before(time.Now()) { //&& key != MyIP {
			log.Println("Found disconnect")

			delete(timeStampMap, key)
			return key
		}
	}
	return ""
}

func printAliveList(timeStampChan map[string]time.Time) {
	var ipListAlive []string
	for key, _ := range timeStampChan {
		ipListAlive = append(ipListAlive, key)
	}
	log.Println(ipListAlive)
}

func getUDPconListener(disconnected chan bool) *net.UDPConn{
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		log.Println("connListener failed1")
		disconnected <- true
		return nil
	}
	log.Println(udpAddr)
	udpListen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println("connListener failed2")
		disconnected <- true
	}
	return udpListen
}

func getUDPcon(disconnected chan bool) *net.UDPConn{
	serverAddr, err := net.ResolveUDPAddr("udp",myBIP+":"+port)
	    if err != nil {
	        log.Println("getUDPcon failed1")
	        disconnected <- true
			return nil 
	   }
	
	con, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
        log.Println("getUDPcon failed2")
        disconnected <- true
		return nil
	}
    return con
}

func CheckDisconnection(timeStampMap map[string]time.Time, messageStruct Message) string{
	updateIP(timeStampMap, messageStruct.Name)
	return timeStampCheck(timeStampMap)
	//printAliveList(timeStampMap)
}