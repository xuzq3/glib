package multicast

import (
	"context"
	"fmt"
	"time"

	"github.com/xuzq3/glib/logx"
)

func Log() Handler {
	return func(c *MessageContext) {
		stime := time.Now()
		logx.Info("multicast serve from:%s cmd:%s seq:%s data:%s",
			c.Message.Src.IP.String(), c.Body.Cmd, c.Body.Seqno, string(c.Body.Data))

		c.Next()

		if c.Error != nil {
			logx.Error("multicast serve failed cmd:%s seq:%s latency:%v err:%s",
				c.Body.Cmd, c.Body.Seqno, time.Since(stime), c.Error.Error())
		} else {
			logx.Debug("multicast serve success cmd:%s seq:%s latency:%v",
				c.Body.Cmd, c.Body.Seqno, time.Since(stime))
		}
	}
}

func Timeout() Handler {
	return func(c *MessageContext) {
		ctx, cancel := context.WithTimeout(c.Context(), time.Minute*5)
		c.WithContext(ctx)

		c.Next()

		if ctx.Err() == context.DeadlineExceeded {
			c.AbortWithError(ctx.Err())
		}
		cancel()
	}
}

func Recovery() Handler {
	return func(c *MessageContext) {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic: %v", r)
				c.AbortWithError(err)
			}
		}()

		c.Next()
	}
}
