package driver

//inspirasjon fra mortenfyhn, fiks matrix-imputs i funksjonene slik at index'ene stemmer.

import (
	. "../config"
	"log"
	"time"
)

//Mulig legge alle const-deklarasjoner inn i en config.go

type ELEV_BUTTON_TYPE int
type ELEV_MOTOR_DIR int

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

//Initialized and descends the lift to a defined state (until the lift reaches a floor)
func Elev_init() (int, error) {
	if Io_init() == false {
		Restart.Run()
		log.Fatal("Hardware driver: ioInit() failed!")
	}
	for f := 0; f < N_FLOORS; f++ {
		if f != 0 {
			Elev_set_button_lamp(BTN_DOWN, f, OFF)
		}
		if f != N_FLOORS-1 {
			Elev_set_button_lamp(BTN_UP, f, OFF)
		}
		Elev_set_button_lamp(BTN_COMMAND, f, OFF)
	}

	Elev_set_stop_lamp(OFF)
	Elev_set_door_open_lamp(OFF)

	Elev_set_motor_direction(DIR_DOWN)
	floor := Elev_get_floor_sensor_signal()
	for floor == NOT_ANY_FLOOR {
		floor = Elev_get_floor_sensor_signal()
	}

	Elev_set_motor_direction(DIR_STOP)
	Elev_set_floor_indicator(floor)
	time.Sleep(7 * time.Second)
	log.Println(ColG, "Hardware Initialized", ColN)
	return floor, nil
}

func Elev_set_motor_direction(dirn ELEV_MOTOR_DIR) int {
	if dirn == DIR_STOP {
		Io_write_analog(MOTOR, DIR_STOP)
		return DIR_STOP
	} else if dirn > DIR_STOP {
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR, 2800)
		return DIR_UP
	} else if dirn < DIR_STOP {
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR, 2800)
		return DIR_DOWN
	}
	return DIR_STOP
}

func Elev_set_button_lamp(button int, floor int, value int) {
	if floor < FLOOR_1 || floor >= N_FLOORS {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return
	}
	if button == BTN_UP && floor == N_FLOORS-1 {
		log.Println("Button up from top floor does not exist!")
		return
	}
	if button == BTN_DOWN && floor == FLOOR_1 {
		log.Println("Button down from ground floor does not exist!")
		return
	}
	if button != BTN_UP &&
		button != BTN_DOWN &&
		button != BTN_COMMAND {
		log.Printf("Invalid button %d\n", button)
		return
	}
	if value != OFF {
		Io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		Io_clear_bit(lamp_channel_matrix[floor][button])
	}
}

func Elev_set_floor_indicator(floor int) {
	if floor < FLOOR_1 || floor >= N_FLOORS {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return
	}
	// Binary encoding. One light must always be on.
	if floor&0x02 > 0 {
		Io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND1)
	}

	if floor&0x01 > 0 {
		Io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND2)
	}
}

func Elev_set_door_open_lamp(value int) {
	if value != OFF {
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func Elev_set_stop_lamp(value int) {
	if value != OFF {
		Io_set_bit(LIGHT_STOP)
	} else {
		Io_clear_bit(LIGHT_STOP)
	}
}

func Elev_get_button_signal(button int, floor int) bool {
	if floor < FLOOR_1 || floor >= N_FLOORS {
		log.Printf("Error: Floor %d out of range!\n", floor)
		return false
	}
	if button < 0 || button >= N_BUTTONS {
		log.Printf("Error: Button %d out of range!\n", button)
		return false
	}
	if button == BTN_UP && floor == N_FLOORS-1 {
		log.Println("Button up from top floor does not exist!")
		return false
	}
	if button == BTN_DOWN && floor == FLOOR_1 {
		log.Println("Button down from ground floor does not exist!")
		return false
	}

	//fmt.Printf("button %i, floor %i\n", button, floor)
	if Io_read_bit(button_channel_matrix[floor][button]) != 0 {
		return true
	} else {
		return false
	}
}

func Elev_get_floor_sensor_signal() int {
	if Io_read_bit(SENSOR_FLOOR1) != 0 {
		return FLOOR_1
	} else if Io_read_bit(SENSOR_FLOOR2) != 0 {
		return FLOOR_2
	} else if Io_read_bit(SENSOR_FLOOR3) != 0 {
		return FLOOR_3
	} else if Io_read_bit(SENSOR_FLOOR4) != 0 {
		return FLOOR_4
	} else {
		return NOT_ANY_FLOOR
	}
}

func Elev_get_stop_signal() int {
	return Io_read_bit(STOP)
}

func Elev_get_obstruction_signal() int {
	return Io_read_bit(OBSTRUCTION)
}
