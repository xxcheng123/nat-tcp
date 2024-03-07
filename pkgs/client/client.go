package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/xxcheng123/nat-tcp/pkgs/ip"
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
	_, err = conn.Write([]byte("nat_info"))
	if err != nil {
		return nil, err
	}
	natInfo, err := read(conn)
	c.natInfos[localPort] = natInfo
	go func() {
		var failTimes = 0
		for {
			_, err = conn.Write([]byte("nat_info"))
			natInfo, err = read(conn)
			if err == nil {
				failTimes = 0
				c.natInfos[localPort] = natInfo
			} else {
				failTimes++
				fmt.Println(err)
			}
			time.Sleep(time.Second * 1)
			if failTimes == 3 {
				fmt.Println("网络错误，开始重试")
				_ = conn.(*net.TCPConn).SetLinger(0)
				<-time.After(time.Duration(10) * time.Second)
				conn.Close()
				for i := 0; i < 3; i++ {
					conn, err = dialer.Dial("tcp", c.RemoteAddr)
					if err != nil {
						time.Sleep(time.Second * 5)
					} else {
						break
					}
				}
				if err != nil {
					break
				} else {
					failTimes = 0
				}
			}
		}
		fmt.Println("关闭链接")
	}()
	return natInfo, err
}

func read(conn net.Conn) (*NatInfo, error) {
	buf := make([]byte, 1024)
	err := conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	if err != nil {
		return nil, err
	}
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
	fmt.Printf("%s %s:%d\n", time.Now().String(), info.Host, info.Port)
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
