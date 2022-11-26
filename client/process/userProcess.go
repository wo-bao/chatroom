package process

import (
	"chatRoom/client/utils"
	"chatRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type UserProcess struct {

}

func (u *UserProcess)Login(id, psd string) (err error) {

	conn, err := net.Dial("tcp", "10.16.16.176:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return err
	}

	defer conn.Close()

	var mes message.Message
	mes.Type = message.LoginMesType

	var loginMes message.LoginMes
	loginMes.UserId = id
	loginMes.UserPsd = psd

	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	n, err := conn.Write(buf[0:4])
	if n !=4 || err != nil {
		fmt.Println("conn.Write(bytes) fail=", err)
		return
	}
	//fmt.Println("mes.len sends suc", pkgLen)

	n, err = conn.Write(data)
	if n !=len(data) || err != nil {
		fmt.Println("conn.Write(data) fail=", err)
		return
	}
	//fmt.Println("mes sends suc")

	t := &utils.Transfer{
		Conn: conn,
	}
	mes, err = t.ReadPkg() //读返回的信息
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}

	var loginResMes message.LoginResMes

	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {

		CurUser.Conn = conn
		CurUser.UserId = id
		CurUser.UserStatus = message.UserOnline

		fmt.Println("login success")
		fmt.Println("online users:")
		for _, v := range loginResMes.Users {
			if v == id {
				continue
			}
			fmt.Println(v)
			user := &message.User{
				UserId: v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}

		go ServerProcessMes(conn)

		for {
			flag := ShowMenu()
			if flag == -1 {
				break
			}
		}
	} else {
		fmt.Println(loginResMes.Error)
	}

	return
}

func (u *UserProcess)Register(id, psd, name string) (err error) {
	conn, err := net.Dial("tcp", "10.16.16.176:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return err
	}
	defer conn.Close()

	var mes message.Message
	mes.Type = message.RegisterMesType

	var registerMes message.RegisterMes
	registerMes.UserId = id
	registerMes.UserPsd = psd
	registerMes.UserName = name

	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	n, err := conn.Write(buf[0:4])
	if n !=4 || err != nil {
		fmt.Println("conn.Write(bytes) fail=", err)
		return
	}
	fmt.Println("mes.len sends suc", pkgLen)

	n, err = conn.Write(data)
	if n !=len(data) || err != nil {
		fmt.Println("conn.Write(data) fail=", err)
		return
	}
	fmt.Println("mes sends suc")

	t := &utils.Transfer{
		Conn: conn,
	}
	mes, err = t.ReadPkg() //读返回的信息
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}

	var registerResMes message.RegisterResMes

	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 100 {
		fmt.Println("register success,please login to get service.")
		return
	} else {
		fmt.Println(registerResMes.Error)
		err = errors.New(registerResMes.Error)
		return
	}
}