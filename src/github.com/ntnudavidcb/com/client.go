package com

import (
	"log"
	"net"
	"time"
	//"encoding/json"
)

func BroadcastUdp(addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer udpBroadcast.Close()

	for {
		udpBroadcast.Write(CreateJSON("Alive", "shit", time.Now()))

		time.Sleep(1000 * time.Millisecond)
		log.Println("Hei")
	}
}
