package com

import (
	"log"
	"net"
	"time"
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
		udpBroadcast.Write([]byte("Not master"))
		time.Sleep(1000 * time.Millisecond)
		log.Println("Hei")
	}
}
