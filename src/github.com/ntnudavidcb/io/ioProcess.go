package io

import (
	"../driver"
	"../config"
	//"../queue"
	"time"
	"log"
)

const(
	UP_1 int = iota
	DOWN_2
	UP_2
	DOWN_3
	UP_3
	DOWN_4
	CMD_1 
	CMD_2
	CMD_3
	CMD_4
)

type ElevState struct {
	floor int
	direction int
}



func ReadAllButtons(PressedButtons map[int]bool){
	for{
		if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_1){
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_1, 1)
			PressedButtons[6] = true
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_2){
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_2, 1)
			PressedButtons[7] = true
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_3){
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_3, 1)
			PressedButtons[8] = true
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_4){
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_4, 1)
			PressedButtons[9] = true

		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_1) {
			driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_1, 1)
			PressedButtons[0] = true
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_2){
			driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_2, 1)
			PressedButtons[1] = true
		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_2){
			driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_2, 1)
			PressedButtons[2] = true
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_3){
			driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_3, 1)
			PressedButtons[3] = true
		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_3){
			driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_3, 1)
			PressedButtons[4] = true
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_4){
			driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_4, 1)
			PressedButtons[5] = true
		}
		log.Println("Local queue: ", PressedButtons)
	}
}

func TurnOffLight(elevState *ElevState){
	if elevState.direction == config.DIR_UP{
		driver.Elev_set_button_lamp(config.BTN_UP, elevState.floor, 0)
		driver.Elev_set_button_lamp(config.BTN_COMMAND, elevState.floor, 0)
	}else if elevState.direction == config.DIR_DOWN{
		driver.Elev_set_button_lamp(config.BTN_DOWN, elevState.floor, 0)
		driver.Elev_set_button_lamp(config.BTN_COMMAND, elevState.floor, 0)
	}
}


func RemoveFromPressedButtonList(elevState *ElevState, PressedButtons map[int]bool){
	index1, index2 := convertDirAndFloorToMapIndex(elevState)
	PressedButtons[index1] = false
	PressedButtons[index2] = false
}

func convertDirAndFloorToMapIndex(elevState *ElevState) (int, int){
	if elevState.floor == config.FLOOR_1{
		return UP_1, CMD_1;
	}else if elevState.floor == config.FLOOR_2 && elevState.direction == config.DIR_DOWN{
		return DOWN_2, CMD_2
	}else if elevState.floor == config.FLOOR_2 && elevState.direction == config.DIR_UP{
		return UP_2, CMD_2
	}else if elevState.floor == config.FLOOR_3 && elevState.direction == config.DIR_DOWN{
		return DOWN_3, CMD_3
	}else if elevState.floor == config.FLOOR_3 && elevState.direction == config.DIR_UP{
		return UP_3, CMD_3
	}else if elevState.floor == config.FLOOR_4{
		return DOWN_4, CMD_4
	}
	return -1, -1

}

func CheckOrder(elevState *ElevState, PressedButtons map[int]bool) bool{
	index1, index2 := convertDirAndFloorToMapIndex(elevState)
	return PressedButtons[index2] || PressedButtons[index1]

}


func Testrun2(){
	PressedButtons :=make(map[int]bool)
	for i := 0; i < 10; i++ {
		PressedButtons[i] = false
	}
	var s ElevState
	go ReadAllButtons(PressedButtons)
	elevState := &s
	elevState.direction = config.DIR_UP
	log.Println(config.ColC, "Test Run 2 Initialized", config.ColN)

	for{
		for driver.Elev_get_floor_sensor_signal() != config.FLOOR_1{
			elevState.direction = driver.Elev_set_motor_direction(config.DIR_DOWN)
			elevState.floor = driver.Elev_get_floor_sensor_signal()
			if CheckOrder(elevState, PressedButtons){
				RemoveFromPressedButtonList(elevState, PressedButtons)
				TurnOffLight(elevState)
				elevState.direction = driver.Elev_set_motor_direction(config.DIR_STOP)
				time.Sleep(1000 * time.Millisecond)
			}
		}
		log.Println("Reached 1st floor")
		elevState.direction = driver.Elev_set_motor_direction(config.DIR_STOP)
		time.Sleep(1000 * time.Millisecond)
		for driver.Elev_get_floor_sensor_signal() != config.FLOOR_4{
			elevState.direction = driver.Elev_set_motor_direction(config.DIR_UP)
			elevState.floor = driver.Elev_get_floor_sensor_signal()
			if CheckOrder(elevState, PressedButtons){
				RemoveFromPressedButtonList(elevState, PressedButtons)
				TurnOffLight(elevState)
				elevState.direction = driver.Elev_set_motor_direction(config.DIR_STOP)
				time.Sleep(1000 * time.Millisecond)
			}

		}
		elevState.direction = driver.Elev_set_motor_direction(config.DIR_STOP)
		log.Println("Reached 4th floor")
		time.Sleep(1000 * time.Millisecond)

	}
}
