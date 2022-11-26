package model

import (
	"chatRoom/common/message"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var (
	MyUserDao *UserDao
)

type UserDao struct {
	Pool *redis.Pool
}

func (ud *UserDao)GetUserById(conn redis.Conn, id string) (user User, err error) {
	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
			return
		}
	}

	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

func (ud *UserDao)FindRepetition(conn redis.Conn, user *message.RegisterMes) (userId string, err error) {
	_, err = redis.String(conn.Do("hget", "users", user.UserId))
	userId = user.UserId
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		return
	} else {
		err = ERROR_USER_REPETITION
		return
	}
}

func (ud *UserDao)Login(userId string, userPsd string) (user User, err error) {
	conn := ud.Pool.Get()
	defer conn.Close()
	user, err = ud.GetUserById(conn, userId)
	if err != nil {
		return
	}
	if user.UserPsd == userPsd {
		fmt.Println("Password id correct, Login success!")
		return
	} else {
		fmt.Println("repetition ahahha.....")
		err = ERROR_USER_PSDWRONG
		return
	}
}

func (ud *UserDao)Register(user *message.RegisterMes) (userId string, err error) {
	conn := ud.Pool.Get()
	defer conn.Close()
	userId, err = ud.FindRepetition(conn, user)
	if err != nil {
		return
	}
	userInfoM, err := json.Marshal(*user)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	userInfo := string(userInfoM)
	fmt.Println(userId)

	_, err = redis.String(conn.Do("hset", "users", userId, userInfo))
	return
}

func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		Pool: pool,
	}
	return userDao
}