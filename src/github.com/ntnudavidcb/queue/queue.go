package queue

import (
	"../config"
	"../driver"
	//"time"
	"log"
)

var local_queue[]int

var main_queue[]int

func Queue_testing(){

}

func ReadAllButtons(){
	log.Println("KJsh")
	for{
		if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_1) && !ifInQueue(config.FLOOR_1, local_queue){
			local_queue = append(local_queue, config.FLOOR_1)
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_2)&& !ifInQueue(config.FLOOR_2, local_queue){
			local_queue = append(local_queue, config.FLOOR_2)
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_3)&& !ifInQueue(config.FLOOR_3, local_queue){
			local_queue = append(local_queue, config.FLOOR_3)
		}else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_4)&& !ifInQueue(config.FLOOR_4, local_queue){
			local_queue = append(local_queue, config.FLOOR_4)

		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_1) && !ifInQueue(config.FLOOR_1, main_queue){
			main_queue = append(main_queue, config.FLOOR_1)
		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_2)&& !ifInQueue(config.FLOOR_2, main_queue){
			main_queue = append(main_queue, config.FLOOR_2)
		}else if driver.Elev_get_button_signal(config.BTN_UP, config.FLOOR_3)&& !ifInQueue(config.FLOOR_3, main_queue){
			main_queue = append(main_queue, config.FLOOR_3)
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_4)&& !ifInQueue(config.FLOOR_4, main_queue){
			main_queue = append(main_queue, config.FLOOR_4)
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_3)&& !ifInQueue(config.FLOOR_3, main_queue){
			main_queue = append(main_queue, config.FLOOR_3)
		}else if driver.Elev_get_button_signal(config.BTN_DOWN, config.FLOOR_2)&& !ifInQueue(config.FLOOR_2, main_queue){
			main_queue = append(main_queue, config.FLOOR_2)
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

func Queue_add_to_main(){
	//Legger til bestillinger fra UP/DOWN-knappene
}

func Queue_check_main(){
	//Heis sjekker hovedkøen (10 ganger i sekunder f.eks)
}

func Queue_remove_order(){

}