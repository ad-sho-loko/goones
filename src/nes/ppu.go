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
	scrollFirst bool  // for 0x2005
	PpuScrollX byte // 0x2005(1)
	PpuScrollY byte // 0x2005(1)
	PpuAddr word   // 0x2006
	PpuData byte   // 0x2007
	cycle   uint64
	ram     Mem
	bus     *Bus
	renderer *Renderer

	// setting
	isHorizontalMirror bool

	// Sprite RAM
	spriteId int
	spriteCounter int
	spriteBuffer [64]*Sprite
	y byte
	tileIndex byte
	attr byte
	x byte
}

func NewPpu(bus *Bus, chrRom []byte, r *Renderer, isHorizontalMirror bool) *Ppu{
	return &Ppu{
		PpuCtrl:0x00,
		PpuMask:0x00,
		PpuStatus:0x00,
		PpuAddr:0x00,
		PpuData:0x00,
		scrollFirst:true,
		cycle:0,
		ram:NewRamInit(0x4000, chrRom),
		bus:bus,
		renderer:r,
		isHorizontalMirror:isHorizontalMirror,
	}
}

// VBlank時にNMI割込の発生(1:On, 0:Off)
func (p *Ppu) isAbleNmiVblank() bool{
	return p.PpuCtrl & 0x80 != 0
}

func (p *Ppu) getIncrementCount() word{
	if p.PpuCtrl & 0x40 == 1{
		return 32
	}else{
		return 1
	}
}

// $0x2000
func (p *Ppu) writePpuCtrl(b byte){
	p.PpuCtrl = b
}

// $0x2001
func (p *Ppu) writePpuMask(b byte){
	p.PpuMask = b
}

func (p *Ppu) writePpuScroll(b byte){
	if p.scrollFirst {
		p.PpuScrollX = b
	} else{
		p.PpuScrollY = b
	}
	p.scrollFirst = !p.scrollFirst
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
	addr := p.calcVramAddr(p.PpuAddr)
	// fmt.Printf("0x%x => 0x%x\n", p.PpuAddr, addr)
	p.ram.store(addr, b)
	p.PpuAddr += p.getIncrementCount()
}

// 0x2002
func (p *Ppu) readPpuStatus() byte{
	// reset scroll x and scroll y
	p.scrollFirst = true
	return p.PpuStatus
}

func (p *Ppu) readOamData() byte{
	return p.OamData
}

func (p *Ppu) readPpuData() byte{
	addr := p.calcVramAddr(p.PpuAddr)
	p.PpuAddr += p.getIncrementCount()
	return p.ram.load(addr)
}

func (p *Ppu) calcVramAddr(addr word) word{
	if p.PpuAddr >= 0x3000 && p.PpuAddr <= 0x3EFF{
		// 0x3000 - 0x3EFF is mirror of 0x2000 - 0x2EFF
		return p.PpuAddr - 0x1000
	}else if p.PpuAddr >= 0x3F20 && p.PpuAddr <= 0x3FFF{
		return p.PpuAddr - 0x0010
	}else{
		return addr
	}
}

func (p *Ppu) onVblank(){
	p.PpuStatus |= 0x80
	if p.isAbleNmiVblank(){
		p.bus.cpu.InterruptNmi()
	}
}

func (p *Ppu) notOnVblank() {
	p.renderer.sprites = p.spriteBuffer
	p.renderer.backgroundPalette = p.getBackgroundPalette()
	p.renderer.spritePalette = p.getSpritePalette()
	p.renderer.line = 0
	p.PpuStatus &= 0x7F
	// not cool
	// p.bus.cpu.unsetBit(Irq)
}

func (p *Ppu) isBackgroundEnable() bool {
	return p.PpuMask & 0x08 != 0
}

func (p *Ppu) isSpriteEnable() bool {
	return p.PpuMask & 0x10 != 0
}

func (p *Ppu) hasHitSprite() bool{
	return p.renderer.line == int(p.y) && p.isBackgroundEnable() && p.isSpriteEnable()
}

func (p *Ppu) run(cycle uint64) bool{
	p.cycle+= cycle * 3

	if p.cycle >= 341 {
		p.cycle -= 341
		p.renderer.line++

		// 0爆弾
		if p.hasHitSprite(){
			p.PpuStatus |= 0x40
		}

		// 262本のスキャンライン
		// 1 - 240 => bg/spritesの描画
		// 240 - 262 => vblank期間。vramの書き換えができる.
		if p.renderer.line <= 240 && p.renderer.line % 8 == 0 {
			p.buildBackground(p.renderer.line - 1, p.renderer.tiles)
		}

		// Start Vblank
		if p.renderer.line == 241{
			p.onVblank()
		}

		// End Vblank
		if p.renderer.line == 262 {
			p.notOnVblank()
			p.PpuStatus &= 0xFF - 0x40
			return true
		}
	}

	return false
}

