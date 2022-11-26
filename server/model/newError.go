package model

import "errors"

var (
	ERROR_USER_NOTEXISTS = errors.New("user is not exists")
	ERROR_USER_REPETITION = errors.New("user already exists")
	ERROR_USER_PSDWRONG = errors.New("password is wrong")
)