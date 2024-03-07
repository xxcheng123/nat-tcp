package main_test

import (
	"fmt"
	"net"
	"testing"
	"time"
)

var tP1 = 12421
var tP2 = 13512

func TestCloseTCP(t *testing.T) {
	dialer1 := net.Dialer{
		LocalAddr: &net.TCPAddr{
			Port: tP1,
		},
	}
	dialer2 := net.Dialer{
		LocalAddr: &net.TCPAddr{
			Port: tP2,
		},
	}
	t.Run("old", func(tt *testing.T) {
		conn, err := dialer1.Dial("tcp", "www.qq.com:80")
		if err != nil {
			tt.Error(err)
			return
		}
		conn.Close()
		time.Sleep(time.Second * 10)
		_, err = dialer1.Dial("tcp", "www.qq.com:80")
		if err != nil {
			tt.Error(err)
			return
		}
		fmt.Println(tt.Name(), "测试通过")
		conn.Close()
	})
	t.Run("new", func(tt *testing.T) {
		conn, err := dialer2.Dial("tcp", "www.qq.com:80")
		if err != nil {
			tt.Error(err)
			return
		}
		err = conn.(*net.TCPConn).SetLinger(0)
		if err != nil {
			tt.Error(err)
			return
		}
		conn.Close()
		time.Sleep(time.Second * 10)
		fmt.Println(t.Name(), "第一阶段通过")
		_, err = dialer2.Dial("tcp", "www.qq.com:80")
		_ = conn.(*net.TCPConn).SetLinger(0)
		if err != nil {
			tt.Error(err)
		}
		if err != nil {
			tt.Error(err)
			return
		}
		fmt.Println(t.Name(), "第二阶段通过")
	})
}
