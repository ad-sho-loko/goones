package main

import (
	"fmt"
	"go_nes/src/nes"
	"go_nes/src/ui"
	"os"
)

func usage(){
	fmt.Println("gooby [.nes file]")
}

func main(){

	if len(os.Args) < 2{
		usage()
		os.Exit(1)
	}
	m, err:= nes.NewCassette("../resource/roms/" + os.Args[1])

	if err != nil{
		fmt.Println(err)
	}
	n := nes.NewNes(m)
	ui.RunUi(n)
}
