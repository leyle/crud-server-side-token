package main

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/rs/zerolog"
	"gitlab.com/fabric-crud-v2/crud-server-side-token/sstapp"
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

func HandlerWrapper(f func(ctx *AppOption), ctx *AppOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx.C = c
		f(ctx)
	}
}
