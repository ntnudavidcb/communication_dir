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

const MIN_COST int = 0

func minIntegerFunc(integer1 int, integer2 int) int {
	if integer1 < integer2 {
		return integer1
	} else {
		return integer2
	}
}

func LowestCostElevator(myDir int, myFloor int, myActive bool, button int) bool {
	if button == CMD_1 || button == CMD_2 || button == CMD_3 || button == CMD_4 {
		return true
	}
	//Hvis en heis er active så prioriteres den ikke
	return false
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

	//Button equivalents if direction is MOVING
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
		}
	}

	//Button equivalents if direction is STOP
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
