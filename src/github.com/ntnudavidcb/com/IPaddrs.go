package com

import (
	"log"
	"net"
	"strings"
)

func GetMyIP() string {
	allIPs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("IP receiving errors!\n")
		return ""
	}
	return strings.Split(allIPs[1].String(), "/")[0]
}

func GetBIP(MyIP string) string {
	IP := strings.Split(MyIP, ".")
	return IP[0] + "." + IP[1] + "." + IP[2] + ".255"
}

func TimeStampCheck(list chan map[string]time.Time, deletedIP chan string, MyIP string) {
	var IPlist map[string]time.Time
	for {
		IPlist = <-list
		for key, val := range IPlist {
			if val.Before(time.Now()) && key != MyIP {
				deletedIP <- key
				delete(IPlist, key)
				break
			}
		}
		list <- IPlist
	}
}

//Add code under to add a timer of 1200 milliseconds to a connection
/*
IPlist[newMsg.from]=time.Now().Add(1200*time.Millisecond)
aliveMapChan<-IPlist
*/
