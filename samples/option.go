package main

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/server-side-token/sstapp"
	"github.com/rs/zerolog"
)

type AppOption struct {
	C      *gin.Context
	Logger *zerolog.Logger
	Conf   *Config
	SST    *sstapp.SSTokenOption
}

func NewAppOption(conf *Config, sst *sstapp.SSTokenOption, logger zerolog.Logger) *AppOption {
	op := &AppOption{
		Conf:   conf,
		SST:    sst,
		Logger: &logger,
	}
	return op
}

func (op *AppOption) New(c *gin.Context) *AppOption {
	logger := zerolog.Ctx(c.Request.Context())
	ctx := &AppOption{
		C:      c,
		Logger: logger,
		Conf:   op.Conf,
		SST:    op.SST.New(logger),
	}
	return ctx
}

func HandlerWrapper(f func(ctx *AppOption), ctx *AppOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		nctx := ctx.New(c)
		f(nctx)
	}
}
