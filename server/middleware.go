package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("panic recover err: %v", err)

				ReturnError(c, NewError(500, fmt.Sprintf("panic recover err: %v", err)))
				c.Abort()
			}
		}()
		c.Next()
	}
}
