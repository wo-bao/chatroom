package process

import (
	"chatRoom/client/model"
	"chatRoom/common/message"
	"fmt"
)

var onlineUsers = make(map[string]*message.User, 10)
var CurUser model.CurUser

func updateUserStatus(mes *message.NotifyUserStatusMes) {
	user, ok := onlineUsers[mes.UserId]
	if ok {
		user.UserStatus = mes.Status
	} else {
		user = &message.User{
			UserId: mes.UserId,
		}
		user.UserStatus = mes.Status
	}
	onlineUsers[mes.UserId] = user
	outputOnlineUsers()
}

func outputOnlineUsers() {
	for id, user := range onlineUsers {
		if user.UserStatus == 0 {
			fmt.Printf("%s is online\n", id)
		} else if user.UserStatus == 2 {
			fmt.Printf("%s is busy\n", id)
		} else {
			fmt.Printf("%s is ????\n", id)
		}
	}
}