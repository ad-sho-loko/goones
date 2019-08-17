package nes

import (
	"log"
	"os"
)

type word uint16

func abort(format string, v ...interface{}){
	log.Fatalf(format, v)
	os.Exit(1)
}

