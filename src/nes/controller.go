package nes

type Controller struct {
	reset byte
	counter int
	buttons [8]bool
}

type Key uint

func NewController() *Controller{
	return &Controller{
		counter:0,
	}
}

func (c *Controller) SetButton(b [8]bool){
	c.buttons = b
}

func (c *Controller) read() byte{
	b := byte(0)

	if c.counter < 8 && c.buttons[c.counter] {
		b = 1
	}

	c.counter++
	if c.reset&1 == 1 {
		c.counter = 0
	}

	return b
}

func (c *Controller) write(b byte){
	c.reset = b
	if c.reset & 1 == 1 {
		c.counter = 0
	}
}
