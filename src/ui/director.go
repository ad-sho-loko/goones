package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"go_nes/src/nes"
	"image"
)

type Director struct {
	nes *nes.Nes
	window *glfw.Window
	gameView View
	ch chan *image.RGBA
}

func newDirector(nes *nes.Nes, window *glfw.Window) *Director {
	return &Director{
		nes:nes,
		window:window,
	}
}

func (d *Director) setKeyCallback(){
	callback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey){
		switch key {
		case glfw.KeyUp : d.nes.PushKey(nes.Up)
		case glfw.KeyDown : d.nes.PushKey(nes.Down)
		case glfw.KeyLeft : d.nes.PushKey(nes.Left)
		case glfw.KeyRight : d.nes.PushKey(nes.Right)
		}
	}
	d.window.SetKeyCallback(callback)
}

func (d *Director) start(){
	d.nes.Init()
	d.setKeyCallback()
	d.playGame()

	// main loop
	for !d.window.ShouldClose() {
		// Do OpenGL stuff.
		d.update()
		// d.window.SwapBuffers()
		glfw.PollEvents()
	}
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