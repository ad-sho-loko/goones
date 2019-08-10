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
}

func NewNes(cassette Ines) *Nes {
	wram := NewRam(0x800)
	renderer := NewRenderer()
	bus := NewBus(wram, cassette.PrgRom())
	cpu := NewCpu(bus)
	ppu := NewPpu(bus, cassette.ChrRom(), renderer)
	bus.cpu = cpu
	bus.ppu = ppu
	return &Nes{
		cassette: cassette,
		cpu:      cpu,
		ppu:      ppu,
		bus:      bus,
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

	return n.ppu.run(cycle * 3)
}

func (n *Nes) Buffer() *image.RGBA {
	return n.ppu.renderer.Buffer()
}
