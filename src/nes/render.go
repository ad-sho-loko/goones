package nes

import (
	"image"
	"image/color"
)

var(
	UpLeft = image.Point{X:0,Y:0}
	DownRight = image.Point{X: 256, Y: 240}
)

type Renderer struct {
	line int
	tiles []*Tile
	palette [16]color.RGBA
	img *image.RGBA
}

func NewRenderer() *Renderer{
	initSystemPallete()
	return &Renderer{
		line:0,
		tiles:make([]*Tile, 30*32),
		// the first image before rendering
		img:image.NewRGBA(image.Rectangle{Min: UpLeft, Max: DownRight}),
	}
}

func (r *Renderer) Buffer() *image.RGBA{
	r.render()
	return r.img
}

func (r *Renderer) render(){
	r.renderBackground(r.tiles)
}

func (r *Renderer) renderBackground(background []*Tile){
	for i := 0; i < len(background); i++ {
		x := (i % 32) * 8
		y := int(i / 32) * 8
		r.renderTile(background[i], x, y)
	}
}

// render 1 tile (8px * 8px)
// TileXはそのタイルの左上のx座標を指している
func (r *Renderer) renderTile(tile *Tile, tileX, tileY int){
	for i := 0; i < 8; i++ {
		for j:= 0; j < 8; j++ {
			paletteIdx := tile.paletteId * 4 + int(tile.sprite[i][j])
			rgba := r.palette[paletteIdx]
			x := tileX + j
			y := tileY + i
			r.img.SetRGBA(x, y, rgba)
			// if x >= 0 && 0xFF >= x && y >= 0 && y <= 224 { }
		}
	}
}

var SystemPalette [64]color.RGBA

func initSystemPallete() {
	colors := []uint32{
		0x666666, 0x002A88, 0x1412A7, 0x3B00A4, 0x5C007E, 0x6E0040, 0x6C0600, 0x561D00,
		0x333500, 0x0B4800, 0x005200, 0x004F08, 0x00404D, 0x000000, 0x000000, 0x000000,
		0xADADAD, 0x155FD9, 0x4240FF, 0x7527FE, 0xA01ACC, 0xB71E7B, 0xB53120, 0x994E00,
		0x6B6D00, 0x388700, 0x0C9300, 0x008F32, 0x007C8D, 0x000000, 0x000000, 0x000000,
		0xFFFEFF, 0x64B0FF, 0x9290FF, 0xC676FF, 0xF36AFF, 0xFE6ECC, 0xFE8170, 0xEA9E22,
		0xBCBE00, 0x88D800, 0x5CE430, 0x45E082, 0x48CDDE, 0x4F4F4F, 0x000000, 0x000000,
		0xFFFEFF, 0xC0DFFF, 0xD3D2FF, 0xE8C8FF, 0xFBC2FF, 0xFEC4EA, 0xFECCC5, 0xF7D8A5,
		0xE4E594, 0xCFEF96, 0xBDF4AB, 0xB3F3CC, 0xB5EBF2, 0xB8B8B8, 0x000000, 0x000000,
	}

	for i, c := range colors {
		r := byte(c >> 16)
		g := byte(c >> 8)
		b := byte(c)
		SystemPalette[i] = color.RGBA{R: r, G: g, B: b, A: 0xFF}
	}
}

