package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var Pool *redis.Pool

func initPool(maxIdle, maxActive int, idleTime time.Duration, address string) {
	Pool = &redis.Pool {
		MaxIdle: maxIdle,
		MaxActive: maxActive,
		IdleTimeout: idleTime,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
