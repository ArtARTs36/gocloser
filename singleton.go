package closer

import (
	"context"
	"time"
)

var closer = NewCloser()

func Close(ctx context.Context) {
	closer.Close(ctx)
}

func Bind(name string, fn func() error) {
	closer.Bind(name, fn)
}

func BindWithCtx(name string, fn func(ctx context.Context) error) {
	closer.BindWithCtx(name, fn)
}

func Subscribe(closeTimeout time.Duration) {
	closer.Subscribe(closeTimeout)
}
