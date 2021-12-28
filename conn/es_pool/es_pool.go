package es_pool

import (
	"github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
)

type clientPool struct {
	client *elastic.Client
	err    error
}

func NewPool(options *Options) *clientPool {
	pool := new(clientPool)
	conf := &config.Config{
		URL:         options.URL,
		Index:       options.Index,
		Username:    options.Username,
		Password:    options.Password,
		Shards:      options.Shards,
		Replicas:    options.Replicas,
		Sniff:       options.Sniff,
		Healthcheck: options.Healthcheck,
		Infolog:     options.Infolog,
		Errorlog:    options.Errorlog,
		Tracelog:    options.Tracelog,
	}
	client, err := elastic.NewClientFromConfig(conf)
	pool.client = client
	if err != nil {
		panic(err)
	}
	return pool
}
func (cl *clientPool) GetClient() *elastic.Client {
	return cl.client
}
func (cl *clientPool) Close() {
	cl.client.Stop()
}
