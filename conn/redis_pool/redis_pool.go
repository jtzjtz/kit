package redis_pool

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func NewRedisPool(o *Options) (*redis.Pool, error) {

	if err := o.validate(); err != nil {
		return nil, err
	}

	pool := &redis.Pool{
		MaxActive:   o.MaxCap,
		MaxIdle:     o.InitCap,
		IdleTimeout: o.IdleTimeout,
		Wait:        o.IsWait,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp",
				o.Host,
				redis.DialDatabase(o.Database),
				redis.DialPassword(o.PassWord),
				redis.DialConnectTimeout(o.DialConnectTimeout),
				redis.DialReadTimeout(o.DialReadTimeout),
				redis.DialWriteTimeout(o.DialWriteTimeout),
			)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return conn, nil
		},
	}
	return pool, nil
}
