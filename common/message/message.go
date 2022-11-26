package message

const (
	LoginMesType = "LoginMes"
	LoginResMesType = "LoginResMes"
	RegisterMesType = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatus"
	SmsMesType = "SmsMes"
	LogoffMesType = "LogoffMes"
	PrivateMesType = "PrivateMes"
	PrivateResMesType = "PrivateResMes"
)

const (
	UserOnline = iota
	UserOffline
	UserBusy
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type LoginMes struct {
	UserId string `json:"userid"`
	UserPsd string `json:"userpsd"`
	UserName string `json:"username"`
}

type LoginResMes struct {
	Code int `json:"code"`//500表示无此用户 200表示登陆成功 300表示密码错误
	Users []string `json:"users"`
	Error string `json:"error"`//返回错误信息
}

type Register struct {
	User User `json:"user"`
}

type RegisterMes struct {
	UserId string `json:"userid"`
	UserPsd string `json:"userpsd"`
	UserName string `json:"username"`
}

type RegisterResMes struct {
	Code int `json:"code"`//400表示次Id已被占用, 100表示注册成功
	Error string `json:"error"`//
}

type NotifyUserStatusMes struct {
	UserId string `json:"userid"`
	Status int `json:"status"`
}

type SmsMes struct {
	Content string `json:"content"`
	User
}

type PrivateMes struct {
	ToUser string `json:"touser"`
	Content string `json:"content"`
	User
}

type LogoffMes struct {
	UserId string `json:"userid"`
}