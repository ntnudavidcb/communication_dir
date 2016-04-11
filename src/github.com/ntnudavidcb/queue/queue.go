package queue

import (
	//"../config"
	//"../driver"
	//"time"
	"../costFunc"
	"../io"
	//"log"
)

var localQueue = [10]bool{}
var costQueue = [10]int{}
var dummyLocalQueue = [10]bool{}
var sortedQueue = [6]int{}

/*Organized as follows: 
UP_1, UP_2, UP_3, DOWN_4, DOWN_3, DOWN_2, 
CMD_1, CMD_2, CMD_3,CMD_4*/

func CheckOrder() bool {
	buttonPressed1, buttonPressed2 := io.ConvertDirAndFloorToMapIndex()
	return inLocalQueue(buttonPressed1) || inLocalQueue(buttonPressed2)
}

func CheckQueue() bool {
	buttonPressed1, buttonPressed2 := io.ConvertDirAndFloorToMapIndex()
	
}

func inLocalQueue(buttonPressed int) bool {
	if buttonPressed == -1 {
		return false
	}
	return localQueue[buttonPressed]
}

func UpdateQueueWithButton(buttonPressed int) {
	localQueue[buttonPressed] = true
}

func UpdateQueueButtonPushed(buttonPushed int) {
	UpdateQueueWithButton(buttonPushed)
	updateCostQueue()
	SortQueue()
}
func UpdateQueueFloorReached() {
	updateCostQueue()
	SortQueue()
}

func InitQueue() {
	SortQueue()
}

func updateCostQueue() {
	currentFloor, currentDir, _ := io.GetElevState()
	for button := 0; button < 10; button++ {
		costQueue[button] = costFunc.CostFunc(currentDir, currentFloor, button)
	}
}

func convertButtonCMD(buttonPressed int) (int, int) {
	if buttonPressed > 8 {
		return 5, 5
	} else if buttonPressed > 7 {
		return 3, 4
	} else if buttonPressed > 6 {
		return 1, 2
	} else if buttonPressed > 5 {
		return 0, 0
	} else {
		return buttonPressed, buttonPressed
	}
}

func RemoveFromQueue(pressedButtons map[int]bool) {
	for key, _ := range pressedButtons {
		localQueue[key] = pressedButtons[key]
	}
	updateCostQueue()
	SortQueue()
}

func GetNextOrder() int {
	return sortedQueue[0]
}

func SortQueue() {
	minButton := -1
	for elem := 0; elem < 10; elem++ {
		dummyLocalQueue[elem] = localQueue[elem]
	}
	for i := 0; i < 6; i++ {
		min := 100
		for j := 0; j < 10; j++ {
			if dummyLocalQueue[j] {
				if min > costQueue[j] {
					min = costQueue[j]
					minButton = j
				}
			}
		}
		if minButton != -1 {
			dummyLocalQueue[minButton] = false
		}
		//log.Println("MinButton: ", minButton)
		sortedQueue[i] = convertButtonToFloor(minButton)
	}
}

func convertButtonToFloor(button int) int {
	if button == costFunc.CMD_4 || button == costFunc.DOWN_4 {
		return 3
	} else if button == costFunc.CMD_3 || button == costFunc.DOWN_3 || button == costFunc.UP_3 {
		return 2
	} else if button == costFunc.CMD_2 || button == costFunc.DOWN_2 || button == costFunc.UP_2 {
		return 1
	} else if button == costFunc.CMD_1 || button == costFunc.UP_1 {
		return 0
	}
	return -1
}
