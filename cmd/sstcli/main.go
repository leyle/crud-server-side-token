package main

import (
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/leyle/server-side-token/sstapp"
)

func main() {
	// cli command
	// create token: ./sstcli -createToken userid
	// revoke token: ./sstcli -revokeToken token
	// verify token: ./sstcli -verifyToken token

	aesKey := "abc"
	dbFile := "/tmp/test.db"
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetStdout)
	sst, _ := sstapp.NewSSTokenOption(aesKey, dbFile, logger)
	sst.GenerateToken("userid")
}
