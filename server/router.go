package server

import "github.com/gin-gonic/gin"

func (b *Boat) NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(Recovery())

	r.POST("/push", b.Send)

	return r
}
