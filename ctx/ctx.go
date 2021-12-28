package ctx

import (
	"context"
	"time"
)

func CreateContext() (context.Context, context.CancelFunc) {
	return CreateContextWithTime(time.Second * 30)
}

func CreateContextWithTime(timeOut time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeOut)
}
