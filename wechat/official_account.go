package wechat

import (
	"context"
	"github.com/xuzq3/glib/lock"
	"math/rand"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/credential"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
)

type OfficialAccountConfig struct {
	AppID          string `json:"app_id"`
	AppSecret      string `json:"app_secret"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
}

type OfficialAccount struct {
	officialaccount.OfficialAccount
	accessTokenHandle *DistributedAccessTokenHandle
}

func NewOfficialAccount(config *OfficialAccountConfig, cache Cache, distributedLock lock.Locker) *OfficialAccount {
	// 初始化微信
	wc := wechat.NewWechat()
	wc.SetCache(cache)

	// 初始化公众号
	officialAccount := wc.GetOfficialAccount(&offConfig.Config{
		AppID:          config.AppID,
		AppSecret:      config.AppSecret,
		Token:          config.Token,
		EncodingAESKey: config.EncodingAESKey,
	})
	accessTokenHandle := NewDistributedAccessTokenHandle(
		config.AppID,
		config.AppSecret,
		credential.CacheKeyOfficialAccountPrefix,
		cache,
		distributedLock,
	)
	officialAccount.SetAccessTokenHandle(accessTokenHandle)

	oa := &OfficialAccount{
		OfficialAccount:   *officialAccount,
		accessTokenHandle: accessTokenHandle,
	}
	oa.StartTokenChecking()
	return oa
}

func NewOfficialAccountByRedis(config *OfficialAccountConfig, redisClient redis.UniversalClient) *OfficialAccount {
	cache := NewRedis(context.Background(), redisClient)
	distributedLock := lock.NewRedisLocker(redisClient)
	return NewOfficialAccount(config, cache, distributedLock)
}

func NewOfficialAccountByMemory(config *OfficialAccountConfig) *OfficialAccount {
	cache := cache.NewMemory()
	return NewOfficialAccount(config, cache, nil)
}

func (oa *OfficialAccount) StartTokenChecking() {
	go func() {
		ticker := time.NewTicker(time.Second*30 + time.Millisecond*time.Duration(rand.Intn(1000))) // 固定时间加随机时间，防止多进程同时执行
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				oa.checkToken()
			}
		}
	}()
}

func (oa *OfficialAccount) checkToken() {
	res, err := oa.OfficialAccount.GetBasic().GetAPIDomainIP()
	if err != nil && strings.Contains(err.Error(), "errcode=40001") {
		// token失效，需要刷新token
		_ = oa.accessTokenHandle.RefreshAccessToken()
	}
	_ = res
}
