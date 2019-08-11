package main

import (
	"fmt"
	"go_nes/src/nes"
	"go_nes/src/ui"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func prof(){
	l, err := net.Listen("tcp", ":52362")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on %s\n", l.Addr())
	go http.Serve(l, nil)
}

func usage(){
	fmt.Println("gooby [.nes file]")
}

func main(){
	// prof()

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
