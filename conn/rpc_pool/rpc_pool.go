package rpc_pool

import (
	"errors"
	"github.com/jtzjtz/kit/conn/grpc_pool"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

//rpc 连接池
type RpcPool struct {
	rpcPool    *grpc_pool.GRPCPool
	rpcOptions *grpc_pool.Options
	mu         sync.Mutex
	getMu      sync.Mutex
}

//取rpc连接 （支持连接池重连）
func (current RpcPool) Get() (*grpc.ClientConn, error) {
	current.getMu.Lock()
	defer current.getMu.Unlock()
	if current.rpcPool == nil {
		err := current.reConnect()
		if err != nil {
			return nil, err
		}
	}
	con, errRpc := current.rpcPool.Get()
	if errRpc != nil { //之前连接池有值，现在连接不上 重连
		errCon := current.reConnect()
		if errCon != nil {
			return nil, errors.New("之前连接池有值，现在连接不上 重连失败 原因：" + errCon.Error())
		} else if current.rpcPool != nil { //重连成功
			return current.rpcPool.Get()
		} else {
			return nil, errors.New("之前连接池有值，现在连接不上 重连后 连接池为空")
		}

	} else {
		return con, nil
	}

}

//释放rpc连接  建议用此方法 名字好，容易理解
func (current RpcPool) Dispose(rpcConnection *grpc.ClientConn) error {
	return current.Put(rpcConnection)
}

//释放rpc连接
func (current RpcPool) Put(rpcConnection *grpc.ClientConn) error {
	if current.rpcPool != nil {
		return current.rpcPool.Put(rpcConnection)
	}
	return errors.New("不存在的连接池不能放回连接")
}

//重新连接
func (current RpcPool) reConnect() error {
	if current.rpcOptions == nil {
		return errors.New("rpc 初始值为空，请重启应用赋值")

	}
	_, err := current.Connect(current.rpcOptions)
	time.Sleep(time.Second * 1) //初始化连接池为异步，等待1s
	if current.rpcPool == nil { //如果还为空 返回错误
		if err != nil {
			return errors.New(" GRPC 连接重连失败，请稍后重试 原因：" + err.Error())
		}
		return errors.New(" GRPC 连接重连失败，请稍后重试")
	}
	return err

}

//初始化连接池
func (current RpcPool) Connect(rpcOptions *grpc_pool.Options) (RpcPool, error) {
	var err error
	current.rpcOptions = rpcOptions
	if current.rpcPool == nil {
		current.mu.Lock()
		if current.rpcPool == nil {
			current.rpcPool, err = grpc_pool.NewGRPCPool(rpcOptions, grpc.WithInsecure(), grpc.WithBlock())
		}
		current.mu.Unlock()
	}

	if err != nil {
		if len(rpcOptions.InitTargets) > 0 {
			log.Printf("getGrpcPool err:%#v url:%#v\n", err, rpcOptions.InitTargets[0])
		} else {
			log.Printf("getGrpcPool err:%#v url: null", err)
		}
	}
	return current, err

}
