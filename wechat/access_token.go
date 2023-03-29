package wechat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/silenceper/wechat/v2/credential"
	"github.com/xuzq3/glib/lock"
)

const (
	// AccessTokenURL 获取access_token的接口
	accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

var _ credential.AccessTokenContextHandle = (*DistributedAccessTokenHandle)(nil)

// DistributedAccessTokenHandle AccessToken 获取
// @Description 重写credential.DefaultAccessToken，实现分布式锁，支持多实例部署
type DistributedAccessTokenHandle struct {
	appID           string
	appSecret       string
	cacheKeyPrefix  string
	cache           Cache
	localLock       *sync.Mutex
	distributedLock lock.Locker
}

// NewDistributedAccessTokenHandle new DistributedAccessTokenHandle
func NewDistributedAccessTokenHandle(
	appID, appSecret, cacheKeyPrefix string, cache Cache, distributedLock lock.Locker,
) *DistributedAccessTokenHandle {
	return &DistributedAccessTokenHandle{
		appID:           appID,
		appSecret:       appSecret,
		cache:           cache,
		cacheKeyPrefix:  cacheKeyPrefix,
		localLock:       new(sync.Mutex),
		distributedLock: distributedLock,
	}
}

// GetAccessToken 获取access_token,先从cache中获取，没有则从服务端获取
func (ak *DistributedAccessTokenHandle) GetAccessToken() (accessToken string, err error) {
	return ak.GetAccessTokenContext(context.Background())
}

func (ak *DistributedAccessTokenHandle) accessTokenCacheKey() string {
	return fmt.Sprintf("%s_access_token:%s", ak.cacheKeyPrefix, ak.appID)
}

func (ak *DistributedAccessTokenHandle) getLockKey() string {
	return fmt.Sprintf("%s_access_token_lock:%s", ak.cacheKeyPrefix, ak.appID)
}

// GetAccessTokenContext 获取access_token,先从cache中获取，没有则从服务端获取
func (ak *DistributedAccessTokenHandle) GetAccessTokenContext(ctx context.Context) (accessToken string, err error) {
	// 先从cache中取
	accessTokenCacheKey := ak.accessTokenCacheKey()
	if val := ak.cache.Get(accessTokenCacheKey); val != nil {
		return val.(string), nil
	}

	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从微信服务器上获取到不同token
	ak.localLock.Lock()
	defer ak.localLock.Unlock()

	// 加上分布式锁，防止多实例同时获取token
	if ak.distributedLock != nil {
		lockKey := ak.getLockKey()
		lock, err := ak.distributedLock.Lock(ctx, lockKey, time.Second*5)
		if err != nil {
			return "", err
		}
		defer func() {
			_ = lock.Unlock(ctx)
		}()
	}

	// 双检，防止重复从微信服务器获取
	if val := ak.cache.Get(accessTokenCacheKey); val != nil {
		return val.(string), nil
	}

	// cache失效，从微信服务器获取
	accessToken, err = ak.doGetAccessTokenContext(ctx)
	return
}

// RefreshAccessToken 刷新access_token
func (ak *DistributedAccessTokenHandle) RefreshAccessToken() (err error) {
	ctx := context.Background()

	ak.localLock.Lock()
	defer ak.localLock.Unlock()

	// 加上分布式锁，防止多实例同时获取token
	if ak.distributedLock != nil {
		lockKey := ak.getLockKey()
		ok, lock, err := ak.distributedLock.TryLock(ctx, lockKey, time.Second*5)
		if err != nil {
			return err
		} else if !ok {
			return nil
		}
		defer func() {
			_ = lock.Unlock(ctx)
		}()
	}

	_, err = ak.doGetAccessTokenContext(ctx)
	return
}

func (ak *DistributedAccessTokenHandle) doGetAccessTokenContext(ctx context.Context) (accessToken string, err error) {
	accessTokenCacheKey := ak.accessTokenCacheKey()

	var resAccessToken credential.ResAccessToken
	resAccessToken, err = credential.GetTokenFromServerContext(ctx, fmt.Sprintf(accessTokenURL, ak.appID, ak.appSecret))
	if err != nil {
		return
	}

	expires := resAccessToken.ExpiresIn - 1500
	err = ak.cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)
	if err != nil {
		return
	}
	accessToken = resAccessToken.AccessToken
	return
}
