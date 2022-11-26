package process

import "errors"

var (
	userMgr *UserMgr
)

type UserMgr struct {
	OnlineUsers map[string]*UserProcess
}

func init() {
	userMgr = &UserMgr{
		OnlineUsers: make(map[string]*UserProcess, 1024),
	}
}

func (um *UserMgr)AddOnlineUser(up *UserProcess) {
	um.OnlineUsers[up.UserId] = up
}

func (um *UserMgr)DelOnlineUser(userId string) {
	delete(um.OnlineUsers, userId)
}

func (um *UserMgr)GetAllOnlineUser() (usermap map[string]*UserProcess) {
	usermap = um.OnlineUsers
	return usermap
}

func (um *UserMgr)GetOnlineUserById(userId string) (up *UserProcess, err error) {
	up, ok := um.OnlineUsers[userId]
	if !ok {
		err = errors.New("user" + userId + "is not online")
		return
	}
	return
}
