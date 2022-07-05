package sst

import (
	"github.com/leyle/go-api-starter/logmiddleware"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole)
	aesKey := "^ct9<.yT3CK*MQ6j/V"
	sqlFile := "test.db"
	sst, err := NewSSTokenOption(aesKey, sqlFile, logger)
	if err != nil {
		t.Error(err)
	}

	userId := "cdi-service"

	token, err := sst.GenerateToken(userId)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
}

func TestRevokeToken(t *testing.T) {
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole)
	aesKey := "^ct9<.yT3CK*MQ6j/V"
	sqlFile := "/tmp/test.db"
	sst, err := NewSSTokenOption(aesKey, sqlFile, logger)
	if err != nil {
		t.Error(err)
	}

	userId := "cdi-service"

	token, err := sst.GenerateToken(userId)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)

	result := sst.RevokeToken(token)
	t.Log(result)

	for _, rv := range sst.revokeList {
		t.Log(rv)
	}

}

func TestSQLiteOpt(t *testing.T) {
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole)
	aesKey := "^ct9<.yT3CK*MQ6j/V"
	sqlFile := "/tmp/test.db"
	sst, err := NewSSTokenOption(aesKey, sqlFile, logger)
	if err != nil {
		t.Error(err)
	}

	userId := "cdi-service"

	token, err := sst.GenerateToken(userId)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)

	// token = "_xMqRdnhyzBgYbnXeAuEa8CN5sMN8O_zIlIcDYLqAZQvp84zG93SsHGIILQb5wgU"
	err = sst.InsertIntoRevokeList(token, time.Now().Unix())
	if err != nil {
		t.Error(err)
	}
}
