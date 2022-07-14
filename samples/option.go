package main

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/leyle/server-side-token/sstapp"
	"github.com/rs/zerolog"
)

type AppOption struct {
	C      *gin.Context
	Logger *zerolog.Logger
	Conf   *Config
	SST    *sstapp.SSTokenOption
}

func NewAppOption(conf *Config, sst *sstapp.SSTokenOption) *AppOption {
	logFormat := logmiddleware.LogTargetConsole
	if conf.Log.Format == LogFormatJson {
		logFormat = logmiddleware.LogTargetStdout
	}
	logger := logmiddleware.GetLogger(logFormat)
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
		SST:    op.SST,
	}
	return ctx
}

func HandlerWrapper(f func(ctx *AppOption), ctx *AppOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		nctx := ctx.New(c)
		f(nctx)
	}
}
