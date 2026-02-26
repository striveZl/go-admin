package rbac

import (
	"context"
	"sync/atomic"
	"time"
)

type Casbinx struct {
	enforcer *atomic.Value `wire:"-"`
	ticker   *time.Ticker  `wire:"-"`
}

func (a *Casbinx) Release(ctx context.Context) error {
	if a.ticker != nil {
		a.ticker.Stop()
	}
	return nil
}
