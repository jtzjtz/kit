package redis_pool

import (
	"errors"
	"time"
)

var (
	errInvalid = errors.New("invalid config")
)

//redis连接配置
type Options struct {
	Host               string
	PassWord           string
	Database           int
	InitCap            int           // 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
	MaxCap             int           // 最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）
	IsWait             bool          // 当超过最大连接数 是报错还是等待， true 等待 false 报错
	IdleTimeout        time.Duration //空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用
	DialConnectTimeout time.Duration
	DialReadTimeout    time.Duration
	DialWriteTimeout   time.Duration
}

//实例化一个默认配置
func NewOptions() *Options {
	o := &Options{}
	o.InitCap = 10
	o.MaxCap = 100
	o.IsWait = true
	o.Database = 0
	o.IdleTimeout = 5 * time.Second
	o.DialConnectTimeout = 5 * time.Second
	o.DialReadTimeout = 5 * time.Second
	o.DialWriteTimeout = 5 * time.Second
	return o
}

//配置验证
func (o *Options) validate() error {
	if o.InitCap <= 0 ||
		o.MaxCap <= 0 ||
		o.InitCap > o.MaxCap ||
		o.Host == "" ||
		o.IdleTimeout == 0 ||
		o.DialConnectTimeout == 0 ||
		o.DialReadTimeout == 0 ||
		o.DialWriteTimeout == 0 {
		return errInvalid
	}
	return nil
}
