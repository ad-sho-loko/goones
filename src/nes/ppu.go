package nes

import (
	"image/color"
)

type Ppu struct {
	PpuAddr word // 0x2006
	PpuData byte // 0x2007
	cycle   uint64
	ram     Mem
	bus     *Bus
	renderer *Renderer
}

func NewPpu(bus *Bus, chrRom []byte, r *Renderer) *Ppu{
	return &Ppu{
		PpuAddr:0x00,
		PpuData:0x00,
		cycle:0,
		ram:NewRamInit(0x4000, chrRom),
		bus:bus,
		renderer:r,
	}
}

func (p *Ppu) writePpuAddr(b byte){
	p.PpuAddr = p.PpuAddr << 8 | word(b)
}

func (p *Ppu) writePpuData(b byte){
	// fmt.Printf("[ppu] 0x%x => 0x%x(%d) \n", p.PpuAddr, b, b)
	p.ram.store(p.PpuAddr, b)
	p.PpuAddr++
}

func (p *Ppu) readPpuData() byte{
	b := p.ram.load(p.PpuAddr)
	p.PpuAddr++
	return b
}

func (p *Ppu) run(cycle uint64) bool{
	p.cycle+= cycle * 3

	if p.cycle >= 341 {
		p.cycle -= 341
		p.renderer.line++

		if p.renderer.line <= 240 && p.renderer.line % 8 == 0 {
			p.buildBackground((p.renderer.line - 1) / 8, p.renderer.tiles)
		}

		if p.renderer.line == 262 {
			p.renderer.palette = p.getPalette()
			p.renderer.line = 0
			return true
		}
	}

	return false
}

type Tile struct {
	paletteId int
	sprite [][]byte
}

// build the background in one line.
func (p *Ppu) buildBackground(y int, renderTiles []*Tile){
	// 30 loops in outer methods.
	for x:=0; x<32; x++{
		sprite, palleteId := p.buildTile(x, y)
		renderTiles[y*32+x] = &Tile{
			sprite:sprite,
			paletteId:palleteId,
		}
	}
}

// Tile is 8px * 8px
func (p *Ppu) buildTile(x, y int)([][]byte, int){
	spriteId := p.getSpriteId(x, y)
	blockId := p.getBlockId(x, y)
	attr := p.getAttribute(x, y)
	paletteId := (attr >> uint(blockId) * 2) & 0x03
	// fmt.Printf("(%d,%d) blockId:%d, spriteId:%d attr:%d palleteId:%d\n", x, y, blockId, spriteId, attr, paletteId)
	return p.buildSprite(spriteId), paletteId
}

func (p *Ppu) getAttribute(x, y int) int{
	// 0x2000 -  => ネームテーブル
	// 0x23C0 -  => 属性テーブル
	addr := word(int(x / 4) + (int(y / 4) * 8) + 0x2000 + 0x03C0)
	return int(p.ram.load(addr))
}

// blockId sets up as follow:
// 16 tiles * 15 tiles = 260
// (0)(1)
// (2)(3)
func (p *Ppu) getBlockId(x, y int) int{
	return int((x % 4) / 2) + (int((y % 4) / 2)) * 2
}

// spriteId is which element in name table use in that sprite
// スプライトIDはタイルにどのスプライトを適用させるか。0x2000以降に入っている。
func (p *Ppu) getSpriteId(x, y int) int{
	// https://github.com/bokuweb/flownes/blob/3603b6d05ebf37d55b4b44236cc124c53667ce7b/src/ppu/index.js#L229
	addr := word(y * 32 + x + 0x2000)
	return int(p.ram.load(addr))
}

// SPRITE is 8px * 8px (built by 64bit + 64bit)
func (p *Ppu) buildSprite(spriteId int) [][]byte{
	sprite := new2DimArray(8,8)
	var i, j word
	for i = 0; i<16; i++{
		for  j = 0; j<8; j++{
			addr := word(spriteId) * 16 + i
			b := p.ram.load(addr) // load pattern table
			if b & (0x80 >> j) != 0x00{
				sprite[i%8][j] += 0x01 << uint(i/8) // 0, 1, 3
			}
		}
	}

	return sprite
}

func new2DimArray(maxX, maxY int) [][]byte{
	outer := make([][]byte, maxY)
	for i:=0; i<maxY; i++{
		outer[i] = make([]byte, maxX)
	}
	return outer
}

func (p *Ppu) getPalette() [16]color.RGBA{
	var currentPalette [16]color.RGBA
	for i, b := range p.ram.slice(0x3F00, 0x3F10){
		currentPalette[i] = SystemPalette[b]
	}
	return currentPalette
}
