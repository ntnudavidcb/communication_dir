package queue

import (
	"../config"
	"../driver"
	//"time"
	"log"
)

const(
	CMD_1 int = iota
	CMD_2
	CMD_3
	CMD_4
	UP_1
	DOWN_2
	UP_2
	DOWN_3
	UP_3
	DOWN_4
)

local_queue := make(map[int]bool)

var main_queue[]int

func Queue_testing(){

}

func ReadAllButtons(){
	log.Println("KJsh")
	for{
		if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_1) && !ifInQueue(config.FLOOR_1, local_queue){
			local_queue[CMD_1] = true
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_2)&& !ifInQueue(config.FLOOR_2, local_queue){
			local_queue[CMD_2] = true
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_3)&& !ifInQueue(config.FLOOR_3, local_queue){
			local_queue[CMD_3] = true
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_4)&& !ifInQueue(config.FLOOR_4, local_queue){
			local_queue[CMD_4] = true

		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_1) && !ifInQueue(UP_1, main_queue){
			main_queue = append(main_queue, UP_1)
		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_2)&& !ifInQueue(UP_2, main_queue){
			main_queue = append(main_queue, UP_2)
		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_3)&& !ifInQueue(UP_3, main_queue){
			main_queue = append(main_queue, UP_3)
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_4)&& !ifInQueue(DOWN_4, main_queue){
			main_queue = append(main_queue, DOWN_4)
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_3)&& !ifInQueue(DOWN_3, main_queue){
			main_queue = append(main_queue, DOWN_3)
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_2)&& !ifInQueue(DOWN_2, main_queue){
			main_queue = append(main_queue, DOWN_2)
		}
		log.Println("Local queue: ", local_queue)
		log.Println("Main queue: ", main_queue)
	}
}

func ifInQueue(number int ,queue []int) bool {
	for _, j := range queue{
		if j == number{
			return true
		}
	}
	return false
}

func CheckOrder(prevFloor int, floor int) bool{
	return local_queue[floor]
}


func removeFromQueue(index int, queue []int){
	queue = append(queue[:index], queue[index+1:]...)
}