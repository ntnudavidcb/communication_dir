package costFunc

import (
	"../config"
	"log"
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

const (
	CMD_BTN = 1
	OUTSIDE_BTN = 2
)

type ElevState struct {
	Floor int
	Direction int
}

var ElevStateMap = make(map[string]ElevState) //string = IP
var MyIP string

const MIN_COST int = 0

func minIntegerFunc(integer1 int, integer2 int) int {
	if integer1 < integer2 {
		return integer1
	} else {
		return integer2
	}
}

func LowestCostElevator(button int) (bool, int) {
	if button == CMD_1 || button == CMD_2 || button == CMD_3 || button == CMD_4 || len(ElevStateMap) == 0{
		return true, CMD_BTN
	} else {
		smallestIPStruct := ElevStateMap[MyIP]
		smallestIPList := []string{}
		for IPs, elevState := range ElevStateMap{
			if CostFunc(smallestIPStruct.Direction, smallestIPStruct.Floor, button) > CostFunc(elevState.Direction, elevState.Floor, button){
				smallestIPStruct = elevState
				smallestIPList = []string{IPs}
			} else if len(smallestIPList) == 0 {
				smallestIPList = append(smallestIPList, IPs)
			} else if CostFunc(smallestIPStruct.Direction, smallestIPStruct.Floor, button) == CostFunc(elevState.Direction, elevState.Floor, button){
				smallestIPList = append(smallestIPList, IPs)
			} 
		}
		if len(smallestIPList) == 1 && smallestIPList[0] == MyIP{
			return true, OUTSIDE_BTN
		} else if len(smallestIPList) == 1{
			return false, -1
		}
		smallestIP := smallestIPList[0]
		for i, _ := range smallestIPList {
			if smallestIP > smallestIPList[i]{
				smallestIP = smallestIPList[i]
			}
		}
		if smallestIP == MyIP{
			return true, OUTSIDE_BTN
		} else{
			return false, -1
		}
	}
	return false, -1
}

func CostFunc(currentDir int, currentFloor int, button int) int {
	costMap := [][]int{
		[]int{0, 1, 2, 3, 4, 5},
		[]int{5, 0, 1, 2, 3, 4},
		[]int{4, 5, 0, 1, 2, 3},
		[]int{3, 4, 5, 0, 1, 2},
		[]int{2, 3, 4, 5, 0, 1},
		[]int{1, 2, 3, 4, 5, 0},
	}
	log.Println("currentDir: ", currentDir)
	log.Println("currentFloor: ", currentFloor)
	log.Println("button: ", button)

	//Button equivalents if Direction is MOVING
	buttonEquivalent := button
	if button == CMD_1 {
		buttonEquivalent = UP_1
	} else if button == CMD_4 {
		buttonEquivalent = DOWN_4
	}

	if currentDir == config.DIR_UP {
		if button == CMD_2 {
			return costMap[currentFloor][minIntegerFunc(UP_2, DOWN_2)]
		} else if button == CMD_3 {
			return costMap[currentFloor][minIntegerFunc(UP_3, DOWN_3)]
		}
		return costMap[currentFloor][buttonEquivalent]
	} else if currentDir == config.DIR_DOWN {
		if currentFloor == config.FLOOR_1 {
			if button == CMD_2 {
				return costMap[currentFloor][minIntegerFunc(UP_2, DOWN_2)]
			} else if button == CMD_3 {
				return costMap[currentFloor][minIntegerFunc(UP_3, DOWN_3)]
			}
			return costMap[config.FLOOR_1][buttonEquivalent]
		} else if currentFloor == config.FLOOR_2 {
			if button == CMD_2 {
				return costMap[5][minIntegerFunc(UP_2, DOWN_2)]
			} else if button == CMD_3 {
				return costMap[5][minIntegerFunc(UP_3, DOWN_3)]
			}
			return costMap[5][buttonEquivalent]
		} else if currentFloor == config.FLOOR_3 {
			if button == CMD_2 {
				return costMap[4][minIntegerFunc(UP_2, DOWN_2)]
			} else if button == CMD_3 {
				return costMap[4][minIntegerFunc(UP_3, DOWN_3)]
			}
			return costMap[4][buttonEquivalent]
		} else if currentFloor == config.FLOOR_4 {
			if button == CMD_2 {
				return costMap[config.FLOOR_4][minIntegerFunc(UP_2, DOWN_2)]
			} else if button == CMD_3 {
				return costMap[config.FLOOR_4][minIntegerFunc(UP_3, DOWN_3)]
			}
			return costMap[config.FLOOR_4][buttonEquivalent]
		}
	}

	//Button equivalents if Direction is STOP
	if button == DOWN_2 || button == CMD_2 {
		buttonEquivalent = UP_2
	} else if button == DOWN_3 || button == CMD_3 {
		buttonEquivalent = UP_3
	}

	if currentDir == config.DIR_STOP {
		if currentFloor == config.FLOOR_1 {
			if currentFloor < buttonEquivalent {
				return buttonEquivalent - currentFloor
			} else {
				return currentFloor - buttonEquivalent
			}
		} else if currentFloor == config.FLOOR_2 {
			if currentFloor < buttonEquivalent {
				return buttonEquivalent - currentFloor
			} else {
				return currentFloor - buttonEquivalent
			}
		} else if currentFloor == config.FLOOR_3 {
			if currentFloor < buttonEquivalent {
				return buttonEquivalent - currentFloor
			} else {
				return currentFloor - buttonEquivalent
			}
		} else if currentFloor == config.FLOOR_4 {
			if currentFloor < buttonEquivalent {
				return buttonEquivalent - currentFloor
			} else {
				return currentFloor - buttonEquivalent
			}
		}
	}
	//If everything goes wrong
	log.Println("CostFunc not working correctly")
	return 1000
}
