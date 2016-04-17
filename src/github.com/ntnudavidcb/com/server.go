package com

import (
	. "../config"
	"log"
	"net"
	"time"
)

var myBIP string
var myIP string
var port string
var timeStampMap = make(map[string]time.Time) //Holde styr pa timestamps paa IP adressene som blir sendt inn

//Kjorer hele tiden, holder hele nettverkssystemet oppe og styrer hele showet
func Server(ipListChannel chan []string, sendAliveMessage chan Message, msgRecievedChan chan Message, disconnected chan bool) {
	port = "20010"
	myIP = GetMyIP()
	myBIP = GetBIP(myIP)

	udpBroadcast := getUDPcon(disconnected)

	defer udpBroadcast.Close()

	go connListener(msgRecievedChan, disconnected)

	for {
		msg := <-sendAliveMessage
		_, err := udpBroadcast.Write(CreateJSON(msg))
		if err != nil {
			log.Println("sendAliveMessage failed")
			disconnected <- true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func connListener(msgRecievedChan chan Message, disconnected chan bool) {
	udpListen := getUDPconListener(disconnected)
	if udpListen == nil {
		return
	}
	defer udpListen.Close()

	buffer := make([]byte, 1024)
	for {
		lenOfBuffer, _, err := udpListen.ReadFromUDP(buffer)
		if err != nil {
			log.Println("connListenersd failed")
			disconnected <- true
			udpListen.Close()
			break
		}

		msg := DecodeJSON(buffer[:lenOfBuffer])

		if msg.Name != GetMyIP() {
			msgRecievedChan <- msg
		}
		select {
		case <-disconnected:
			udpListen.Close()
		default:
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

//Oppdaterer IP-listen med den IP'en som sendes inn
func UpdateIP(IPAddr string) {
	timeStampMap[IPAddr] = time.Now().Add(6 * time.Second)
}

func timeStampCheck(timeStampMap map[string]time.Time) string {

	for key, val := range timeStampMap {
		if val.Before(time.Now()) {
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

func getUDPconListener(disconnected chan bool) *net.UDPConn {
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

func getUDPcon(disconnected chan bool) *net.UDPConn {
	serverAddr, err := net.ResolveUDPAddr("udp", myBIP+":"+port)
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

func CheckForDisconnect() string {
	for key, val := range timeStampMap {
		if val.Before(time.Now()) { //Disconnected
			log.Println(ColB, "Elevator with IP ", key, " disconnected from the network or might be stuck between floors.", ColN)
			delete(timeStampMap, key)
			return key
		}
	}
	return ""
}
