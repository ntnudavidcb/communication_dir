package queue

import (
	"../config"
	"../driver"
	"time"
	"log"
)

var local_queue[]int



func Queue_testing(){

}

func Queue_add_to_local(){
	//Legger til bestillinger fra COMMAND-knappene
	n := len(local_queue)
	for{
	if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_1){
		local_queue = append(local_queue, config.FLOOR_1)
	}
	else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_2){
		local_queue = append(local_queue, config.FLOOR_2)
	}
	else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_3){
		local_queue = append(local_queue, config.FLOOR_3)
	}
	else if driver.Elev_get_button_signal(config.BTN_COMMAND, config.FLOOR_4){
		local_queue = append(local_queue, config.FLOOR_4)
	}
}
}

func Queue_add_to_main(){
	//Legger til bestillinger fra UP/DOWN-knappene
}

func Queue_check_main(){
	//Heis sjekker hovedkøen (10 ganger i sekunder f.eks)
}

func Queue_remove_order(){

}