package lock

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var _ Locker = (*redisLocker)(nil)

type redisLocker struct {
	redis redis.Cmdable
}

func NewRedisLocker(client redis.Cmdable) *redisLocker {
	return &redisLocker{
		redis: client,
	}
}

func (r *redisLocker) TryLock(ctx context.Context, key string, expiration time.Duration) (bool, UnLocker, error) {
	id := uuid.New().String()
	ok, err := r.redis.SetNX(ctx, key, id, expiration).Result()
	if err != nil {
		return false, nil, err
	}
	if !ok {
		return false, nil, nil
	}
	lock := &unlocker{
		client: r,
		key:    key,
		value:  id,
	}
	return true, lock, nil
}

func (r *redisLocker) Lock(ctx context.Context, key string, expiration time.Duration) (UnLocker, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		ok, lock, err := r.TryLock(ctx, key, expiration)
		if err != nil {
			return nil, err
		}
		if ok {
			return lock, nil
		}
		time.Sleep(time.Millisecond * 50)
	}
}

const (
	deleteLua = `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`
)

var _ UnLocker = (*unlocker)(nil)

type unlocker struct {
	client *redisLocker
	key    string
	value  string
}

func (r *unlocker) Unlock() error {
	return r.UnlockWithCtx(context.Background())
}

func (r *unlocker) UnlockWithCtx(ctx context.Context) error {
	res, err := r.client.redis.Eval(ctx, deleteLua, []string{r.key}, r.value).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return err
	}
	_ = res
	return nil
}
