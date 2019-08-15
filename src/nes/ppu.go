package nes

import (
	"fmt"
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
	spriteBuffer [64]*Sprite
	spriteRam [64][4]byte // 0=y, 1=idx, 2=attr, 3=x
}

func NewPpu(bus *Bus, chrRom []byte, r *Renderer, isHorizontalMirror bool) *Ppu{
	return &Ppu{
		PpuCtrl:            0x00,
		PpuMask:            0x00,
		PpuStatus:          0x00,
		OamAddr:            0,
		OamData:            0,
		scrollFirst:        true,
		PpuScrollX:         0,
		PpuScrollY:         0,
		PpuAddr:            0x00,
		PpuData:            0x00,
		cycle:              0,
		ram:                NewRamInit(0x4000, chrRom),
		bus:                bus,
		renderer:           r,
		isHorizontalMirror: isHorizontalMirror,
		spriteBuffer:       [64]*Sprite{},
		spriteRam:          [64][4]byte{},
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

// 0x2002
func (p *Ppu) readPpuStatus() byte{
	p.scrollFirst = true // reset scroll register($0x2006)
	b := p.PpuStatus
	p.clearVblank()
	// p.unsetSpriteHit()
	return b
}

// $0x2003
func (p *Ppu) writeOamAddr(b byte){
	p.OamAddr = b
}

// $0x2004
func (p *Ppu) readOamData() byte{
	return p.OamData
}

func (p *Ppu) writeOamData(b byte){
	if p.OamAddr % 4 == 0{
		p.spriteRam[p.OamAddr / 4][0] = b
	} else if p.OamAddr %4 == 1{
		p.spriteRam[p.OamAddr / 4][1] = b
	} else if p.OamAddr %4 == 2 {
		p.spriteRam[p.OamAddr / 4][2] = b
	} else{
		p.spriteRam[p.OamAddr / 4][3] = b
	}
	p.OamAddr++
}


// $0x2005
func (p *Ppu) writePpuScroll(b byte){
	if p.scrollFirst {
		p.PpuScrollX = b
	} else{
		p.PpuScrollY = b
	}

	p.scrollFirst = !p.scrollFirst
}

// $0x2006
func (p *Ppu) writePpuAddr(b byte){
	p.PpuAddr = p.PpuAddr << 8 | word(b)
}

// $0x2007
func (p *Ppu) readPpuData() byte{
	addr := p.calcVramAddr()
	fmt.Printf("[read] 0x%x => 0x%x\n", p.PpuAddr, addr)
	p.PpuAddr += p.getIncrementCount()
	return p.ram.load(addr)
}

func (p *Ppu) writePpuData(b byte){
	addr := p.calcVramAddr()
	fmt.Printf("[write] 0x%x => 0x%x\n", p.PpuAddr, addr)
	p.ram.store(addr, b)
	p.PpuAddr += p.getIncrementCount()
}

func (p *Ppu) clearVblank(){
	p.PpuStatus &= 0x7F
}

func (p *Ppu) calcVramAddr() word{

	if p.PpuAddr == 0x3F10 || p.PpuAddr == 0x3F14 || p.PpuAddr == 0x3F18 || p.PpuAddr == 0x3F1C{
		// $3F10/$3F14/$3F18/$3F1C are mirror of $3F00/$3F04/$3F08/$3F0C.
		return p.PpuAddr - 0x10
	}

	/*
	// mirroring
	if !p.isHorizontalMirror && p.PpuAddr >= 0x2800 && p.PpuAddr < 0x2FFF {
		return p.PpuAddr - 0x0800
	}
	*/

	if p.PpuAddr >= 0x3000 && p.PpuAddr < 0x3F00 {
		// 0x3000 - 0x3EFF is mirror of 0x2000 - 0x2EFF
		return p.PpuAddr - 0x1000
	}

	if p.PpuAddr >= 0x3F20 && p.PpuAddr <= 0x3FFF {
		// 0x3F20 - 0x3FFF is mirror of 0x3F00 - 0x3F1F
		return p.PpuAddr - p.PpuAddr % 0x20
	}

	if p.PpuAddr > 0x4000{
		// Valid addresses are $0000-$3FFF; higher addresses will be mirrored down.
		return p.PpuAddr % 0x4000
	}

	return p.PpuAddr
}

func (p *Ppu) enterVblank(){
	p.PpuStatus |= 0x80
	if p.isAbleNmiVblank(){
		p.bus.cpu.InterruptNmi()
	}
}

func (p *Ppu) leaveVblank() {
	p.renderer.sprites = p.spriteBuffer
	p.renderer.backgroundPalette = p.getBackgroundPalette()
	p.renderer.spritePalette = p.getSpritePalette()
	p.renderer.line = 0
	p.clearVblank()
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
	zeroSpriteY := p.spriteRam[0][0]
	return p.renderer.line == int(zeroSpriteY) && p.isBackgroundEnable() && p.isSpriteEnable()
}

func (p *Ppu) hitSprite(){
	p.PpuStatus |= 0x40
}

func (p *Ppu) endHitSprite(){
	p.PpuStatus &= 0xBF
}

func (p *Ppu) run(cycle uint64) bool{
	p.cycle+= cycle * 3

	if p.renderer.line == 0{
		for i, r := range p.spriteRam{
			p.spriteBuffer[i] = p.getSprite(r[0], r[1], r[2], r[3])
		}
	}

	if p.cycle >= 341 {
		p.cycle -= 341
		p.renderer.line++

		// 0爆弾
		if p.hasHitSprite(){
			p.hitSprite()
		}

		// 262本のスキャンライン
		// 1 - 240 => bg/spritesの描画
		// 240 - 262 => vblank期間。vramの書き換えができる.
		if p.renderer.line <= 240 && p.renderer.line % 8 == 0 {
			p.buildBackground(p.renderer.line - 1, p.renderer.tiles)
		}

		// Start Vblank
		if p.renderer.line == 241{
			p.enterVblank()
		}

		// End Vblank
		if p.renderer.line == 262 {
			p.leaveVblank()
			p.endHitSprite()
			// p.bus.cpu.intrrupt = nil
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
		// Horizontal offsets range from 0 to 255. "Normal" vertical offsets range from 0 to 239,
		// while values of 240 to 255 are treated as -16 through -1 in a way,
		// but tile data is incorrectly fetched from the attribute table.
		// By changing the values here across several frames and writing tiles to newly revealed areas of the nametables,
		// one can achieve the effect of a camera panning over a large background.
		adjusted -= 256
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
	if (addr >= 0x2400 && addr < 0x2800) || addr >= 0x2C00{
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

func (p *Ppu) getSprite(y, spriteId, attr, x byte) *Sprite{
	bytes := p.buildSprite(int(spriteId), p.fetchSpriteChrTable())

	/*
	// INFO: Offset sprite Y position, because First and last 8line is not rendered.
	if y < 8{
		return &Sprite{}
	}
	y-=8
	*/

	return &Sprite{
		// NOTE : Sprite data is delayed by one scanline;
		// you must subtract 1 from the sprite's Y coordinate before writing it here.
		// Hide a sprite by writing any values in $EF-$FF here.
		// Sprites are never displayed on the first line of the picture,
		// and it is impossible to place a sprite partially off the top of the screen.
		y:y-1,
		x:x,
		bytes:bytes,
		isVerticalReverse:attr & 0x80 != 0,
		isHorizontalReverse:attr & 0x40 != 0,
		isUseBg:attr & 0x20 != 0,
		paletteId:attr & 0x03,
	}
}