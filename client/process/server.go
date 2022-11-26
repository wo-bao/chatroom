package process

import (
	"chatRoom/client/utils"
	"chatRoom/common/message"
	"encoding/json"
	"fmt"
	"net"
)

func ShowMenu() (flag int) {
	fmt.Printf("------------恭喜 %s 登录成功------------\n", CurUser.UserId)
	fmt.Println("------------1. 显示在线用户------------")
	fmt.Println("------------2. 发送消息   ------------")
	fmt.Println("------------3. 信息列表   ------------")
	fmt.Println("------------4. 退出系统   ------------")
	fmt.Println("------------请选择(1-4)   ------------")
	var key int
	var content string
	sp := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
		case 1:
			fmt.Println("显示在线用户列表")
			outputOnlineUsers()
		case 2:
			var MesType int
			fmt.Println("请输入消息类型: 1-群发消息, 2-私发消息")
			fmt.Scanf("%d\n", &MesType)
			switch MesType {
			case 1:
				fmt.Println("请输入消息:")
				fmt.Scanf("%s\n", &content)
				err := sp.SendGroupMes(content)
				if err != nil {
					fmt.Println("SendGroupMes failed =",err)
				}
			case 2:
				var userID string
				fmt.Println("请输入要发送的用户名:")
				fmt.Scanf("%s\n", &userID)
				if userID == CurUser.UserId {
					fmt.Println("can not send message to yourself! ")
				} else {
					fmt.Println("请输入消息:")
					fmt.Scanf("%s\n", &content)
					err := sp.SendPrivateMes(userID, content)
					if err != nil {
						fmt.Println("SendGroupMes failed =",err)
					}
				}
			}
		case 3:
			fmt.Println("信息列表")
		case 4:
			fmt.Println("退出系统...")
			flag = -1
			err := sp.SendLogoffMes()
			if err != nil {
				fmt.Println("SendLogoffMes failed =",err)
			}
			return
		default:
			fmt.Println("输入错误, 请重新输入")
	}
	return
}

func ServerProcessMes(conn net.Conn) {
	t := utils.Transfer{
		Conn: conn,
	}
	for {
		//fmt.Println("client is reading from server...")
		mes, err := t.ReadPkg()
		if err!= nil {
			fmt.Println("connection is broken, err=", err)
			return
		}
		switch mes.Type {
			case message.NotifyUserStatusMesType:
				var mesNotifyStatus message.NotifyUserStatusMes
				err = json.Unmarshal([]byte(mes.Data), &mesNotifyStatus)
				if err != nil {
					fmt.Println("json.Unmarshal err=", err)
				} else {
					updateUserStatus(&mesNotifyStatus)
				}
			case message.SmsMesType:
				var groupMes message.SmsMes
				err = json.Unmarshal([]byte(mes.Data), &groupMes)
				if err != nil {
					fmt.Println("json.Unmarshal in receive group mes err=", err)
				}
				userid := groupMes.UserId
				fmt.Printf("user %s sends a group mes: \n", userid)
				fmt.Println(groupMes.Content)
			case message.LogoffMesType:
				var mesNotifyLogoff message.LogoffMes
				err = json.Unmarshal([]byte(mes.Data), &mesNotifyLogoff)
				delete(onlineUsers, mesNotifyLogoff.UserId)
				if err != nil {
					fmt.Println("json.Unmarshal in receive logoff mes err=", err)
				} else {
					fmt.Printf("%s is offline\n", mesNotifyLogoff.UserId)
				}
			case message.PrivateMesType:
				var groupMes message.PrivateMes
				err = json.Unmarshal([]byte(mes.Data), &groupMes)
				if err != nil {
					fmt.Println("json.Unmarshal in receive private mes err=", err)
				}
				userid := groupMes.UserId
				fmt.Printf("user %s sends a private mes: \n", userid)
				fmt.Println(groupMes.Content)
			case message.PrivateResMesType:
				fmt.Println(mes.Data)
			default:
				fmt.Println("server sends a fucking mes type")
		}
	}
}