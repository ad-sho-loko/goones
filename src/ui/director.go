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

func (d *Director) start(){
	d.nes.Init()
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