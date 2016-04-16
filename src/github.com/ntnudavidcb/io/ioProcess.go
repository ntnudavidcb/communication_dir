package io

import (
	. "../config"
	"../converter"
	"../driver"
	"log"
	"time"
)

var elevState struct {
	floor     int
	direction int
	reserved  int
}

//Prøve å gjøre denne privat
var PressedButtons = make(map[int]bool)

var takeBackup = make(chan bool, CMD_4+1)

func ReadAllButtons(buttonPressed chan int) {
	for {
		if driver.Elev_get_button_signal(BTN_COMMAND, FLOOR_1) {
			driver.Elev_set_button_lamp(BTN_COMMAND, FLOOR_1, ON)
			PressedButtons[CMD_1] = true
			buttonPressed <- CMD_1
		} else if driver.Elev_get_button_signal(BTN_COMMAND, FLOOR_2) {
			driver.Elev_set_button_lamp(BTN_COMMAND, FLOOR_2, ON)
			PressedButtons[CMD_2] = true
			buttonPressed <- CMD_2
		} else if driver.Elev_get_button_signal(BTN_COMMAND, FLOOR_3) {
			driver.Elev_set_button_lamp(BTN_COMMAND, FLOOR_3, ON)
			PressedButtons[CMD_3] = true
			buttonPressed <- CMD_3
		} else if driver.Elev_get_button_signal(BTN_COMMAND, FLOOR_4) {
			driver.Elev_set_button_lamp(BTN_COMMAND, FLOOR_4, ON)
			PressedButtons[CMD_4] = true
			buttonPressed <- CMD_4
		} else if driver.Elev_get_button_signal(BTN_UP, FLOOR_1) {
			driver.Elev_set_button_lamp(BTN_UP, FLOOR_1, ON)
			PressedButtons[UP_1] = true
			buttonPressed <- UP_1
		} else if driver.Elev_get_button_signal(BTN_UP, FLOOR_2) {
			driver.Elev_set_button_lamp(BTN_UP, FLOOR_2, ON)
			PressedButtons[UP_2] = true
			buttonPressed <- UP_2
		} else if driver.Elev_get_button_signal(BTN_UP, FLOOR_3) {
			driver.Elev_set_button_lamp(BTN_UP, FLOOR_3, ON)
			PressedButtons[UP_3] = true
			buttonPressed <- UP_3
		} else if driver.Elev_get_button_signal(BTN_DOWN, FLOOR_4) {
			driver.Elev_set_button_lamp(BTN_DOWN, FLOOR_4, ON)
			PressedButtons[DOWN_4] = true
			buttonPressed <- DOWN_4
		} else if driver.Elev_get_button_signal(BTN_DOWN, FLOOR_3) {
			driver.Elev_set_button_lamp(BTN_DOWN, FLOOR_3, ON)
			PressedButtons[DOWN_3] = true
			buttonPressed <- DOWN_3
		} else if driver.Elev_get_button_signal(BTN_DOWN, FLOOR_2) {
			driver.Elev_set_button_lamp(BTN_DOWN, FLOOR_2, ON)
			PressedButtons[DOWN_2] = true
			buttonPressed <- DOWN_2
		} else {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		takeBackup <- true
		time.Sleep(10 * time.Millisecond)
	}
}

func SetPressedButton(button int, value bool) {
	PressedButtons[button] = value
}

func SetElevState(floor int, direction int, reserved int) {
	elevState.floor = floor
	elevState.direction = direction
	elevState.reserved = reserved
}

func SetElevStateFloor(floor int) {
	elevState.floor = floor
}

func SetElevStateDir(direction int) {
	elevState.direction = direction
}

func SetElevStateReserved(reserved int) {
	elevState.reserved = reserved
}

func TurnOffLight() {
	log.Println("TurnOffLight: elevState:", elevState)
	if elevState.direction == DIR_UP {
		driver.Elev_set_button_lamp(BTN_UP, elevState.floor, OFF)
		driver.Elev_set_button_lamp(BTN_COMMAND, elevState.floor, OFF)
		if elevState.floor == FLOOR_4 {
			driver.Elev_set_button_lamp(BTN_DOWN, elevState.floor, OFF)
		}
	} else if elevState.direction == DIR_DOWN {
		driver.Elev_set_button_lamp(BTN_DOWN, elevState.floor, OFF)
		driver.Elev_set_button_lamp(BTN_COMMAND, elevState.floor, OFF)
		if elevState.floor == FLOOR_1 {
			driver.Elev_set_button_lamp(BTN_UP, elevState.floor, OFF)
		}
	}
}

func UpdateLightsWithRecievedMessage(floor int, direction int) {
	if direction == DIR_UP {
		driver.Elev_set_button_lamp(BTN_UP, floor, OFF)
		if floor == FLOOR_4 {
			driver.Elev_set_button_lamp(BTN_DOWN, floor, OFF)
		}
	} else if direction == DIR_DOWN {
		driver.Elev_set_button_lamp(BTN_DOWN, floor, OFF)
		if floor == FLOOR_1 {
			driver.Elev_set_button_lamp(BTN_UP, floor, OFF)
		}
	}
}

func RemoveButtonFromPressedButtonList(button int) {
	if button == NOT_ANY_BUTTON {
		return
	}
	PressedButtons[button] = false
	takeBackup <- true
}

func RemoveFromPressedButtonList() {
	buttonUp, buttonDown, buttonCMD := converter.ConvertDirAndFloorToMapIndex(elevState.floor, elevState.direction)
	log.Println(ColG, "Floor, direction: ", elevState.floor, elevState.direction, ColN)
	log.Println(ColR, "buttonUp: ", buttonUp, "buttonDown: ", buttonDown, "buttonCMD: ", buttonCMD, ColN)
	if buttonUp == NOT_ANY_BUTTON {
		PressedButtons[buttonDown] = false
		PressedButtons[buttonCMD] = false
	} else if buttonDown == NOT_ANY_BUTTON {
		PressedButtons[buttonUp] = false
		PressedButtons[buttonCMD] = false
	} else if buttonCMD == NOT_ANY_BUTTON {
		PressedButtons[buttonDown] = false
		PressedButtons[buttonUp] = false
	} else {
		if PressedButtons[buttonUp] == true {
			PressedButtons[buttonUp] = false
			PressedButtons[buttonCMD] = false
		} else if PressedButtons[buttonDown] == true {
			PressedButtons[buttonDown] = false
			PressedButtons[buttonCMD] = false
		} else {
			PressedButtons[buttonCMD] = false
		}
	}
	takeBackup <- true
}

func FloorSignalListener(floorReached chan bool) {
	for {
		if NOT_ANY_FLOOR != driver.Elev_get_floor_sensor_signal() {
			elevState.floor = driver.Elev_get_floor_sensor_signal()
			floorReached <- true
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func GetElevState() (int, int, int) {
	return elevState.floor, elevState.direction, elevState.reserved
}

func GetElevStateFloor() int {
	return elevState.floor
}

func GetElevStateDir() int {
	return elevState.direction
}

func GetElevStateReserved() int {
	return elevState.reserved
}

func InitElevState(floor int) {
	elevState.floor = floor
	elevState.direction = DIR_STOP
	elevState.reserved = NOT_ANY_BUTTON
}

func GoToNextFloor(nextFloor int) int {
	if nextFloor == -1 {
		elevState.direction = driver.Elev_set_motor_direction(DIR_STOP)
	} else if nextFloor > elevState.floor {
		elevState.direction = driver.Elev_set_motor_direction(DIR_UP)
	} else if nextFloor < elevState.floor {
		elevState.direction = driver.Elev_set_motor_direction(DIR_DOWN)
	} else if nextFloor == elevState.floor && elevState.direction == DIR_DOWN {
		elevState.direction = driver.Elev_set_motor_direction(DIR_UP)
	} else if nextFloor == elevState.floor && elevState.direction == DIR_UP {
		elevState.direction = driver.Elev_set_motor_direction(DIR_DOWN)
	} else if nextFloor == elevState.floor && elevState.direction == DIR_STOP {
		buttonUp, buttonDown, _ := converter.ConvertDirAndFloorToMapIndex(elevState.floor, elevState.direction)
		if buttonUp != NOT_ANY_BUTTON {
			elevState.direction = driver.Elev_set_motor_direction(DIR_UP)
		} else if buttonDown != NOT_ANY_BUTTON {
			elevState.direction = driver.Elev_set_motor_direction(DIR_DOWN)
		} else {
			elevState.direction = driver.Elev_set_motor_direction(DIR_STOP)
			log.Println(ColR, "Failed in GoToNextFloor", ColN)

		}
	} else {
		log.Println(ColR, "Direction: STOP, in GoToNextFloor, should not happen", ColN)
		elevState.direction = driver.Elev_set_motor_direction(DIR_STOP)
	}
	log.Println("GoToNextFloor: Direction: ", elevState.direction)
	return elevState.direction
}

func InitButtonAndFloorListeners(buttonPressed chan int, floorReached chan bool) {
	go ReadAllButtons(buttonPressed)
	go FloorSignalListener(floorReached)
}

func HandleWantedFloorReached() {
	RemoveFromPressedButtonList()
	log.Println("PressedButtons: ", PressedButtons)
	if driver.Elev_get_floor_sensor_signal() != NOT_ANY_FLOOR {
		elevState.floor = driver.Elev_get_floor_sensor_signal()
	}
	TurnOffLight()

	driver.Elev_set_motor_direction(DIR_STOP)
	driver.Elev_set_door_open_lamp(ON)
	time.Sleep(1500 * time.Millisecond)
	driver.Elev_set_door_open_lamp(OFF)
}

func GetPressedButtons() map[int]bool {
	return PressedButtons
}

func isEmptyMap(buttonMap map[int]bool) bool {
	for _, order := range buttonMap {
		if order {
			return false
		}
	}
	return true
}
