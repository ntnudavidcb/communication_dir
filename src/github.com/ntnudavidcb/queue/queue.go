package queue

import (
	//"../config"
	//"../driver"
	//"time"
	"../io"
	"log"
)

var localQueue = [10]bool{}
func CheckOrder(floor int) bool{ //elevState *ElevState
	/*index1, index2 := io.ConvertDirAndFloorToMapIndex()
	return io.PressedButtons[index2] || io.PressedButtons[index1]*/
	if inLocalQueue(floor){
		return true
	}
	return false
}

func UpdateQueueWithButton(buttonPressed int){
	localQueue[buttonPressed] = true
}

func UpdateQueue(){
	//run some cost function here and sort the queue with it
}

func convertButtonCMD(buttonPressed int) (int, int){
	if buttonPressed > 8{
		return 5, 5
	}else if buttonPressed > 7{
		return 3, 4
	}
	}else if buttonPressed > 6{
		return 1, 2
	}else if buttonPressed > 5{
		return 0, 0
	}else{
		return buttonPressed, buttonPressed
	}
}


func AddToQueue(){
	//Dette blir allerede gjort fra IO, noe som skal fikses
	for key, value := range io.PressedButtons{
		floor := io.ConvertMapIndexToFloor(key)
		if value && !inLocalQueue(floor){
			localQueue = append(localQueue, floor)
			log.Println(localQueue)
		}
	}
}

func removeFromQueue(){

}

func GetNextOrder() int{
	if len(localQueue) == 0{
		return -1 //Ingenting i koen
	}else{
		return localQueue[0]
	}
}

func inLocalQueue(buttonPressed int) bool{
	return localQueue[buttonPressed]
}

