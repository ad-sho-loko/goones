package nes

import (
	"image/color"
)

type Ppu struct {
	// Core
	PpuCtrl byte   // 0x2000
	PpuMask byte   // 0x2001
	PpuStatus byte // 0x2002
	OamAddr byte   // 0x2003
	OamData byte   // 0x2004
	PpuAddr word   // 0x2006
	PpuData byte   // 0x2007
	cycle   uint64
	ram     Mem
	bus     *Bus
	renderer *Renderer

	// Sprite RAM
	spriteId int
	spriteCounter int
	spriteBuffer [64]*Sprite
	y byte
	tileIndex byte
	attr byte
	x byte
}

func NewPpu(bus *Bus, chrRom []byte, r *Renderer) *Ppu{
	return &Ppu{
		PpuCtrl:0x00,
		PpuMask:0x00,
		PpuStatus:0x00,
		PpuAddr:0x00,
		PpuData:0x00,
		cycle:0,
		ram:NewRamInit(0x4000, chrRom),
		bus:bus,
		renderer:r,
	}
}

// VBlank時にNMI割込の発生(1:On, 0:Off)
func (p *Ppu) isAbleNmiVblank() bool{
	return p.PpuCtrl & 0x80 != 0
}

func (p *Ppu) writePpuCtrl(b byte){
	p.PpuCtrl = b
}

func (p *Ppu) writePpuMask(b byte){
	p.PpuMask = b
}

func (p *Ppu) writeOamAddr(b byte){
	p.OamAddr = b
}

func (p *Ppu) writeOamData(b byte){

	if p.spriteCounter == 0{
		p.y = b
	}else if p.spriteCounter == 1{
		p.tileIndex = b
	}else if p.spriteCounter == 2{
		p.attr = b
	}else if p.spriteCounter == 3 {
		p.x = b
		p.spriteBuffer[p.spriteId] = p.getSprite(int(p.tileIndex))
		p.spriteId++
		p.spriteId%=64
	}

	p.spriteCounter++
	p.spriteCounter %= 4
}

func (p *Ppu) writePpuAddr(b byte){
	p.PpuAddr = p.PpuAddr << 8 | word(b)
}

func (p *Ppu) writePpuData(b byte){
	// fmt.Printf("[ppu] 0x%x => 0x%x(%d) \n", p.PpuAddr, b, b)
	p.ram.store(p.PpuAddr, b)
	p.PpuAddr++
}

func (p *Ppu) readPpuStatus() byte{
	return p.PpuStatus
}

func (p *Ppu) readOamData() byte{
	return p.OamData
}

func (p *Ppu) readPpuData() byte{
	b := p.ram.load(p.PpuAddr)
	p.PpuAddr++
	return b
}

func (p *Ppu) onVrank(){
	p.PpuStatus |= 0x80
}

func (p *Ppu) notOnVrank() {
	p.PpuStatus &= 0x7F
}

func (p *Ppu) run(cycle uint64) bool{
	p.cycle+= cycle * 3

	if p.cycle >= 341 {
		p.cycle -= 341
		p.renderer.line++

		if p.renderer.line == 241{
			if p.isAbleNmiVblank(){
				p.bus.cpu.InterruptNmi()
			}
		}

		// これ毎回やる必要ないのでは？
		if p.renderer.line > 240 && p.renderer.line <= 262{
			p.onVrank()
		}else{
			p.notOnVrank()
		}

		if p.renderer.line <= 240 && p.renderer.line % 8 == 0 {
			p.buildBackground((p.renderer.line - 1) / 8, p.renderer.tiles)
		}

		if p.renderer.line == 262 {
			p.renderer.sprites = p.spriteBuffer
			p.renderer.backgroundPalette = p.getBackgroundPalette()
			p.renderer.spritePalette = p.getSpritePalette()

			p.renderer.line = 0
			return true
		}
	}

	return false
}

type Tile struct {
	paletteId int
	bytes     [8][8]byte
}

// build the background in one line.
func (p *Ppu) buildBackground(y int, renderTiles []*Tile){
	// 30 loops in outer methods.
	for x:=0; x<32; x++{
		sprite, palleteId := p.buildTile(x, y)
		renderTiles[y*32+x] = &Tile{
			bytes:     sprite,
			paletteId: palleteId,
		}
	}
}

// Tile is 8px * 8px
func (p *Ppu) buildTile(x, y int)([8][8]byte, int){
	spriteId := p.getSpriteId(x, y)
	blockId := p.getBlockId(x, y)
	attr := p.getAttribute(x, y)
	paletteId := (attr >> uint(blockId) * 2) & 0x03
	// fmt.Printf("(%d,%d) blockId:%d, spriteId:%d attr:%d palleteId:%d\n", x, y, blockId, spriteId, attr, paletteId)
	return p.buildSprite(0x0000, spriteId), paletteId
}

func (p *Ppu) getAttribute(x, y int) int{
	// 0x2000 -  => ネームテーブル + 0x23C0 -  => 属性テーブル

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

// spriteId is which element in name table use in that bytes
// スプライトIDはタイルにどのスプライトを適用させるか。0x2000以降に入っている。
func (p *Ppu) getSpriteId(x, y int) int{
	addr := word(y * 32 + x + 0x2000)
	return int(p.ram.load(addr))
}

// SPRITE is 8px * 8px (built by 64bit + 64bit)
func (p *Ppu) buildSprite(offset word, spriteId int) [8][8]byte{
	var sprite [8][8]byte
	var i, j word
	for i = 0; i<16; i++{
		for  j = 0; j<8; j++{
			// calculates the pattern table address.
			addr := offset + word(spriteId) * 16 + i
			b := p.ram.load(addr)
			if b & (0x80 >> j) != 0x00{
				sprite[i%8][j] += 0x01 << uint(i/8) // 0, 1, 3
			}
		}
	}

	return sprite
}

func (p *Ppu) getBackgroundPalette() [16]color.RGBA{
	var currentPalette [16]color.RGBA
	// background palette
	for i, b := range p.ram.slice(0x3F00, 0x3F10){
		currentPalette[i] = SystemPalette[b]
	}
	return currentPalette
}

func (p *Ppu) getSpritePalette() [16]color.RGBA{
	var currentPalette [16]color.RGBA
	// sprite palette
	for i, b := range p.ram.slice(0x3F10, 0x3F20){
		currentPalette[i] = SystemPalette[b]
	}
	return currentPalette
}

type Sprite struct {
	y byte
	x byte
	bytes [8][8]byte
	isVerticalReverse bool
	isHorizontalReverse bool
	isUseBg bool
	paletteId byte
}

func (p *Ppu) getSprite(tileIndex int) *Sprite{
	// tileIndexはどのタイル(.spr)を使うかを決めるもの
	// spriteIdはレンダリングのために連番をつけただけ
	bytes := p.buildSprite(0x1000, tileIndex)

	return &Sprite{
		y:p.y,
		x:p.x,
		bytes:bytes,
		isVerticalReverse:p.attr & 0x80 != 0,
		isHorizontalReverse:p.attr & 0x40 != 0,
		isUseBg:p.attr & 0x20 != 0,
		paletteId:p.attr & 0x03,
	}
}
