package io

import (
	. "../config"
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
		log.Println(ColR, "gob.Encode() error: Failed to backup.", ColN)
		return err
	}
	return nil
}

func loadFromDrive(PressedButtons map[int]bool, filename string) error {
	if _, err := os.Stat(filename); err == nil {
		log.Println(ColG, "Backup file found, processing...", ColN)

		file, _ := os.Open(filename)
		dec := gob.NewDecoder(file)
		err := dec.Decode(&PressedButtons)
		if err != nil {
			log.Println(ColR, "gob.Decode() error: Failed to decode file.", err, ColN)
		}
	}
	log.Println(ColG, "runBackup: backup1 map:", PressedButtons)
	return nil
}

func RunBackup() {
	const filename = "elevatorBackup"

	backup := make(map[int]bool)
	loadFromDrive(backup, filename)
	log.Println(ColR, "runBackup: backup map:", backup)

	log.Println(ColM, "runBackup: backup:", backup)
	//if !isEmptyMap(backup) {
	for order := UP_1; order < CMD_4+1; order++ {
		if value, ok := backup[order]; ok {
			log.Println(ColM, "runBackup: Index: ", order)
			SetPressedButton(order, value)
		} else {
			SetPressedButton(order, false)
		}
	}
	//}

	go func() {
		for {
			<-takeBackup

			dummyLocalQueue := make(map[int]bool)

			for elem := UP_1; elem < CMD_4+1; elem++ {
				dummyLocalQueue[elem] = PressedButtons[elem]
			}
			for elem := UP_1; elem < CMD_1; elem++ {
				dummyLocalQueue[elem] = false
			}
			if err := saveToDrive(dummyLocalQueue, filename); err != nil {
				log.Println(ColR, err, ColN)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
}
