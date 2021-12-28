package es_pool

type Options struct {
	URL         string
	Index       string
	Username    string
	Password    string
	Shards      int
	Replicas    int
	Sniff       *bool
	Healthcheck *bool
	Infolog     string
	Errorlog    string
	Tracelog    string
}
