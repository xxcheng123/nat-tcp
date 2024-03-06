package main

import (
	"flag"
	"fmt"
	"nat-tcp/pkgs/client"
	"os"
	"os/signal"
)

var remoteAddr string
var localPort int

func init() {
	flag.StringVar(&remoteAddr, "addr", "localhost:9933", "服务器地址")
	flag.IntVar(&localPort, "port", 9986, "本地端口")
}

func main() {
	natClient := client.NewClient(remoteAddr)
	info, err := natClient.Call(1999)
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		fmt.Println("bye")
	}
}
