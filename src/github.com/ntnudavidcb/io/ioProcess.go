package io

import (
	"../config"
	"../driver"
	"../converter"
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
	reserved int
}

//Prøve å gjøre denne privat
var PressedButtons = make(map[int]bool)

func ReadAllButtons(buttonPressed chan int) {
	for i := 0; i < 10; i++ {
		PressedButtons[i] = false
	}
	for {
		//log.Println("PressedButtons: ", PressedButtons)
		//log.Println("ElevState: (floor, dir): ", elevState.floor, elevState.direction)
		//log.Println("Local queue: ", PressedButtons)
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
	}
}

func SetPressedButton(button int){
	PressedButtons[button] = true
}

func SetElevState(floor int, direction int, reserved int) {
	elevState.floor = floor
	elevState.direction = direction
	elevState.reserved = reserved
}

func SetElevStateFloor(floor int){
	elevState.floor = floor
}

func SetElevStateDir(direction int) {
	elevState.direction = direction
}

func SetElevStateReserved(reserved int){
	elevState.reserved = reserved
}

func TurnOffLight() { //elevState *ElevState
	log.Println("TurnOffLight: elevState:", elevState)
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

func RemoveButtonFromPressedButtonList(button int){
	if button == config.NOT_ANY_BUTTON{
		return
	}
	PressedButtons[button] = false
}

func RemoveFromPressedButtonList() { //elevState *ElevState,
	buttonUp, buttonDown, buttonCMD := converter.ConvertDirAndFloorToMapIndex(elevState.floor, elevState.direction)
	log.Println(config.ColG, "Floor, direction: ", elevState.floor, elevState.direction, config.ColN)
	log.Println(config.ColR, "buttonUp: ", buttonUp, "buttonDown: ", buttonDown, "buttonCMD: ", buttonCMD, config.ColN)
	if buttonUp == config.NOT_ANY_BUTTON {
		PressedButtons[buttonDown] = false
		PressedButtons[buttonCMD] = false
	} else if buttonDown == config.NOT_ANY_BUTTON {
		PressedButtons[buttonUp] = false
		PressedButtons[buttonCMD] = false
	} else if buttonCMD == config.NOT_ANY_BUTTON {
		PressedButtons[buttonDown] = false
		PressedButtons[buttonUp] = false
	} else{
		if PressedButtons[buttonUp] == true{
			PressedButtons[buttonUp] = false
			PressedButtons[buttonCMD] = false
		} else if PressedButtons[buttonDown] == true{
			PressedButtons[buttonDown] = false
			PressedButtons[buttonCMD] = false
		} else{
			PressedButtons[buttonCMD] = false
		}
	}
}

func FloorSignalListener(floorReached chan bool) {
	for {
		if config.NOT_ANY_FLOOR != driver.Elev_get_floor_sensor_signal() {
			elevState.floor = driver.Elev_get_floor_sensor_signal()
			floorReached <- true
		}
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

func GetElevStateReserved() int{
	return elevState.reserved
}

func InitElevState(floor int){
	elevState.floor = floor
	elevState.direction = config.DIR_STOP
	elevState.reserved = config.NOT_ANY_BUTTON
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
	} else if nextFloor == elevState.floor && elevState.direction == config.DIR_STOP{
		buttonUp, buttonDown, _ := converter.ConvertDirAndFloorToMapIndex(elevState.floor, elevState.direction)
		if buttonUp != config.NOT_ANY_BUTTON{
			elevState.direction = driver.Elev_set_motor_direction(config.DIR_UP)
		} else if buttonDown != config.NOT_ANY_BUTTON{
			elevState.direction = driver.Elev_set_motor_direction(config.DIR_DOWN)
		} else{
			elevState.direction = driver.Elev_set_motor_direction(config.DIR_STOP)
			log.Println(config.ColR,"Failed in GoToNextFloor",config.ColN)

		}
	} else {
		log.Println(config.ColR,"Direction: STOP, in GoToNextFloor, should not happen",config.ColN)
		elevState.direction = driver.Elev_set_motor_direction(config.DIR_STOP)
	}
	log.Println("GoToNextFloor: Direction: ", elevState.direction)
}

func InitButtonAndFloorListeners(buttonPressed chan int, floorReached chan bool) {
	go ReadAllButtons(buttonPressed)
	go FloorSignalListener(floorReached)
}

func HandleWantedFloorReached() {
	RemoveFromPressedButtonList()
	log.Println("PressedButtons: ", PressedButtons)
	if driver.Elev_get_floor_sensor_signal() != -1{
		elevState.floor = driver.Elev_get_floor_sensor_signal()
	}
	TurnOffLight()

	driver.Elev_set_motor_direction(config.DIR_STOP)
	driver.Elev_set_door_open_lamp(1)
	time.Sleep(1500 * time.Millisecond)
	driver.Elev_set_door_open_lamp(0)
}

func GetPressedButtons() map[int]bool {
	return PressedButtons
}