package costFunc

import (
	"math"
	"../config"
)

const (
	UP_1 int iota
	UP_2
	UP_3
	DOWN_4
	DOWN_3
	DOWN_2
	CMD_1
	CMD_2
	CMD_3
	CMD_4
)

const MAX_COST int = 5

func costFunc(currentDir int, currentFloor int, button int) int { 
	var costMap [6][6]int {}
	costMap[0] = {1,2,3,4,5}
	for key := 1; key < 6; key++ {
		for val := range costMap[key]{
			if costMap[key-1][val] == MAX_COST {
				costMap[key][val] = 0
			} else {
				costMap[key][val] = costMap[key-1][val] - 1	
			}
		}
	}

	var buttonEquivalent := button
	if button == CMD_1{
		buttonEquivalent = UP_1
	} else if button == CMD_4{
		buttonEquivalent = DOWN_4
	}

	if currentDir == config.DIR_UP{
		if button == CMD_2 {
			return costMap[currentFloor][math.Min(UP_2, DOWN_2)]
		} else if button == CMD_3 {
			return costMap[currentFloor][math.Min(UP_3, DOWN_3)]
		} 
		return costMap[currentFloor][buttonEquivalent]
	} else if currentDir == config.DIR_DOWN {
		if currentFloor == config.FLOOR_1 {
			if button == CMD_2 {
			return costMap[currentFloor][math.Min(UP_2, DOWN_2)]
		} else if button == CMD_3 {
			return costMap[currentFloor][math.Min(UP_3, DOWN_3)]
		} 
			return costMap[FLOOR_1][buttonEquivalent]
		} else if currentFloor == config.FLOOR_2{
			if button == CMD_2 {
			return costMap[5][math.Min(UP_2, DOWN_2)]
		} else if button == CMD_3 {
			return costMap[5][math.Min(UP_3, DOWN_3)]
		} 
			return costMap[5][buttonEquivalent]
		} else if currentFloor == config.FLOOR_3{
			if button == CMD_2 {
			return costMap[4][math.Min(UP_2, DOWN_2)]
		} else if button == CMD_3 {
			return costMap[4][math.Min(UP_3, DOWN_3)]
		} 
			return costMap[4][buttonEquivalent]
		}
	} 
}