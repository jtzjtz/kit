package mongo_pool

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

const (
	INITIAL_CONNECTION = 5
	MAX_CONNECTION     = 50
	AVAILABLE          = false
	USED               = true
)

var mu sync.RWMutex

type mongodata struct {
	client *mongo.Client
	pos    int
	flag   bool
}

type ClientPool struct {
	uri        string
	clientList [MAX_CONNECTION]mongodata
	size       int
}

func NewClientPool(uri string) (*ClientPool, error) {
	if len(uri) < 1 {
		return nil, errors.New("connect uri is empty")
	}
	cp := &ClientPool{
		uri:        uri,
		clientList: [50]mongodata{},
		size:       0,
	}
	for size := 0; size < INITIAL_CONNECTION || size < MAX_CONNECTION; size++ {
		err := cp.allocateCToPool(size)
		log.Fatal("init - initial create the connect conn failed, size:", size, err)
		return nil, err
	}
	return cp, nil
}

func (cp *ClientPool) Dbconnect() (client *mongo.Client, err error) {
	client, err = mongo.NewClient(options.Client().ApplyURI(cp.uri))
	if err != nil {
		log.Fatal("Dbconnect - connect mongodb failed", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Dbconnect - connect mongodb ctx failed", err)
		return nil, err
	}

	return client, nil
}

func (cp *ClientPool) Dbdisconnect(client *mongo.Client) (err error) {
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal("Dbdisconnect - disconnect mongodb failed", err)
	}
	return err
}

func (cp *ClientPool) allocateCToPool(pos int) (err error) {
	cp.clientList[pos].client, err = cp.Dbconnect()
	if err != nil {
		log.Fatal("allocateCToPool - allocateCToPool failed,position: ", pos, err)
		return err
	}

	cp.clientList[pos].flag = USED
	cp.clientList[pos].pos = pos
	return nil
}

func (cp *ClientPool) getCToPool(pos int) {
	cp.clientList[pos].flag = USED
}

func (cp *ClientPool) putCBackPool(pos int) {
	cp.clientList[pos].flag = AVAILABLE
}

func (cp *ClientPool) GetClient() (mongoclient *mongodata, err error) {
	mu.RLock()
	for i := 1; i < cp.size; i++ {
		if cp.clientList[i].flag == AVAILABLE {
			return &cp.clientList[i], nil
		}
	}
	mu.RUnlock()
	mu.Lock()
	defer mu.Unlock()
	if cp.size < MAX_CONNECTION {
		err = cp.allocateCToPool(cp.size)
		if err != nil {
			log.Fatal("GetClient - DB pooling allocate failed", err)
			return nil, err
		}
		pos := cp.size
		cp.size++
		return &cp.clientList[pos], nil
	} else {
		log.Fatal("GetClient - DB pooling is fulled")
		return nil, errors.New("DB pooling is fulled")
	}
}

func (cp *ClientPool) ReleaseClient(mongoclient *mongodata) {
	mu.Lock()
	cp.putCBackPool(mongoclient.pos)
	mu.Unlock()
}
