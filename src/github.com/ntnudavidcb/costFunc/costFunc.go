package costFunc

import (
	. "../config"
	"log"
	"strconv"
	"strings"
)

type ElevState struct {
	Floor     int
	Direction int
	Reserved  int
}

var ElevStateMap = make(map[string]ElevState)
var MyIP string

const MIN_COST int = 0

func DelElevStateMap() {
	for key, _ := range ElevStateMap {
		delete(ElevStateMap, key)
	}
}

func minIntegerFunc(integer1 int, integer2 int) int {
	if integer1 < integer2 {
		return integer1
	} else {
		return integer2
	}
}

func LowestCostElevator(button int) (bool, int) {
	if button > CMD_4 && button < NOT_ANY_BUTTON {
		Restart.Run()
		log.Fatal("LowestCostElevator recieved wrong input, button: ", button)
	}

	log.Println("button, ElevStateMap", button, ElevStateMap)
	if button == CMD_1 || button == CMD_2 || button == CMD_3 || button == CMD_4 || len(ElevStateMap) == 0 {
		return true, CMD_BTN
	} else {
		smallestIPStruct := ElevStateMap[MyIP]
		smallestIPList := []string{}
		for IPs, elevState := range ElevStateMap {
			if CostFunc(smallestIPStruct.Direction, smallestIPStruct.Floor, button) > CostFunc(elevState.Direction, elevState.Floor, button) {
				smallestIPStruct = elevState
				smallestIPList = []string{IPs}
			} else if len(smallestIPList) == 0 {
				smallestIPList = append(smallestIPList, IPs)
			} else if CostFunc(smallestIPStruct.Direction, smallestIPStruct.Floor, button) == CostFunc(elevState.Direction, elevState.Floor, button) {
				smallestIPList = append(smallestIPList, IPs)
			}
		}
		if len(smallestIPList) == 1 && smallestIPList[0] == MyIP {
			return true, OUTSIDE_BTN
		} else if len(smallestIPList) == 1 {
			return false, -1
		}
		smallestIP := smallestIPList[0]
		for i, _ := range smallestIPList {
			str1, _ := strconv.Atoi(strings.Split(smallestIP, ".")[3])
			str2, _ := strconv.Atoi(strings.Split(smallestIPList[i], ".")[3])
			if str1 > str2 {
				smallestIP = smallestIPList[i]
			}
		}
		if smallestIP == MyIP {
			return true, OUTSIDE_BTN
		} else {
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
	log.Println("ElevStateMap: ", ElevStateMap)
	if button > CMD_4 && button < NOT_ANY_BUTTON && currentFloor < FLOOR_1 && currentFloor > FLOOR_4 {
		Restart.Run()
		log.Fatal("CostFunc recieved wrong input, currentDir, currentFloor, button: ", currentDir, currentFloor, button)
	}

	//Button equivalents if Direction is MOVING
	buttonEquivalent := button
	if button == CMD_1 {
		buttonEquivalent = UP_1
	} else if button == CMD_4 {
		buttonEquivalent = DOWN_4
	}

	if currentDir == DIR_UP {
		if button == CMD_2 {
			return costMap[currentFloor][minIntegerFunc(UP_2, DOWN_2)]
		} else if button == CMD_3 {
			return costMap[currentFloor][minIntegerFunc(UP_3, DOWN_3)]
		}
		return costMap[currentFloor][buttonEquivalent]
	} else if currentDir == DIR_DOWN {
		if currentFloor == FLOOR_1 {
			if button == CMD_2 {
				return costMap[currentFloor][minIntegerFunc(UP_2, DOWN_2)]
			} else if button == CMD_3 {
				return costMap[currentFloor][minIntegerFunc(UP_3, DOWN_3)]
			}
			return costMap[FLOOR_1][buttonEquivalent]
		} else if currentFloor == FLOOR_2 {
			if button == CMD_2 {
				return costMap[DOWN_2][minIntegerFunc(UP_2, DOWN_2)]
			} else if button == CMD_3 {
				return costMap[DOWN_2][minIntegerFunc(UP_3, DOWN_3)]
			}
			return costMap[DOWN_2][buttonEquivalent]
		} else if currentFloor == FLOOR_3 {
			if button == CMD_2 {
				return costMap[DOWN_3][minIntegerFunc(UP_2, DOWN_2)]
			} else if button == CMD_3 {
				return costMap[DOWN_3][minIntegerFunc(UP_3, DOWN_3)]
			}
			return costMap[DOWN_3][buttonEquivalent]
		} else if currentFloor == FLOOR_4 {
			if button == CMD_2 {
				return costMap[FLOOR_4][minIntegerFunc(UP_2, DOWN_2)]
			} else if button == CMD_3 {
				return costMap[FLOOR_4][minIntegerFunc(UP_3, DOWN_3)]
			}
			return costMap[FLOOR_4][buttonEquivalent]
		}
	}

	//Button equivalents if Direction is STOP
	if button == DOWN_2 || button == CMD_2 {
		buttonEquivalent = UP_2
	} else if button == DOWN_3 || button == CMD_3 {
		buttonEquivalent = UP_3
	}

	if currentDir == DIR_STOP {
		if currentFloor == FLOOR_1 {
			if currentFloor < buttonEquivalent {
				return buttonEquivalent - currentFloor
			} else {
				return currentFloor - buttonEquivalent
			}
		} else if currentFloor == FLOOR_2 {
			if currentFloor < buttonEquivalent {
				return buttonEquivalent - currentFloor
			} else {
				return currentFloor - buttonEquivalent
			}
		} else if currentFloor == FLOOR_3 {
			if currentFloor < buttonEquivalent {
				return buttonEquivalent - currentFloor
			} else {
				return currentFloor - buttonEquivalent
			}
		} else if currentFloor == FLOOR_4 {
			if currentFloor < buttonEquivalent {
				return buttonEquivalent - currentFloor
			} else {
				return currentFloor - buttonEquivalent
			}
		}
	}
	//If everything goes wrong
	Restart.Run()
	log.Fatal("CostFunc not working correctly")

	return 1000
}
