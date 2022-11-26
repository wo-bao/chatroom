package process

import (
	"chatRoom/client/utils"
	"chatRoom/common/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

func (sp *SmsProcess)SendGroupMes(content string) (err error) {
	var mes message.Message
	mes.Type = message.SmsMesType

	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal failed =", err)
		return
	}

	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal failed =", err)
		return
	}

	t := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal failed =", err)
		return
	}
	return
}

func (sp *SmsProcess)SendLogoffMes() (err error) {
	var mes message.Message
	mes.Type = message.LogoffMesType
	var logoffMes message.LogoffMes
	logoffMes.UserId = CurUser.UserId
	data, err := json.Marshal(logoffMes)
	if err != nil {
		fmt.Println("SendLogoffMes json.Marshal failed =", err)
		return
	}

	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendLogoffMes json.Marshal failed =", err)
		return
	}

	t := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("SendLogoffMes json.Marshal failed =", err)
		return
	}
	return
}

func (sp *SmsProcess)SendPrivateMes(userID string, content string) (err error) {
	var mes message.Message
	mes.Type = message.PrivateMesType

	var privateMes message.PrivateMes
	privateMes.Content = content
	privateMes.UserId = CurUser.UserId
	privateMes.UserStatus = CurUser.UserStatus
	privateMes.ToUser = userID

	data, err := json.Marshal(privateMes)
	if err != nil {
		fmt.Println("Send privateMes json.Marshal failed =", err)
		return
	}

	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("Send privateMes json.Marshal failed =", err)
		return
	}

	t := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("Send privateMes json.Marshal failed =", err)
		return
	}
	return
}