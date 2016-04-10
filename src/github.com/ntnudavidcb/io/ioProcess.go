package io

import (
	"../config"
	"../driver"
	//"../queue"
	"log"
	"time"
)

const (
	UP_1   = 0
	UP_2   = 1
	UP_3   = 2
	DOWN_4 = 3
	DOWN_3 = 4
	DOWN_2 = 5
	CMD_1  = 6
	CMD_2  = 7
	CMD_3  = 8
	CMD_4  = 9
)

var elevState struct {
	floor     int
	direction int
}

var PressedButtons = make(map[int]bool)

func ReadAllButtons(buttonPressed chan int) {
	for {
		//log.Println("PressedButtons: ", PressedButtons)
		//log.Println("ElevState: (floor, dir): ", elevState.floor, elevState.direction)
		if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_1) {
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_1, 1)
			PressedButtons[6] = true
			buttonPressed <- 6

		} else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_2) {
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_2, 1)
			PressedButtons[7] = true
			buttonPressed <- 7
		} else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_3) {
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_3, 1)
			PressedButtons[8] = true
			buttonPressed <- 8
		} else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_4) {
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_4, 1)
			PressedButtons[9] = true
			buttonPressed <- 9
		} else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_1) {
			driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_1, 1)
			PressedButtons[0] = true
			buttonPressed <- 0
		} else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_2) {
			driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_2, 1)
			PressedButtons[1] = true
			buttonPressed <- 1
		} else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_3) {
			driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_3, 1)
			PressedButtons[2] = true
			buttonPressed <- 2
		} else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_4) {
			driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_4, 1)
			PressedButtons[3] = true
			buttonPressed <- 3
		} else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_3) {
			driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_3, 1)
			PressedButtons[4] = true
			buttonPressed <- 4
		} else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_2) {
			driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_2, 1)
			PressedButtons[5] = true
			buttonPressed <- 5
		}
		//log.Println("Local queue: ", PressedButtons)
	}
}

func SetElevState(floor int, direction int) {
	elevState.floor = floor
	elevState.direction = direction
}

func SetElevStateDir(direction int) {
	elevState.direction = direction
}

func TurnOffLight() { //elevState *ElevState
	if elevState.direction == config.DIR_UP {
		driver.Elev_set_button_lamp(config.BTN_UP, elevState.floor, 0)
		driver.Elev_set_button_lamp(config.BTN_COMMAND, elevState.floor, 0)
		if elevState.floor == config.FLOOR_4 {
			driver.Elev_set_button_lamp(config.BTN_DOWN, elevState.floor, 0)
		}
	} else if elevState.direction == config.DIR_DOWN {
		driver.Elev_set_button_lamp(config.BTN_DOWN, elevState.floor, 0)
		driver.Elev_set_button_lamp(config.BTN_COMMAND, elevState.floor, 0)
		if elevState.floor == config.FLOOR_1 {
			driver.Elev_set_button_lamp(config.BTN_UP, elevState.floor, 0)
		}
	}
}

func RemoveFromPressedButtonList() { //elevState *ElevState,
	index1, index2 := ConvertDirAndFloorToMapIndex()
	log.Println("Index1, Index2: ", index1, index2)
	PressedButtons[index1] = false
	PressedButtons[index2] = false
}

func ConvertDirAndFloorToMapIndex() (int, int) { //elevState *ElevState
	if elevState.floor == config.FLOOR_1 {
		return UP_1, CMD_1
	} else if elevState.floor == config.FLOOR_2 && elevState.direction == config.DIR_DOWN {
		return DOWN_2, CMD_2
	} else if elevState.floor == config.FLOOR_2 && elevState.direction == config.DIR_UP {
		return UP_2, CMD_2
	} else if elevState.floor == config.FLOOR_3 && elevState.direction == config.DIR_DOWN {
		return DOWN_3, CMD_3
	} else if elevState.floor == config.FLOOR_3 && elevState.direction == config.DIR_UP {
		return UP_3, CMD_3
	} else if elevState.floor == config.FLOOR_4 {
		return DOWN_4, CMD_4
	}
	return -1, -1
}

func ConvertMapIndexToFloor(mapIndex int) int {
	if mapIndex == 6 {
		return 0
	} else if mapIndex == 7 {
		return 1
	} else if mapIndex == 8 {
		return 2
	} else if mapIndex == 9 {
		return 3
	} else if mapIndex == 3 {
		return 3
	} else if mapIndex == 2 || mapIndex == 4 {
		return 2
	} else if mapIndex == 1 || mapIndex == 5 {
		return 1
	} else {
		return 0
	}
}

func floorSignalListener(floorReached chan bool) {
	for {
		if -1 != driver.Elev_get_floor_sensor_signal() {
			elevState.floor = driver.Elev_get_floor_sensor_signal()
			floorReached <- true
		}
	}
}

func GetElevState() (int, int) {
	return elevState.floor, elevState.direction
}

func GoToNextFloor(nextFloor int) {
	if nextFloor == -1 {
		elevState.direction = driver.Elev_set_motor_direction(config.DIR_STOP)
	} else if nextFloor > elevState.floor {
		elevState.direction = driver.Elev_set_motor_direction(config.DIR_UP)
	} else if nextFloor < elevState.floor {
		elevState.direction = driver.Elev_set_motor_direction(config.DIR_DOWN)
	} else if nextFloor == elevState.floor && elevState.direction == config.DIR_DOWN {
		elevState.direction = driver.Elev_set_motor_direction(config.DIR_UP)
	} else if nextFloor == elevState.floor && elevState.direction == config.DIR_UP {
		elevState.direction = driver.Elev_set_motor_direction(config.DIR_DOWN)
	} else {
		log.Println("Direction: STOP, in GoToNextFloor, should not happen")
		elevState.direction = driver.Elev_set_motor_direction(config.DIR_STOP)
	}

}

func InitListeners(buttonPressed chan int, floorReached chan bool) {
	for i := 0; i < 10; i++ {
		PressedButtons[i] = false
	}
	go ReadAllButtons(buttonPressed)
	go floorSignalListener(floorReached)
}

func WantedFloorReached() {
	log.Println("WantedFloorReached")
	RemoveFromPressedButtonList()
	TurnOffLight()

	driver.Elev_set_motor_direction(config.DIR_STOP)
	time.Sleep(1000 * time.Millisecond)
}

func GetPressedButtons() map[int]bool {
	return PressedButtons
}
