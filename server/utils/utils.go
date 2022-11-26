package utils

import (
	"chatRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Transfer struct {
	Conn net.Conn
	Buf [4096]byte // 传输使用的缓冲
}


func (t *Transfer)ReadPkg() (mes message.Message, err error) {
	//buf := make([]byte, 4096)
	n, err := t.Conn.Read(t.Buf[:4]) //读不到东西会阻塞,conn关闭了就不会阻塞
	if n!=4 || err != nil {
		fmt.Println("conn.Read err=",err)
		return
	}
	fmt.Println("conn.Read buf=", t.Buf[:4])

	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(t.Buf[:4])

	n, err = t.Conn.Read(t.Buf[:pkgLen])
	if n!= int(pkgLen) || err!= nil {
		fmt.Println("conn.Read err=", err)
		return
	}

	err = json.Unmarshal(t.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
	}

	return mes, err
}

func (t *Transfer)WritePkg(data []byte) (err error) {

	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(t.Buf[:4], pkgLen)

	n, err := t.Conn.Write(t.Buf[:4])
	if n !=4 || err != nil {
		fmt.Println("conn.Write(bytes) fail=", err)
		return
	}
	fmt.Println("mes.len sends suc", pkgLen)

	n, err = t.Conn.Write(data)
	if n !=int(pkgLen) || err != nil {
		fmt.Println("conn.Write(data) fail=", err)
		return
	}
	fmt.Println("mes sends suc", pkgLen)

	return
}
