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

var elevState struct {
	floor int
	direction int
}

var PressedButtons = make(map[int]bool)

func ReadAllButtons(buttonPressed chan bool){
	for{
		if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_1){
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_1, 1)
			PressedButtons[6] = true
			buttonPressed <- true
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_2){
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_2, 1)
			PressedButtons[7] = true
			buttonPressed <- true
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_3){
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_3, 1)
			PressedButtons[8] = true
			buttonPressed <- true
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_4){
			driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_4, 1)
			PressedButtons[9] = true
			buttonPressed <- true
		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_1) {
			driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_1, 1)
			PressedButtons[0] = true
			buttonPressed <- true
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_2){
			driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_2, 1)
			PressedButtons[1] = true
			buttonPressed <- true
		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_2){
			driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_2, 1)
			PressedButtons[2] = true
			buttonPressed <- true
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_3){
			driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_3, 1)
			PressedButtons[3] = true
			buttonPressed <- true

		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_3){
			driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_3, 1)
			PressedButtons[4] = true
			buttonPressed <- true
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_4){
			driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_4, 1)
			PressedButtons[5] = true
			buttonPressed <- true
		}
		//log.Println("Local queue: ", PressedButtons)
	}
}

func SetElevState(floor int, direction int){
	elevState.floor = floor
	elevState.direction = direction
}

func TurnOffLight(){//elevState *ElevState
	if elevState.direction == config.DIR_UP{
		driver.Elev_set_button_lamp(config.BTN_UP, elevState.floor, 0)
		driver.Elev_set_button_lamp(config.BTN_COMMAND, elevState.floor, 0)
	}else if elevState.direction == config.DIR_DOWN{
		driver.Elev_set_button_lamp(config.BTN_DOWN, elevState.floor, 0)
		driver.Elev_set_button_lamp(config.BTN_COMMAND, elevState.floor, 0)
	}
}


func RemoveFromPressedButtonList(){//elevState *ElevState,
	index1, index2 := ConvertDirAndFloorToMapIndex()
	PressedButtons[index1] = false
	PressedButtons[index2] = false
}

func ConvertDirAndFloorToMapIndex() (int, int){ //elevState *ElevState
	log.Println(elevState.floor)
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

func ConvertMapIndexToFloor(mapIndex int) int{
	if mapIndex == 6{
		return 0
	}else if mapIndex == 7{
		return 1
	}else if mapIndex == 8{
		return 2
	}else if mapIndex == 9{
		return 3
	}else if mapIndex > 4{
		return 3
	}else if mapIndex >2{
		return 2
	}else if mapIndex > 0{
		return 1
	}else{
		return 0
	}
}

func floorSignalListener(floorReached chan bool) {
	for{
		if -1 != driver.Elev_get_floor_sensor_signal(){
			floorReached <- true
		}
	}
}

func GetElevState() (int, int){
	return elevState.floor, elevState.direction
}


func Testrun2(floorReached chan bool, buttonPressed chan bool, nextFloor chan int){
	for i := 0; i < 10; i++ {
		PressedButtons[i] = false
	}
	

	go ReadAllButtons(buttonPressed)
	go floorSignalListener(floorReached)
	log.Println(config.ColC, "Test Run 2 Initialized", config.ColN)

	var nextStopFloor int
	for{
		nextStopFloor =<- nextFloor
		if nextStopFloor == -1{
			driver.Elev_set_motor_direction(config.DIR_STOP)
			continue
		}
		goTowardsFloor(nextStopFloor)
	}
}

func goTowardsFloor(nextFloor int){
	if nextFloor > elevState.floor{
		driver.Elev_set_motor_direction(config.DIR_UP)
	}else if nextFloor < elevState.floor{
		driver.Elev_set_motor_direction(config.DIR_DOWN)
	}else{
		driver.Elev_set_motor_direction(config.DIR_STOP)
	}

}




func WantedFloorReached(){
	log.Println("JJJJJJAAAAA")
	RemoveFromPressedButtonList()
	TurnOffLight()
	elevState.direction = driver.Elev_set_motor_direction(config.DIR_STOP)
	time.Sleep(1000 * time.Millisecond)
}