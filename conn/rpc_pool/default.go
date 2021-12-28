package rpc_pool

import (
	"sync"
	"time"

	"github.com/jtzjtz/kit/conn/grpc_pool"
)

var pool *RpcPool

func Default() *RpcPool {
	return pool
}

var once sync.Once

// InitP 初始化，失败时panic
func InitP(host string, maxCap int) {
	once.Do(func() {
		p, err := RpcPool{}.Connect(&grpc_pool.Options{
			InitTargets: []string{host},
			InitCap:     5,
			MaxCap:      maxCap,
			DialTimeout: time.Second * 30,
			IdleTimeout: time.Second * 60 * 60,
		})
		if err != nil {
			panic(err)
		}

		pool = &p
	})
}
