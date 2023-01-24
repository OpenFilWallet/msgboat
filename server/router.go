package server

import (
	"github.com/gin-gonic/gin"
	"time"
)

func (b *Boat) NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(Recovery())
	r.Use(ContextTimeout(10 * time.Second))
	r.Use(b.RateLimiter())

	r.GET("/status", b.Status)
	r.POST("/send", b.Send)

	return r
}
