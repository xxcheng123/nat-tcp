package client

import (
	"encoding/json"
	"fmt"
	"github.com/xxcheng123/nat-tcp/pkgs/ip"
	"net"
	"strconv"
	"time"
)

type Client struct {
	conns      map[int]net.Conn
	natInfos   map[int]*NatInfo
	RemoteAddr string
	running    bool
}
type NatInfo struct {
	PublicHost  string
	PublicPort  int
	PrivateHost string
	PrivatePort int
}

func (c *Client) Call(localPort int) (*NatInfo, error) {
	oldNatInfo, ok := c.natInfos[localPort]
	if ok {
		return oldNatInfo, nil
	}
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			Port: localPort,
		},
	}
	conn, err := dialer.Dial("tcp", c.RemoteAddr)
	c.conns[localPort] = conn
	if err != nil {
		return nil, err
	}
	_, err = conn.Write([]byte("hello"))
	if err != nil {
		return nil, err
	}
	natInfo, err := read(conn)
	c.natInfos[localPort] = natInfo
	go func() {
		for {
			_, err = conn.Write([]byte("hello"))
			natInfo, err = read(conn)
			if err == nil {
				c.natInfos[localPort] = natInfo
			}
			time.Sleep(time.Second * 10)
		}
	}()
	return natInfo, err
}

func read(conn net.Conn) (*NatInfo, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	info := &ip.Info{}
	if err = json.Unmarshal(buf[:n], info); err != nil {
		return nil, err
	}
	lh, lp, _ := net.SplitHostPort(conn.LocalAddr().String())
	lp2, _ := strconv.Atoi(lp)
	natInfo := &NatInfo{
		PublicHost:  info.Host,
		PublicPort:  info.Port,
		PrivateHost: lh,
		PrivatePort: lp2,
	}
	fmt.Printf("publicAddr:%s:%d\n", info.Host, info.Port)
	return natInfo, nil
}

func NewClient(remoteAddr string) *Client {
	c := &Client{
		RemoteAddr: remoteAddr,
		conns:      map[int]net.Conn{},
		natInfos:   map[int]*NatInfo{},
	}
	return c
}
