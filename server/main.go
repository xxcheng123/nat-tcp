package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"nat-tcp/pkgs/ip"
	"net"
	"strconv"
)

var port = 9933

var SIGNAL_CLOSE = []byte("close")

func init() {
	flag.IntVar(&port, "port", 9933, "服务端口")
}
func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		panic(err)
	}
	defer l.Close()
	fmt.Printf("running in :%d\n", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	write(conn)
	bs := make([]byte, 1024)
	var retryTimes = 0
	for {
		_, err := conn.Read(bs)
		if err != nil {
			retryTimes++
			continue
		} else if retryTimes > 3 {
			break
		}
		if bytes.Compare(bs, SIGNAL_CLOSE) == 0 {
			break
		}
		retryTimes = 0
		write(conn)
	}
}

func write(conn net.Conn) {
	_host, _port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	fmt.Printf("receive reqeust [%s]\n", conn.RemoteAddr().String())
	clientPublicHost := _host
	clientPublicPort, _ := strconv.Atoi(_port)
	body, _ := json.Marshal(ip.Info{Host: clientPublicHost, Port: clientPublicPort})
	_, _ = conn.Write(body)
}