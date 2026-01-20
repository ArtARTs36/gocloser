package closer

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Closer struct {
	mux sync.Mutex

	closers []*resource
}

func NewCloser() *Closer {
	return &Closer{closers: make([]*resource, 0)}
}

type resource struct {
	name  string
	close func(ctx context.Context) error
}

func (c *Closer) Close(ctx context.Context) {
	c.mux.Lock()
	defer c.mux.Unlock()

	for i := len(c.closers) - 1; i >= 0; i-- {
		if err := c.closers[i].close(ctx); err != nil {
			slog.
				With(slog.String("resource", c.closers[i].name)).
				With(slog.Any("err", err)).
				ErrorContext(ctx, "[gocloser] failed to close")
		}
	}
}

func (c *Closer) Bind(name string, fn func() error) {
	c.BindWithCtx(name, func(_ context.Context) error {
		return fn()
	})
}

func (c *Closer) BindWithCtx(name string, fn func(ctx context.Context) error) {
	c.mux.Lock()
	c.closers = append(c.closers, &resource{
		name: name,
		close: func(ctx context.Context) error {
			return fn(ctx)
		},
	})
	c.mux.Unlock()
}

func (c *Closer) Subscribe(closeTimeout time.Duration) {
	channel := make(chan os.Signal, 1)

	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		sig := <-channel
		slog.Info("[gocloser] received signal", slog.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
		defer cancel()
		c.Close(ctx)

		done <- true
	}()
}
