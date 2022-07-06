package main

import "github.com/gin-gonic/gin"

func SampleRouter(ctx *AppOption, g *gin.RouterGroup) {
	auth := g.Group("", HandlerWrapper(Auth, ctx))

	sampleR := auth.Group("/samples")
	{
		sampleR.GET("/sample/:id", HandlerWrapper(GetHandler, ctx))
		sampleR.POST("/sample", HandlerWrapper(PostHandler, ctx))
	}
}

func SSTRouter(ctx *AppOption, g *gin.RouterGroup) {
	auth := g.Group("", HandlerWrapper(Auth, ctx))

	sstR := auth.Group("/ssts")
	{
		// create server side token
		sstR.POST("/sst", HandlerWrapper(CreateSSTHandler, ctx))
	}

}
