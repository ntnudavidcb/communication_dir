package com

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

type Message struct {
	Name string
	Body string
	Time time.Time
}

func CreateJSON(name string, body string, timeSendt time.Time) []byte {
	m := Message{Name: name, Body: body, Time: timeSendt}
	b, err := json.Marshal(m)
	if err != nil {
		log.Println("Failed to encode JSON object")
	}
	return b
}

func DecodeJSON(b []byte) Message {
	var m Message
	err := json.Unmarshal(b, &m)
	if err != nil {
		log.Println("Failed to decode JSON object")
		log.Fatal(err)
	}
	return m
}
func statusUpdater(addr string) {
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
		udpBroadcast.Write(CreateJSON("Alive", "", time.Now()))

		time.Sleep(1000 * time.Millisecond)
		log.Println("Hei")
	}
}
