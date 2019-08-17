package main

import (
	"fmt"
	"go_nes/src/nes"
	"go_nes/src/ui"
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
