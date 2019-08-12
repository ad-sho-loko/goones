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

func (r *Renderer) renderTile(tile *Tile, tileX, tileY int){
	offsetX := tile.scrollX % 8
	// offsetY := tile.scrollY % 8
	for i := 0; i < 8; i++ {
		for j:= 0; j < 8; j++ {
			paletteIdx := int(tile.paletteId) * 4  + int(tile.bytes[i][j])
			rgba := r.backgroundPalette[paletteIdx]
			x := tileX + j - int(offsetX)
			y := tileY + i // - int(offsetY)

			if x >= 0 && x <= 0xFF && y >= 0 && y < 240 {
				r.img.SetRGBA(x, y, rgba)
			}
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

	// fix : 右端のSpriteが更新されないバグあり
	for i := 0; i < 8; i++ {
		for j:= 0; j < 8; j++ {
			/*
			if sprite.bytes[i][j] == 0 {
				continue
			}
			*/

			paletteIdx := int(sprite.paletteId) * 4 + int(sprite.bytes[i][j])
			rgba := r.spritePalette[paletteIdx]
			x := int(sprite.x) + j
			y := int(sprite.y) + i

			r.img.SetRGBA(x, y, rgba)
		}
	}
}