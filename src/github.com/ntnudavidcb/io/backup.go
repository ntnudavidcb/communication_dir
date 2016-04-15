package io

import (
	"../config"
	"encoding/gob"
	"log"
	"os"
	"time"
)


func saveToDrive(PressedButtons map[int]bool, filename string) error {
	file, _ := os.Create(filename)
	defer file.Close()
	enc := gob.NewEncoder(file)
	err := enc.Encode(PressedButtons)
	if err != nil {
		log.Println(config.ColR, "gob.Encode() error: Failed to backup.", config.ColN)
		return err
	}
	return nil
}

func loadFromDrive(PressedButtons map[int]bool, filename string) error {
	if _, err := os.Stat(filename); err == nil {
		log.Println(config.ColG, "Backup file found, processing...", config.ColN)

		file, _ := os.Open(filename)
		dec := gob.NewDecoder(file)
		err := dec.Decode(&PressedButtons)
		if err != nil {
			log.Println(config.ColR, "gob.Decode() error: Failed to decode file.", err, config.ColN)
		}
	}
	log.Println(config.ColG, "runBackup: backup1 map:", PressedButtons)
	return nil
}

func RunBackup(){
	const filename = "elevatorBackup"

	backup := make(map[int]bool)
	loadFromDrive(backup, filename)
	log.Println(config.ColR, "runBackup: backup map:", backup)

	log.Println(config.ColM,"runBackup: backup:", backup)
	//if !isEmptyMap(backup) {
	for order := 0; order < 10; order++ {
		if value, ok := backup[order]; ok {
			log.Println(config.ColM, "runBackup: Index: ", order)
			SetPressedButton(order, value)
		} else {
			SetPressedButton(order, false)
		}
	}
	//}

	go func(){
		for {
			<-takeBackup

			dummyLocalQueue := make(map[int]bool)

			for elem := 0; elem < 10; elem++ {
				dummyLocalQueue[elem] = PressedButtons[elem]
			}
			for elem := 0; elem < 6; elem++ {
				dummyLocalQueue[elem] = false
			}
			if err := saveToDrive(dummyLocalQueue, filename); err != nil {
				log.Println(config.ColR, err, config.ColN)
			}
			time.Sleep(10*time.Millisecond)
		}
	}()
}