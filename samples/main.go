package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/leyle/go-api-starter/confighelper"
	"github.com/leyle/go-api-starter/ginhelper"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/leyle/server-side-token/sstapp"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var err error
	var cfile string
	var selfToken string
	flag.StringVar(&cfile, "c", "", "-c /path/to/config")
	flag.StringVar(&selfToken, "t", "", "usage -t admin")

	flag.Parse()
	if cfile == "" {
		fmt.Println("no config file, usage: -c /path/to/config")
		os.Exit(1)
	}

	var conf *Config
	err = confighelper.LoadConfig(context.Background(), cfile, &conf)
	if err != nil {
		fmt.Println("parse config file failed")
		fmt.Println(err)
		os.Exit(1)
	}

	// set log format
	logFormat := logmiddleware.LogTargetConsole
	if conf.Log.Format == LogFormatJson {
		logFormat = logmiddleware.LogTargetStdout
	}
	logger := logmiddleware.GetLogger(logFormat)

	logger.Debug().Str("selfName", conf.SST.ServiceName).Send()

	// initial server side token
	aesKey := conf.SST.AesKey
	if aesKey == "" {
		fmt.Println("invalid sst(server side token) config values, aesKey is empty")
		os.Exit(1)
	}
	sst, err := sstapp.NewSSTokenOption(conf.SST.ServiceName, aesKey)
	if err != nil {
		fmt.Println("create server side token object failed")
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx = logger.WithContext(ctx)
	if selfToken != "" {
		fmt.Println("create itself server side token...")
		selfTokenVal, err := sst.GenerateToken(ctx, selfToken)
		if err != nil {
			fmt.Println("generate itself server side token failed")
			os.Exit(1)
		}
		fmt.Println("generate itself server side token succeed")
		fmt.Println(selfTokenVal)
		os.Exit(1)
	}

	// start http server
	ap := NewAppOption(conf, sst)

	go httpServer(ap)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	os.Exit(1)
}

func httpServer(ctx *AppOption) {
	e := ginhelper.SetupGin(*ctx.Logger)
	ginhelper.PrintHeaders = ctx.Conf.Debug

	router := e.Group("/api")

	SampleRouter(ctx, router.Group(""))

	// sst manage api
	SSTRouter(ctx, router.Group(""))

	addr := ctx.Conf.Server.ListenServerAddr()
	err := e.Run(addr)
	if err != nil {
		fmt.Println(err)
		ctx.Logger.Error().Err(err).Msg("start http server failed")
		return
	}
}
