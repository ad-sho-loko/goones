package nes

type Controller struct {
	resetFlag bool
	counter int
	A byte
	B byte
	Select byte
	Start byte
	Up byte
	Down byte
	Left byte
	Right byte
}

type Key uint

const(
	A Key = iota
	B
	Select
	Start
	Up
	Down
	Left
	Right
)

func NewController() *Controller{
	return &Controller{
		counter:0,
	}
}

func (c *Controller) read() byte{
	var b byte = 0

	switch c.counter {
	case 0: b = c.A
	case 1: b = c.B
	case 2: b = c.Select
	case 3: b = c.Start
	case 4: b = c.Up
	case 5: b = c.Down
	case 6: b = c.Left
	case 7: b = c.Right
	}

	c.counter++
	c.counter%=8
	return b
}

func (c *Controller) reset(b byte){
	if !c.resetFlag && b == 1{
		c.resetFlag = true
	}else if c.resetFlag && b == 0{
		c.resetFlag = false
		c.counter = 0
		c.A = 0
		c.B  = 0
		c.Select = 0
		c.Start = 0
		c.Up = 0
		c.Down = 0
		c.Left = 0
		c.Right = 0
	}
}

func (c *Controller) set(key Key, b byte){
	switch key {
	case A: c.A = b
	case B: c.B = b
	case Select: c.Select = b
	case Start: c.Start = b
	case Up: c.Up = b
	case Down: c.Down = b
	case Left: c.Left = b
	case Right: c.Right = b
	}
}
