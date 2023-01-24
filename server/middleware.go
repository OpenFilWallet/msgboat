package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("panic recover err: %v", err)

				ReturnError(c, NewError(500, fmt.Sprintf("panic recover err: %v", err)))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}

func (b *Boat) RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		count := b.limiterBuckets.TakeAvailable(1)
		if count == 0 {
			ReturnError(c, NewError(500, "Too many requests, please try again later"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func ContextTimeout(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
