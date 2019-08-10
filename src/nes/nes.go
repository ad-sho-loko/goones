package nes

import (
	"errors"
	"image"
)

type Nes struct {
	cassette Ines
	cpu      *Cpu
	ppu      *Ppu
	bus      *Bus
	r        *Renderer
}

func NewNes(cassette Ines) *Nes {
	wram := NewRam(0x800)
	bus := NewBus(wram, cassette.PrgRom())
	cpu := NewCpu(bus)
	ppu := NewPpu(bus, cassette.ChrRom())
	bus.cpu = cpu
	bus.ppu = ppu
	return &Nes{
		cassette: cassette,
		cpu:      cpu,
		ppu:      ppu,
		bus:      bus,
		r:        NewRenderer(),
	}
}

func (n *Nes) isSetCassette() bool {
	return n.cassette != nil
}

func (n *Nes) Init() error {
	if !n.isSetCassette() {
		return errors.New("cassette must be set")
	}

	n.cpu.PC = 0x8000
	return nil
}

func (n *Nes) Run() bool {
	pc := n.cpu.PC

	// decode
	b := n.bus.Load(pc)
	inst := n.cpu.decode(b)
	cycle := inst.cycle
	wd := n.cpu.solveAddrMode(inst.addrMode)

	// for debug
	// n.cpu.dump(b, wd, inst.mnemonic, inst.addrMode)

	// execute
	n.cpu.advance(inst.addrMode)
	n.cpu.execute(inst, wd)

	n.cpu.cycle += cycle

	// TODO : PPUの処理を移動させる.
	n.ppu.cycle += cycle * 3

	if n.ppu.cycle >= 341 {
		n.ppu.cycle -= 341
		n.r.line++

		if n.r.line <= 240 && n.r.line % 8 == 0 {
			n.ppu.buildBackground((n.r.line - 1) / 8, n.r.tiles)
		}

		if n.r.line == 262 {
			n.r.palette = n.ppu.getPalette()
			n.r.line = 0
			return true
		}
	}

	// is nesessary for rendering?
	return false
}

func (n *Nes) Buffer() *image.RGBA {
	return n.r.Buffer()
}
