package converter

import (
	"../config"
)

func ConvertDirAndFloorToMapIndex(floor int, direction int) (int, int) { //elevState *ElevState
	if floor == config.FLOOR_1 {
		return config.UP_1, config.CMD_1
	} else if floor == config.FLOOR_2 && direction == config.DIR_DOWN {
		return config.DOWN_2, config.CMD_2
	} else if floor == config.FLOOR_2 && direction == config.DIR_UP {
		return config.UP_2, config.CMD_2
	} else if floor == config.FLOOR_2 && direction == config.DIR_STOP {
		return config.UP_2, config.DOWN_2
	} else if floor == config.FLOOR_3 && direction == config.DIR_DOWN {
		return config.DOWN_3, config.CMD_3
	} else if floor == config.FLOOR_3 && direction == config.DIR_UP {
		return config.UP_3, config.CMD_3
	} else if floor == config.FLOOR_2 && direction == config.DIR_STOP {
		return config.UP_3, config.DOWN_3
	} else if floor == config.FLOOR_4 {
		return config.DOWN_4, config.CMD_4
	}
	return -1, -1
}

func ConvertButtonToFloor(button int) int {
	if button == config.CMD_4 || button == config.DOWN_4 {
		return 3
	} else if button == config.CMD_3 || button == config.DOWN_3 || button == config.UP_3 {
		return 2
	} else if button == config.CMD_2 || button == config.DOWN_2 || button == config.UP_2 {
		return 1
	} else if button == config.CMD_1 || button == config.UP_1 {
		return 0
	}
	return -1
}