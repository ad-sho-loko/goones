package nes

import (
	"log"
	"os"
)

type word uint16

type Mem interface{
	load(addr word) byte
	store(addr word, b byte)
	slice(begin int, end int) []byte
}

func abort(format string, v ...interface{}){
	log.Fatalf(format, v)
	os.Exit(1)
}

