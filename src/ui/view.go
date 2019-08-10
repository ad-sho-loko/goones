package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"image"
)

type View interface {
	Enter()
	Update()
}

type GameView struct {
	View
	director *Director
	texture uint32
}

func newGameView(d *Director) View{
	return &GameView{
		director:d,
		texture:createTexture(),
	}
}

func (g *GameView) Enter(){
	gl.ClearColor(0, 0, 0, 1)
}

func (g *GameView) Update(){
	isRender := g.director.nes.Run()

	if isRender{
		rgba := g.director.nes.Buffer()

		gl.BindTexture(gl.TEXTURE_2D, g.texture)

		// porting vram to texture
		setTexture(rgba)

		// draw actually in window
		drawBuffer(g.director.window)
		gl.BindTexture(gl.TEXTURE_2D, 0)

		// 本当にここでいいかは不明
		g.director.window.SwapBuffers()
	}
}


func createTexture() uint32{
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

func setTexture(i *image.RGBA) {
	size := i.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(i.Pix))
}

func drawBuffer(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / 256
	s2 := float32(h) / 240
	f := float32(1)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}
