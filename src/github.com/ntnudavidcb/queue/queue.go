package queue

import (
	"../costFunc"
	"../converter"
	"log"
	"../config"
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

func CheckOrder(floor int, direction int) bool {
	buttonUp, buttonDown, buttonCMD := converter.ConvertDirAndFloorToMapIndex(floor, direction) //costFunc.ElevStateMap[costFunc.MyIP].Floor, costFunc.ElevStateMap[costFunc.MyIP].Direction
	log.Println("CheckOrder: buttonUp, buttonDown, buttonCMD: ", buttonUp, buttonDown, buttonCMD)
	log.Println("CheckOrder: localQueue: ", localQueue)
	return inLocalQueue(buttonUp) || inLocalQueue(buttonDown) || inLocalQueue(buttonCMD)
}

func CheckUpOrDownButton() int{
	buttonUp, buttonDown, _ := converter.ConvertDirAndFloorToMapIndex(costFunc.ElevStateMap[costFunc.MyIP].Floor, costFunc.ElevStateMap[costFunc.MyIP].Direction)
	log.Println("CheckUpOrDownButton: Up, down: ", buttonUp, buttonDown)
	if inLocalQueue(buttonUp){
		return config.BTN_UP
	} else if inLocalQueue(buttonDown){
		return config.BTN_DOWN
	} else{
		return config.BTN_COMMAND
	}
}

func inLocalQueue(buttonPressed int) bool {
	if buttonPressed == -1 {
		return false
	}
	return localQueue[buttonPressed]
}

func AddButtonToQueue(buttonPressed int) {
	localQueue[buttonPressed] = true
}

func UpdateQueue() {
	updateCostQueue()
	sortQueue()
}

func updateCostQueue() {
	currentFloor := costFunc.ElevStateMap[costFunc.MyIP].Floor
	currentDir := costFunc.ElevStateMap[costFunc.MyIP].Direction
	for button := 0; button < 10; button++ {
		costQueue[button] = costFunc.CostFunc(currentDir, currentFloor, button)
	}
}

func RemoveButtonFromQueue(button int){
	if button == config.NOT_ANY_BUTTON{
		return
	}
	localQueue[button] = false
}

func InitQueue(){
	for i := 0; i < 10; i++ {
		costQueue[i] = -1
	}
	for i := 0; i < 10; i++ {
		sortedQueue[i] = -1
	}
}

//Rar funksjon
func SynchronizeQueueWithIO(pressedButtons map[int]bool) {
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

func UpdateElevStateMap(name string, floor int, direction int){
	costFunc.ElevStateMap[name] = costFunc.ElevState{floor, direction}
	log.Println(config.ColB, "UpdateElevStateMap: (name, floor, direction): ", name, floor, direction, config.ColN)
}

func sortQueue() {
	log.Println("sortQueue: localQueue: ", localQueue)
	minButton := config.NOT_ANY_BUTTON
	for elem := 0; elem < config.CMD_4+1; elem++ {
		dummyLocalQueue[elem] = localQueue[elem]
	}
	for i := config.UP_1; i < config.CMD_4+1; i++ {
		min := 100 //Setter en høy cost
		for j := config.UP_1; j < config.CMD_4+1; j++ {
			if dummyLocalQueue[j] {
				if min > costQueue[j] {
					min = costQueue[j]
					minButton = j
				}
			}
		}
		if minButton != config.NOT_ANY_BUTTON {
			dummyLocalQueue[minButton] = false
		}
		sortedQueue[i] = minButton
		minButton = config.NOT_ANY_BUTTON
	}
	log.Println(config.ColC, "SortQueue; sortedQueue: ", sortedQueue, config.ColN)
}
 
func EmptyQueue() bool{
	if sortedQueue[0] == -1{
		return true
	}
	return false
}