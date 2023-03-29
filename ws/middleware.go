package ws

import (
	"context"
	"time"
)

func Timeout() MessageHandler {
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

//func Recovery() Handler {
//	return func(c *MessageContext) {
//		defer func() {
//			if r := recover(); r != nil {
//				err := fmt.Errorf("panic: %v", r)
//				c.AbortWithError(err)
//			}
//		}()
//
//		c.Next()
//	}
//}
