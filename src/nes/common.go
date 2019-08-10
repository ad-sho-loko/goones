package nes

type Word uint16

type Mem interface{
	load(addr Word) byte
	store(addr Word, b byte)
	slice(begin int, end int) []byte
}

func abort(format string, vv ...interface{}){
	// log.Fatalf(format, vv)
	// os.Exit(1)
}