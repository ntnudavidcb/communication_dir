package queue

import (
	"../costFunc"
	"../converter"
	"log"
)

var localQueue = [10]bool{}
var costQueue = [10]int{}
var dummyLocalQueue = [10]bool{}
var sortedQueue = [10]int{}

/*Organized as follows: 
UP_1, UP_2, UP_3, DOWN_4, DOWN_3, DOWN_2, 
CMD_1, CMD_2, CMD_3,CMD_4*/

func SetMyIP(IP string){
	costFunc.MyIP = IP
}

func CheckOrder() bool {
	buttonPressed1, buttonPressed2 := converter.ConvertDirAndFloorToMapIndex(costFunc.ElevStateMap[costFunc.MyIP].Floor, costFunc.ElevStateMap[costFunc.MyIP].Direction)
	log.Println("Buttonpressed1, buttonPressed2: ", buttonPressed1, buttonPressed2)
	log.Println("CheckOrder: ", inLocalQueue(buttonPressed1) || inLocalQueue(buttonPressed2))
	return inLocalQueue(buttonPressed1) || inLocalQueue(buttonPressed2)
}

func inLocalQueue(buttonPressed int) bool {
	if buttonPressed == -1 {
		return false
	}
	log.Println("localQueue: ", localQueue)
	log.Println("inLocalQueue: ", localQueue[buttonPressed])
	return localQueue[buttonPressed]
}

func UpdateQueueWithButton(buttonPressed int) {
	localQueue[buttonPressed] = true
}

func UpdateQueueButtonPushed(buttonPushed int) { //Brukes ikke
	UpdateQueueWithButton(buttonPushed)
	updateCostQueue()
	SortQueue()
}
func UpdateQueueFloorReached() {
	updateCostQueue()
	SortQueue()
}

func updateCostQueue() {
	//currentFloor, currentDir, _ := io.GetElevState()
	currentFloor := costFunc.ElevStateMap[costFunc.MyIP].Floor
	currentDir := costFunc.ElevStateMap[costFunc.MyIP].Direction
	for button := 0; button < 10; button++ {
		costQueue[button] = costFunc.CostFunc(currentDir, currentFloor, button)
	}
}

func RemoveFromQueue(pressedButtons map[int]bool) {
	for key, _ := range pressedButtons {
		localQueue[key] = pressedButtons[key]
	}
}

func RemoveFromLocalQueue(order int){
	localQueue[order] = false
}

func GetNextOrder() (int, int) {
	for _, btn_states := range sortedQueue{
		if btn_states == -1{
			return -1, -1
		} else {
			lowestCost, button := costFunc.LowestCostElevator(btn_states)
			if lowestCost && button == costFunc.CMD_BTN{
				return btn_states, costFunc.CMD_BTN
			} else if lowestCost && button == costFunc.OUTSIDE_BTN {
				return btn_states, button
			}
		}
	}
	return -1, -1
}

func UpdateElevStateMap(name string, direction int, floor int){
	costFunc.ElevStateMap[name] = costFunc.ElevState{floor, direction}
}

func SortQueue() {
	minButton := -1
	for elem := 0; elem < 10; elem++ {
		dummyLocalQueue[elem] = localQueue[elem]
	}
	for i := 0; i < 10; i++ {
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
		sortedQueue[i] = minButton
	}
	//log.Println(config.ColC, "SortQueue; sortedQueue: ", sortedQueue, config.ColN)
	//log.Println(config.ColG, "SortQueue; localQueue: ", localQueue, config.ColN)
}