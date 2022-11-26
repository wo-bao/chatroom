package process

import (
	"chatRoom/common/message"
	"chatRoom/server/model"
	"chatRoom/server/utils"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	UserId string
}

func (u *UserProcess)NotifyOtherOnlineUser(userId string) {
	for id, up := range userMgr.OnlineUsers {
		if id == userId {
			continue
		}
		err := up.NotifyMeOnline(userId)
		if err != nil {
			fmt.Println("Notify err=", err)
		}
	}
}

func (u *UserProcess)NotifyMeOnline(userId string) (err error) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	data, err := json.Marshal(notifyUserStatusMes)
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

	t := &utils.Transfer{
		Conn: u.Conn,
	}
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
	return
}

func (u *UserProcess)ServerProcessLogin(mes *message.Message) (err error) {
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.LoginResMesType
	var loginResMes message.LoginResMes

	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPsd)
	if err != nil {
		if err == model.ERROR_USER_PSDWRONG{
			loginResMes.Code = 403
			loginResMes.Error = "Password is wrong"
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "Inner error of server..."
		}
	} else {
		loginResMes.Code = 200
		u.UserId = loginMes.UserId
		fmt.Print(user, "Login success")
		userMgr.AddOnlineUser(u)
		u.NotifyOtherOnlineUser(loginMes.UserId)
		for id, _ := range userMgr.OnlineUsers {
			loginResMes.Users = append(loginResMes.Users, id)
		}
	}

	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	resMes.Data = string(data)

	data ,err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	t := &utils.Transfer{
		Conn: u.Conn,
	}
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("writePkg err=",err)
		return
	}
	fmt.Println("writePkg success")

	connToRedis := model.MyUserDao.Pool.Get()
	defer connToRedis.Close()
	for {
		SavedMes, errGetSavedMes := redis.String(connToRedis.Do("lpop", u.UserId))
		if errGetSavedMes == nil {
			errSendSavedMes := t.WritePkg([]byte(SavedMes))
			if errSendSavedMes != nil {
				fmt.Println("sendSavedMes err=",errSendSavedMes)
			}
		} else {
			break
		}
	}
	return
}

func (u *UserProcess)ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	user, err := model.MyUserDao.Register(&registerMes)
	if err != nil {
		if err == model.ERROR_USER_REPETITION{
			registerResMes.Code = 400
			registerResMes.Error = "User is Repetition"
		} else {
			registerResMes.Code = 505
			registerResMes.Error = "Inner error of server..."
		}
	} else {
		registerResMes.Code = 100
		fmt.Print(user, "Register success")
	}
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	resMes.Data = string(data)

	data ,err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	t := &utils.Transfer{
		Conn: u.Conn,
	}
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("writePkg err=",err)
		return
	}
	fmt.Println("writePkg success")

	return
}

func (u *UserProcess)ServerProcessGroupMes(mes *message.Message) (err error) {
	data, err := json.Marshal(*mes)
	if err != nil {
		fmt.Println("json.marshal in Distribute GroupMes err=", err)
		return
	}
	users := userMgr.OnlineUsers
	t := utils.Transfer{}
	for _, v := range users {
		if v.Conn == u.Conn {
			continue
		}
		t.Conn = v.Conn
		err = t.WritePkg(data)
		if err != nil {
			fmt.Println("writePkg in Distribute GroupMes err=",err)
			return
		}
	}

	fmt.Println("writePkg in Distribute GroupMes success")
	return
}

func (u *UserProcess)ServerProcessLogoff(mes *message.Message) (err error) {
	var logoffMes message.LogoffMes
	err = json.Unmarshal([]byte(mes.Data), &logoffMes)
	if err != nil {
		fmt.Println("json.Unmarshal in process LogoffMes err=", err)
		return
	}

	t := utils.Transfer{}
	for _, v := range userMgr.OnlineUsers {
		if v.Conn == u.Conn {
			delete(userMgr.OnlineUsers, logoffMes.UserId)
			u.Conn.Close()
			continue
		}
		t.Conn = v.Conn
		logoffNotifyMes, err := ProduceLogoffNotifyMes(mes)
		if err != nil {
			fmt.Println("writePkg in Produce LogoffMes err=",err)
			return err
		}
		err = t.WritePkg(logoffNotifyMes)
		if err != nil {
			fmt.Println("writePkg in Distribute LogoffMes err=",err)
			return err
		}
	}

	fmt.Println("writePkg in Distribute GroupMes success")
	return
}

func ProduceLogoffNotifyMes(logoffMes *message.Message) (logoffNotifyMes []byte,err error) {
	logoffNotifyMes, err = json.Marshal(*logoffMes)
	return
}

