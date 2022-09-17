package main

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/crud-server-side-token/sstapp"
)

func SampleRouter(ctx *AppOption, g *gin.RouterGroup) {
	auth := g.Group("", sstapp.GinAuthMiddleware(ctx.SST))

	sampleR := auth.Group("/samples")
	{
		sampleR.GET("/sample/:id", HandlerWrapper(GetHandler, ctx))
		sampleR.POST("/sample", HandlerWrapper(PostHandler, ctx))
	}
}

func SSTRouter(ctx *AppOption, g *gin.RouterGroup) {
	auth := g.Group("", sstapp.GinAuthMiddleware(ctx.SST))

	sstR := auth.Group("/ssts")
	{
		// create server side token
		sstR.POST("/sst", HandlerWrapper(CreateSSTHandler, ctx))
	}

}

func SST2Router(ctx *AppOption, g *gin.RouterGroup) {
	auth := g.Group("", sstapp.GinAuthMiddleware(ctx.SST))

	sst2R := auth.Group("/sst2")
	{
		// create server side token
		sst2R.POST("/sst", HandlerWrapper(CreateSSTHandler, ctx))
	}
}
