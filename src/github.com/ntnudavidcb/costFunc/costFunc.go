package costFunc

import (
	"math"
	"../config"
	"log"
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

const MIN_COST int = 0

func CostFunc(currentDir int, currentFloor int, button int) int { 
	var costMap [6][6]int {}
	costMap[0] = {0,1,2,3,4,5}
	for key := 1; key < 6; key++ {
		for val := range costMap[key]{
			if costMap[key-1][val] == MIN_COST {
				costMap[key][val] = 5
			} else {
				costMap[key][val] = costMap[key-1][val] - 1	
			}
		}
	}

	//Button equivalents if direction is MOVING
	buttonEquivalent := button
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

	//Button equivalents if direction is STOP
	if button == DOWN_2 || button == CMD_2 {
		buttonEquivalent = UP_2
	} else if button == DOWN_3 || button == CMD_3 {
		buttonEquivalent = UP_3
	}

	if currentDir == config.DIR_STOP{
		if currentFloor == config.FLOOR_1{
			if currentFloor < buttonEquivalent{
				return buttonEquivalent - currentFloor
			} else {return currentFloor - buttonEquivalent}
		} else if currentFloor == config.FLOOR_2{
			if currentFloor < buttonEquivalent{
				return buttonEquivalent - currentFloor
			} else {return currentFloor - buttonEquivalent}
		} else if currentFloor == config.FLOOR_3{
			if currentFloor < buttonEquivalent{
				return buttonEquivalent - currentFloor
			} else {return currentFloor - buttonEquivalent}
		} else if currentFloor == config.FLOOR_4{
			if currentFloor < buttonEquivalent{
				return buttonEquivalent - currentFloor
			} else {return currentFloor - buttonEquivalent}
		}
	}
	//If everything goes wrong
	log.Println("Input parameters are not valid")
	return 1000
}