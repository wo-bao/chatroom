package main

import (
	"chatRoom/client/process"
	"fmt"
)

var Userid string
var Userpsd string
var Username string

func main() {
	var key int
	var loop = true
	for loop {
		fmt.Println("------------欢迎登陆多人聊天系统------------")
		fmt.Println("\t\t 1 登陆聊天室")
		fmt.Println("\t\t 2 注册用户")
		fmt.Println("\t\t 3 退出系统")
		fmt.Println("\t\t 请选择(1-3):")

		fmt.Scanf("%d\n",&key)
		switch key {
			case 1 :
				fmt.Println("登陆聊天室")
				fmt.Println("请输入用户id:")
				fmt.Scanf("%s\n", &Userid)
				fmt.Println("请输入用户密码:")
				fmt.Scanf("%s\n", &Userpsd)
				up := &process.UserProcess{}
				err := up.Login(Userid, Userpsd)
				if err!= nil {
					fmt.Println("log failed, err=", err)
				}
			case 2 :
				fmt.Println("注册用户")
				fmt.Println("请输入注册用户id:")
				fmt.Scanf("%s\n", &Userid)
				fmt.Println("请输入注册用户密码:")
				fmt.Scanf("%s\n", &Userpsd)
				fmt.Println("请输入注册用户名:")
				fmt.Scanf("%s\n", &Username)
				up := &process.UserProcess{}
				err := up.Register(Userid, Userpsd, Username)
				if err!= nil {
					fmt.Println("register failed, err=", err)
				} else {
					fmt.Println("register success")
				}
			case 3 :
				fmt.Println("退出系统")
				loop = false
			default:
				fmt.Println("你的输入有误, 请重新输入")
		}
	}
}