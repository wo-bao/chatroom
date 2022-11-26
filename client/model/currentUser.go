package model

import (
	"chatRoom/common/message"
	"net"
)

type CurUser struct {
	Conn net.Conn
	message.User
}
