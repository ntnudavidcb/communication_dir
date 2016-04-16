package queue

import (
	. "../config"
	"../converter"
	"../costFunc"
	"log"
)

var localQueue = [10]bool{}
var costQueue = [10]int{}
var dummyLocalQueue = [10]bool{}
var sortedQueue = [10]int{}

func SetMyIP(IP string) {
	costFunc.MyIP = IP
}

func CheckOrder(floor int, direction int) bool {
	buttonUp, buttonDown, buttonCMD := converter.ConvertDirAndFloorToMapIndex(floor, direction) //costFunc.ElevStateMap[costFunc.MyIP].Floor, costFunc.ElevStateMap[costFunc.MyIP].Direction
	log.Println("CheckOrder: buttonUp, buttonDown, buttonCMD: ", buttonUp, buttonDown, buttonCMD)
	log.Println("CheckOrder: localQueue: ", localQueue)
	return inLocalQueue(buttonUp) || inLocalQueue(buttonDown) || inLocalQueue(buttonCMD)
}

func CheckUpOrDownButton() int {
	buttonUp, buttonDown, _ := converter.ConvertDirAndFloorToMapIndex(costFunc.ElevStateMap[costFunc.MyIP].Floor, costFunc.ElevStateMap[costFunc.MyIP].Direction)
	log.Println("CheckUpOrDownButton: Up, down: ", buttonUp, buttonDown)
	if inLocalQueue(buttonUp) {
		return BTN_UP
	} else if inLocalQueue(buttonDown) {
		return BTN_DOWN
	} else {
		return BTN_COMMAND
	}
}

func inLocalQueue(buttonPressed int) bool {
	if buttonPressed == NOT_ANY_BUTTON {
		return false
	}
	return localQueue[buttonPressed]
}

func AddButtonToQueue(buttonPressed int) {
	if inLocalQueue(buttonPressed) {
		return
	}
	localQueue[buttonPressed] = true
}

func UpdateQueue() {
	updateCostQueue()
	sortQueue()
}

func updateCostQueue() {
	currentFloor := costFunc.ElevStateMap[costFunc.MyIP].Floor
	currentDir := costFunc.ElevStateMap[costFunc.MyIP].Direction
	for button := UP_1; button < CMD_4+1; button++ {
		costQueue[button] = costFunc.CostFunc(currentDir, currentFloor, button)
	}
}

func RemoveButtonFromQueue(button int) {
	if button == NOT_ANY_BUTTON {
		return
	}
	localQueue[button] = false
}

func InitQueue() {
	for button := UP_1; button < CMD_4+1; button++ {
		costQueue[button] = NOT_ANY_BUTTON
		sortedQueue[button] = NOT_ANY_BUTTON
	}
	log.Println(ColG, "InitQueue: localQueue: ", localQueue)
}

//Rar funksjon
func SynchronizeQueueWithIO(pressedButtons map[int]bool) {
	for key, _ := range pressedButtons {
		localQueue[key] = pressedButtons[key]
	}
}

func RemoveFromLocalQueue(order int) {
	localQueue[order] = false
}

func GetNextOrder() (int, int) {
	for _, btn_states := range sortedQueue {
		if btn_states == NOT_ANY_BUTTON {
			return NOT_ANY_BUTTON, NOT_ANY_BUTTON
		} else {
			lowestCost, button := costFunc.LowestCostElevator(btn_states)
			if lowestCost && button == CMD_BTN {
				return btn_states, CMD_BTN
			} else if lowestCost && button == OUTSIDE_BTN {
				return btn_states, button
			}
		}
	}
	return NOT_ANY_BUTTON, NOT_ANY_BUTTON
}

func UpdateElevStateMap(name string, floor int, direction int, reserved int) {
	costFunc.ElevStateMap[name] = costFunc.ElevState{floor, direction, reserved}
	log.Println(ColB, "UpdateElevStateMap: (name, floor, direction): ", name, floor, direction, ColN)
}

func sortQueue() {
	log.Println("sortQueue: localQueue: ", localQueue)
	minButton := NOT_ANY_BUTTON
	for elem := UP_1; elem < CMD_4+1; elem++ {
		dummyLocalQueue[elem] = localQueue[elem]
	}
	for i := UP_1; i < CMD_4+1; i++ {
		min := 100 //Setter en høy cost
		for j := UP_1; j < CMD_4+1; j++ {
			if dummyLocalQueue[j] {
				if min > costQueue[j] {
					min = costQueue[j]
					minButton = j
				}
			}
		}
		if minButton != NOT_ANY_BUTTON {
			dummyLocalQueue[minButton] = false
		}
		sortedQueue[i] = minButton
		minButton = NOT_ANY_BUTTON
	}
	log.Println(ColC, "SortQueue; sortedQueue: ", sortedQueue, ColN)
}

func EmptyQueue() bool {
	if sortedQueue[0] == NOT_ANY_BUTTON {
		return true
	}
	return false
}
