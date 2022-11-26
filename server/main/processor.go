package main

import (
	"chatRoom/common/message"
	"chatRoom/server/utils"
	"fmt"
	"io"
	"net"
	process2 "chatRoom/server/process"
	"strings"
)

type Processor struct {
	Conn net.Conn
}

func (p *Processor)ServerProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		u := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = u.ServerProcessLogin(mes)
		return
	case message.RegisterMesType:
		u := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = u.ServerProcessRegister(mes)
		return
	case message.SmsMesType:
		u := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = u.ServerProcessGroupMes(mes)
		return
	case message.LogoffMesType:
		u := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = u.ServerProcessLogoff(mes)
		return
	case message.PrivateMesType:
		u := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = u.ServerProcessPrivateMes(mes)
		if err != nil {
			errNew := u.ServerProcessErrInPri(err)
			fmt.Println("server inner err=", errNew)
		}
		err = nil
		return
	default:
		fmt.Println("MesType is not existing...")
	}
	return
}

func (p *Processor)Process3() (err error) {
	t := utils.Transfer{
		Conn: p.Conn,
	}
	for {
		mes, err := t.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("current conn is break")
				return err
			} else{
				if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host"){
					u := &process2.UserProcess{
						Conn: t.Conn,
					}
					err := u.SendLogoffForciblyMes()
					if err != nil {
						fmt.Println("Distribute Logoff Forcibly Mes err=",err)
					}
					return err
				}
				fmt.Println("mes readPkg err=", err)
				return err
			}
		}
		fmt.Println(mes)
		err = p.ServerProcessMes(&mes)
		if err != nil {
			fmt.Println("err = ",err)
			return err
		}
	}
	return
}

