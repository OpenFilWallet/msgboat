package server

import "github.com/gin-gonic/gin"

func (b *Boat) NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(Recovery())

	r.GET("/status", b.Status)
	r.POST("/send", b.Send)

	return r
}
