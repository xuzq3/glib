package wechat

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/silenceper/wechat/v2/cache"
)

var _ cache.Cache = (*Redis)(nil)

// Redis .redis cache
// @description: 实现cache.Cache的Redis类
type Redis struct {
	ctx  context.Context
	conn redis.UniversalClient
}

// NewRedis 实例化
func NewRedis(ctx context.Context, conn redis.UniversalClient) *Redis {
	return &Redis{ctx: ctx, conn: conn}
}

// Get 获取一个值
func (r *Redis) Get(key string) interface{} {
	result, err := r.conn.Do(r.ctx, "GET", key).Result()
	if err != nil {
		return nil
	}
	return result
}

// Set 设置一个值
func (r *Redis) Set(key string, val interface{}, timeout time.Duration) error {
	return r.conn.SetEX(r.ctx, key, val, timeout).Err()
}

// IsExist 判断key是否存在
func (r *Redis) IsExist(key string) bool {
	result, _ := r.conn.Exists(r.ctx, key).Result()

	return result > 0
}

// Delete 删除
func (r *Redis) Delete(key string) error {
	return r.conn.Del(r.ctx, key).Err()
}
