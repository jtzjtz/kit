package mysql_pool

import (
	"errors"
)

var (
	errInvalid = errors.New("invalid config")
)

//Options conn options
type Options struct {
	// init connection
	InitCap int
	// max connections
	MaxCap int

	IsDebug bool

	User     string
	Pass     string
	Host     string
	Port     string
	DataBase string
}

// NewOptions returns a new newOptions instance with sane defaults.
func NewOptions() *Options {
	o := &Options{}
	o.InitCap = 5
	o.MaxCap = 100
	o.IsDebug = true
	return o
}

// validate checks a Config instance.
func (o *Options) validate() error {
	if o.InitCap <= 0 ||
		o.MaxCap <= 0 ||
		o.InitCap > o.MaxCap ||
		o.User == "" ||
		o.Pass == "" ||
		o.Host == "" ||
		o.Port == "" ||
		o.DataBase == "" {
		return errInvalid
	}
	return nil
}
