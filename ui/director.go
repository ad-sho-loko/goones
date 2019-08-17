package ui

import (
	"github.com/ad-sho-loko/goones/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Director struct {
	nes *nes.Nes
	window *glfw.Window
	gameView View
}

func newDirector(nes *nes.Nes, window *glfw.Window) *Director {
	return &Director{
		nes:nes,
		window:window,
	}
}

var keyStates [8]bool

func (d *Director) setKeyCallback(){
	callback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey){
		if !(action == glfw.Press || action == glfw.Release){
			return
		}
		var isPush = action == glfw.Press

		switch key {
		case glfw.KeyA:
			keyStates[0] = isPush
		case glfw.KeyB:
			keyStates[1] = isPush
		case glfw.KeyRightShift :
			keyStates[2] = isPush
		case glfw.KeyEnter :
			keyStates[3] = isPush
		case glfw.KeyUp :
			keyStates[4] = isPush
		case glfw.KeyDown :
			keyStates[5] = isPush
		case glfw.KeyLeft :
			keyStates[6] = isPush
		case glfw.KeyRight :
			keyStates[7] = isPush
		}
		d.nes.PushButton(keyStates)
	}
	d.window.SetKeyCallback(callback)
}

func (d *Director) start(){
	d.nes.Init()
	d.setKeyCallback()
	d.playGame()

	// main loop
	for !d.window.ShouldClose() {
		d.update()
		d.window.SwapBuffers()
		glfw.PollEvents()
	}

	d.setView(nil)
}

func (d *Director) playGame(){
	gameView := newGameView(d)
	d.setView(gameView)
}

func (d *Director) update(){
	gl.Clear(gl.COLOR_BUFFER_BIT)
	d.gameView.Update()
}

func (d *Director) setView(view View){
	if view != nil{
		d.gameView = view
	}
	d.gameView.Enter()
}