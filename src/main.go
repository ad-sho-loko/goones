package main

import (
	"fmt"
	"go_nes/src/nes"
	"go_nes/src/ui"
	"net"
	"net/http"
	// _ "net/http/pprof"
)

func prof(){
	l, err := net.Listen("tcp", ":52362")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on %s\n", l.Addr())
	go http.Serve(l, nil)
}

func main(){
	// prof()
	m, err:= nes.NewCassette("../resource/helloworld/sample1.nes")

	if err != nil{
		fmt.Println(err)
	}
	n := nes.NewNes(m)
	ui.RunUi(n)
}
