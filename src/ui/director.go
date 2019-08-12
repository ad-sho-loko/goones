package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"go_nes/src/nes"
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

func (d *Director) setKeyCallback(){
	callback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey){
		if !(action == glfw.Press || action == glfw.Release){
			return
		}

		// 両方押したとき周りの挙動が少しおかしいので要修正
		var b [8]bool
		var isPush = action == glfw.Press
		switch key {
		case glfw.KeyA:
			b[0] = isPush
			fallthrough
		case glfw.KeyB:
			b[1] = isPush
			fallthrough
		case glfw.KeyRightShift :
			b[2] = isPush
			fallthrough
		case glfw.KeyEnter :
			b[3] = isPush
			fallthrough
		case glfw.KeyUp :
			b[4] = isPush
			fallthrough
		case glfw.KeyDown :
			b[5] = isPush
			fallthrough
		case glfw.KeyLeft :
			b[6] = isPush
			fallthrough
		case glfw.KeyRight :
			b[7] = isPush
		}
		d.nes.PushButton(b)
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