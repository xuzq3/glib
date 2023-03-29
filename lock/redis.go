package lock

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Locker interface {
	TryLock(ctx context.Context, key string, expiration time.Duration) (bool, UnLocker, error)
	Lock(ctx context.Context, key string, expiration time.Duration) (UnLocker, error)
}

type UnLocker interface {
	Unlock(ctx context.Context) error
}

var _ Locker = (*RedisLocker)(nil)

type RedisLocker struct {
	redis redis.Cmdable
}

func NewRedisLocker(client redis.Cmdable) *RedisLocker {
	return &RedisLocker{
		redis: client,
	}
}

func (r *RedisLocker) TryLock(ctx context.Context, key string, expiration time.Duration) (bool, UnLocker, error) {
	id := uuid.New().String()
	ok, err := r.redis.SetNX(ctx, key, id, expiration).Result()
	if err != nil {
		return false, nil, err
	}
	if !ok {
		return false, nil, nil
	}
	lock := &Lock{
		client: r,
		key:    key,
		value:  id,
	}
	return true, lock, nil
}

func (r *RedisLocker) Lock(ctx context.Context, key string, expiration time.Duration) (UnLocker, error) {
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

type Lock struct {
	client *RedisLocker
	key    string
	value  string
}

func (r *Lock) Unlock(ctx context.Context) error {
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
