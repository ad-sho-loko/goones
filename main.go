package main

import (
	"fmt"
	"github.com/ad-sho-loko/goones/nes"
	"github.com/ad-sho-loko/goones/ui"
	"os"
)

func usage(){
	fmt.Println("no rom files specified or found")
}

func main(){

	if len(os.Args) < 2{
		usage()
		os.Exit(1)
	}

	m, err:= nes.NewCassette(os.Args[1])
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	n := nes.NewNes(m)
	ui.RunUi(n)
}