func (u *UserProcess)SendLogoffForciblyMes() (err error) {
	var mes message.Message
	mes.Type = message.LogoffMesType
	var logoffMes message.LogoffMes

	for _, v := range userMgr.OnlineUsers{
		if v.Conn == u.Conn {
			logoffMes.UserId = v.UserId
			delete(userMgr.OnlineUsers, logoffMes.UserId)
			u.Conn.Close()
			continue
		}
	}

	t := utils.Transfer{}
	for _, v := range userMgr.OnlineUsers {
		t.Conn = v.Conn
		data, err := json.Marshal(logoffMes)
		if err != nil {
			fmt.Println("SendLogoffForciblyMes json.Marshal failed =", err)
			return err
		}

		mes.Data = string(data)
		data, err = json.Marshal(mes)
		if err != nil {
			fmt.Println("SendLogoffForciblyMesMes json.Marshal failed =", err)
			return err
		}
		logoffNotifyMes, err := ProduceLogoffNotifyMes(&mes)
		if err != nil {
			fmt.Println("writePkg in Produce Forcibly LogoffMes err=",err)
			return err
		}
		err = t.WritePkg(logoffNotifyMes)
		if err != nil {
			fmt.Println("writePkg in Distribute Forcibly LogoffMes err=",err)
			return err
		}
	}
	return
}

func (u *UserProcess)ServerProcessPrivateMes(mes *message.Message) (err error) {
	var privateMes message.PrivateMes
	err = json.Unmarshal([]byte(mes.Data), &privateMes)
	if err != nil {
		fmt.Println("json.Unmarshal in process LogoffMes err=", err)
		return
	}
	toUserId := privateMes.ToUser
	for i, v := range userMgr.OnlineUsers {
		if i == toUserId {
			err = sendPrMesToUser(v.Conn, mes)
			if err != nil {
				fmt.Println("Transfer in sendPrMes toUser err=", err)
				return
			}
			err = u.ToOnlineUserSucPrivateResMes(toUserId)
			if err != nil {
				fmt.Println("sendPrResMes err=", err)
				return
			}
			return
		}
	}
	connToRedis := model.MyUserDao.Pool.Get()
	defer connToRedis.Close()
	_, errFind := model.MyUserDao.GetUserById(connToRedis, toUserId)
	if errFind != nil {
		errTo := u.ToUnknownUserResMes(toUserId)
		if errTo != nil {
			fmt.Println("sendUnknownUserResMes err=", err)
			return errTo
		}
		return errFind
	}
	err = saveOffLineMesTemp(mes, toUserId, connToRedis)
	if err != nil {
		fmt.Println("saveOfflineMes err=", err)
		return
	}
	err = u.ToOffLineUserSucPrivateResMes(toUserId)
	if err != nil {
		fmt.Println("sendPrResMes err=", err)
		return
	}
	return
}

func sendPrMesToUser (connToUser net.Conn, mes *message.Message) (err error) {
	t := utils.Transfer{}
	t.Conn = connToUser
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Unmarshal in sendPrMes toUser err=", err)
		return
	}
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("Transfer in sendPrMes toUser err=", err)
		return
	}
	return
}

func (u *UserProcess)ToOnlineUserSucPrivateResMes(toUser string) (err error) {
	var mes message.Message
	mes.Type = message.PrivateResMesType
	mes.Data = "Private mes to " + toUser + " success"
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Unmarshal in sendPrMes toUser err=", err)
		return
	}
	t := utils.Transfer{}
	t.Conn = u.Conn
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("sendPrResMes err=", err)
		return
	}
	return
}

func (u *UserProcess)ToOffLineUserSucPrivateResMes(toUser string) (err error) {
	var mes message.Message
	mes.Type = message.PrivateResMesType
	mes.Data = "User " + toUser + " is not online, mes is saved in server until " + toUser + " is online"
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Unmarshal in sendPrMes toUser err=", err)
		return
	}
	t := utils.Transfer{}
	t.Conn = u.Conn
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("sendPrResMes err=", err)
		return
	}
	return
}

func (u *UserProcess)ServerProcessErrInPri(errInPri error) (err error) {
	var mes message.Message
	mes.Type = message.PrivateResMesType
	mes.Data = "Server inner error = " + errInPri.Error()
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Unmarshal in sendPrMes toUser err=", err)
		return
	}
	t := utils.Transfer{}
	t.Conn = u.Conn
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("sendPrResMes err=", err)
		return
	}
	return
}

func (u *UserProcess)ToUnknownUserResMes(toUser string) (err error) {
	var mes message.Message
	mes.Type = message.PrivateResMesType
	mes.Data = "user " + toUser + " is not existing"
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Unmarshal in sendUnknownUserResMes err=", err)
		return
	}
	t := utils.Transfer{}
	t.Conn = u.Conn
	err = t.WritePkg(data)
	if err != nil {
		fmt.Println("sendUnknownUserResMes err=", err)
		return
	}
	return
}

func saveOffLineMesTemp(mes *message.Message, toUser string, conn redis.Conn) (err error) {
	 data, err := json.Marshal(*mes)
	 if err != nil {
	 	fmt.Println("json.Marshal in saveOffLineMes err=", err)
	 	return
	 }
	 _, err = conn.Do("rpush", toUser, string(data))
	 if err != nil {
	 	fmt.Println("saveOffLineMes to redis err=", err)
	 	return
	 }
	 return
}