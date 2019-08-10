package nes

import "testing"

var c *Cpu

func Test_LDA(t *testing.T){
	c.lda(1)
	if c.A != 1{
		t.Fatal()
	}
}

func TestMain(m *testing.M){
	c = &Cpu{}
	p = &Ppu{}
	m.Run()
}
