package lock

import (
	"context"
	"time"
)

// Locker ...
type Locker interface {
	// TryLock no block
	TryLock(ctx context.Context, key string, expiration time.Duration) (bool, UnLocker, error)
	// Lock block until successful or timed out
	Lock(ctx context.Context, key string, expiration time.Duration) (UnLocker, error)
}

// UnLocker ...
type UnLocker interface {
	Unlock() error
}
