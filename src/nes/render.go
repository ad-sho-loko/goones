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
	line              int
	tiles             []*Tile
	backgroundPalette [16]color.RGBA
	spritePalette     [16]color.RGBA
	img               *image.RGBA
	sprites           [64]*Sprite
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
	return r.img
}

func (r *Renderer) render(){
	r.renderBackground(r.tiles)
	r.renderSprites(r.sprites)
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
			paletteIdx := tile.paletteId * 4 + int(tile.bytes[i][j])
			rgba := r.backgroundPalette[paletteIdx]
			x := tileX + j
			y := tileY + i
			r.img.SetRGBA(x, y, rgba)
			// if x >= 0 && 0xFF >= x && y >= 0 && y <= 224 { }
		}
	}
}

func (r *Renderer) renderSprites(sprites [64]*Sprite){
	for _, s := range sprites{
		if s != nil{
			r.renderSprite(s)
		}
	}
}

func (r *Renderer) reverse(b [8][8]byte, isHorizontal bool) [8][8]byte{
	var buf [8][8]byte

	if isHorizontal {
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				buf[i][j] = b[i][7-j]
			}
		}
		return buf
	}

	// vertical
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			buf[i][j] = b[7-i][j]
		}
	}

	return buf
}

func (r *Renderer) renderSprite(sprite *Sprite){
	if sprite.isUseBg{
		return
	}

	if sprite.isHorizontalReverse{
		sprite.bytes = r.reverse(sprite.bytes, true)
	}else if sprite.isVerticalReverse{
		sprite.bytes = r.reverse(sprite.bytes, false)
	}

	for i := 0; i < 8; i++ {
		for j:= 0; j < 8; j++ {

			if sprite.bytes[i][j] == 0 {
				continue
			}

			paletteIdx := int(sprite.paletteId) * 4 + int(sprite.bytes[i][j])
			rgba := r.spritePalette[paletteIdx]
			r.img.SetRGBA(int(sprite.x)+j, int(sprite.y)+i, rgba)

		}
	}
}