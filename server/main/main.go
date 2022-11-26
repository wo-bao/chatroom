package main

import (
	"chatRoom/server/model"
	"fmt"
	"net"
	"time"
)

// net.Conn 是引用类型
func process(conn net.Conn) {
	defer conn.Close()
	p := &Processor{
		Conn: conn,
	}
	err := p.Process3()
	if err != nil {
		fmt.Println("server and client communication err=", err)
		return
	}
}

func initUserDao() {
	model.MyUserDao = model.NewUserDao(Pool)
}

func main() {
	initPool(8, 0, time.Second * 300, "localhost:6379")
	initUserDao()
	fmt.Println("server is monitoring at 8889 port")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer listen.Close()
	for {
		fmt.Println("waiting clients come to connect")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		go process(conn)
		time.Sleep(10)
	}
}