func (p *Ppu) fetchBgChrTable() word{
	if p.PpuCtrl & 0x10 != 0{
		return 0x1000
	}else{
		return 0x0000
	}
}

func (p *Ppu) fetchSpriteChrTable() word{
	if p.PpuCtrl & 0x08 != 0{
		return 0x1000
	}else{
		return 0x0000
	}
}

func (p *Ppu) fetchNameTableId() int{
	return int(p.PpuCtrl & 0x03)
}

type Tile struct {
	paletteId int
	bytes     [8][8]byte
	scrollX byte
	scrollY byte
}


func (p *Ppu) adjustScrollY() int{
	adjusted := int(p.PpuScrollY) + (int(p.fetchNameTableId() / 2) * 240)

	if p.PpuScrollY >= 240 {
		// スクロールレジスタYは255まで値が格納できるが、
		// 実際に描画するのは240までなので調整する必要がある
		adjusted -= 255
	}

	return adjusted
}

func (p *Ppu) adjustScrollX() int{
	return int(p.PpuScrollX) + (int(p.fetchNameTableId() % 2) * 256)
}

func (p *Ppu) evaluateNameTableOffset(x, y int) int{
	tableIdOffset := 0
	if int(y / 30) % 2 == 1{
		tableIdOffset = 2
	}

	nameTableId := int(x / 32) % 2 + tableIdOffset
	return nameTableId * 0x400
}

// build the background in one line.
func (p *Ppu) buildBackground(line int, renderTiles []*Tile){
	y := int(line / 8)
	tileY := y + int(p.adjustScrollY() / 8)
	rotateY := tileY % 30

	for x:=0; x<32; x++{
		tileX := x + int(p.adjustScrollX() / 8)
		rotatedX := tileX % 32
		nameTableOffset := p.evaluateNameTableOffset(tileX, tileY)

		sprite, palleteId := p.buildTile(rotatedX, rotateY, nameTableOffset)
		renderTiles[y * 32 + x] = &Tile{
			bytes:     sprite,
			paletteId: palleteId,
			scrollX: p.PpuScrollX,
			scrollY: p.PpuScrollY,
		}
	}
}

func (p *Ppu) buildTile(x, y, offset int)([8][8]byte, int){
	spriteId := p.getSpriteId(x, y, offset)
	blockId := p.getBlockId(x, y)
	attr := p.getAttribute(x, y, offset)
	paletteId := (attr >> (word(blockId) * 2)) & 0x03
	// fmt.Printf("(%d,%d) blockId:%d, spriteId:%d attr:%d palleteId:%d\n", x, y, blockId, spriteId, attr, paletteId)
	return p.buildSprite(spriteId, p.fetchBgChrTable()), paletteId
}

func (p *Ppu) getBlockId(x, y int) int{
	return int((x % 4) / 2) + (int((y % 4) / 2)) * 2
}

func (p *Ppu) getAttribute(x, y, offset int) int{
	addr := word(int(x / 4) + (int(y / 4) * 8) + 0x23C0 + offset)
	addr = p.downMirror(addr)
	return int(p.ram.load(addr))
}

// Gets sprite ids from the name table.
func (p *Ppu) getSpriteId(x, y, offset int) int{
	addr := word(y * 32 + x + 0x2000 +  offset)
	addr = p.downMirror(addr)
	return int(p.ram.load(addr))
}

func (p *Ppu) downMirror(addr word) word{
	if !p.isHorizontalMirror {
		return addr
	}

	// Is nametable 1 or 3?
	if addr >= 0x2400 && addr < 0x2800 || addr >= 0x2C00{
		return addr - 0x0400
	}

	// name table 2
	return addr
}

func (p *Ppu) buildSprite(spriteId int, offset word) [8][8]byte{
	var sprite [8][8]byte
	var i, j word
	for i = 0; i<16; i++{
		for  j = 0; j<8; j++{
			addr := word(spriteId) * 16 + i + offset
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
	for i, b := range p.ram.slice(0x3F00, 0x3F10){
		if i % 4 == 0 {
			// 0x3F04, 0x3F08, 0x3C0C are ignored by background.
			// Instead of here, use these values in the sprite palette.
			currentPalette[i] = SystemPalette[p.ram.load(0x3F00)]
		}else{
			currentPalette[i] = SystemPalette[b]
		}
	}
	return currentPalette
}

func (p *Ppu) getSpritePalette() [16]color.RGBA{
	var currentPalette [16]color.RGBA
	for i, b := range p.ram.slice(0x3F10, 0x3F20){
		if i % 4 == 0 {
			// 0x3F10, 0x3F14, 0x3F18, 0x3F1C are mirror of 0x3F00, 0x3F04, 0x3F08, 0x3CFC
			currentPalette[i] = SystemPalette[p.ram.load(word(0x3F00+i))]
		}else{
			currentPalette[i] = SystemPalette[b]
		}
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
	bytes := p.buildSprite(tileIndex, p.fetchSpriteChrTable())

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
