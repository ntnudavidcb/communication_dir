package config

import (
	"os/exec"
)

const N_FLOORS = 4
const N_BUTTONS = 3

const (
	OFF = 0
	ON  = 1
)

const (
	NOT_ANY_FLOOR = -1
	FLOOR_1       = 0
	FLOOR_2       = 1
	FLOOR_3       = 2
	FLOOR_4       = 3
)

const (
	BTN_UP      = 0
	BTN_DOWN    = 1
	BTN_COMMAND = 2
)

const (
	DIR_UP   = 1
	DIR_DOWN = -1
	DIR_STOP = 0
)

const (
	CMD_BTN     = 1
	OUTSIDE_BTN = 2
)

const (
	NOT_ANY_BUTTON = -1
	UP_1           = 0
	UP_2           = 1
	UP_3           = 2
	DOWN_4         = 3
	DOWN_3         = 4
	DOWN_2         = 5
	CMD_1          = 6
	CMD_2          = 7
	CMD_3          = 8
	CMD_4          = 9
)

var Restart = exec.Command("gnome-terminal", "-x", "go", "run", "main.go")

// Colours for printing to console
const Col0 = "\x1b[30;1m" // Dark grey
const ColR = "\x1b[31;1m" // Red
const ColG = "\x1b[32;1m" // Green
const ColY = "\x1b[33;1m" // Yellow
const ColB = "\x1b[34;1m" // Blue
const ColM = "\x1b[35;1m" // Magenta
const ColC = "\x1b[36;1m" // Cyan
const ColW = "\x1b[37;1m" // White
const ColN = "\x1b[0m"    // Grey (neutral)
