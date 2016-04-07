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
		log.Fatal(err)
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

//SKriver ut paa broadcast at denne maskinen lever, og sender informasjon med
func statusUpdater(addr string, port string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr+port)
	if err != nil {
		log.Fatal(err)
	}
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpBroadcast.Close()

	var elevStatus string

	for {
		elevStatus = findElevStatus()
		udpBroadcast.Write(CreateJSON("Status", elevStatus, time.Now()))

		time.Sleep(1000 * time.Millisecond)
		log.Println("Status updater")
	}
}

func findElevStatus() string {
	return "Ingen hardware informasjon tilgjengelig"
}
