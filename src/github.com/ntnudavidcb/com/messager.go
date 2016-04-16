package com

import (
	"encoding/json"
	"log"
	//"net"
	"time"
)

type Message struct {
	Name         string //IP
	ButtonPushed int
	Floor        int
	Direction    int
	Reserved     int
	OrderTaken   int
	Time         time.Time
}

func CreateJSON(msg Message) []byte {
	b, err := json.Marshal(msg)
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
func SendMessage(msg Message, disconnected chan bool) {
	udpCon := getUDPcon(disconnected)
	defer udpCon.Close()
	udpCon.Write(CreateJSON(msg))
}
