package main

import (
	//"../config"
	"../driver"
	//"../queue"
	"../io"
	"log"
	//"time"
)


/*func testrun1() {
	log.Println(config.ColB, "Test Run 1 Initialized", config.ColN)

	//Etasje 1 ankommet etter init
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_1, 1)
	driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_1, 1)
	time.Sleep(3000 * time.Millisecond)
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_1, 0)
	driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_1, 0)

	//Etasje 2 neste stopp
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_2, 1)
	driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_2, 1)
	driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_2, 1)
	for driver.Elev_get_floor_sensor_signal() != config.FLOOR_2 {
		driver.Elev_set_motor_elevState.direction(config.DIR_UP)
	}
	driver.Elev_set_motor_elevState.direction(config.DIR_STOP)
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_2, 1)
	driver.Elev_set_floor_indicator(config.FLOOR_2)
	time.Sleep(3000 * time.Millisecond)
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_2, 0)
	driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_2, 0)
	driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_2, 0)

	//Etasje 3 neste stopp
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_3, 1)
	driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_3, 1)
	driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_3, 1)
	for driver.Elev_get_floor_sensor_signal() != config.FLOOR_3 {
		driver.Elev_set_motor_elevState.direction(config.DIR_UP)
	}
	driver.Elev_set_motor_elevState.direction(config.DIR_STOP)
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_3, 1)
	driver.Elev_set_floor_indicator(config.FLOOR_3)
	time.Sleep(3000 * time.Millisecond)
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_3, 0)
	driver.Elev_set_button_lamp(config.BTN_UP, config.FLOOR_3, 0)
	driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_3, 0)

	//Etasje 4 neste stopp
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_4, 1)
	driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_4, 1)
	for driver.Elev_get_floor_sensor_signal() != config.FLOOR_4 {
		driver.Elev_set_motor_elevState.direction(config.DIR_UP)
	}
	driver.Elev_set_motor_elevState.direction(config.DIR_STOP)
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_4, 1)
	driver.Elev_set_floor_indicator(config.FLOOR_4)
	time.Sleep(3000 * time.Millisecond)
	driver.Elev_set_button_lamp(config.BTN_COMMAND, config.FLOOR_4, 0)
	driver.Elev_set_button_lamp(config.BTN_DOWN, config.FLOOR_4, 0)

	log.Println(config.ColB, "Test Run Finished", config.ColN)
}*/


func main() {	
	asd := make(chan int)
	driver.Elev_init()
	
	log.Println("Hei")
	go io.Testrun2()

	asd <- 1
}