package sst

import (
	"github.com/leyle/go-api-starter/logmiddleware"
	"strconv"
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

func TestVerifyToken(t *testing.T) {
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole)
	aesKey := "^ct9<.yT3CK*MQ6j/V"
	sqlFile := "/tmp/test.db"
	sst, err := NewSSTokenOption(aesKey, sqlFile, logger)
	if err != nil {
		t.Error(err)
	}

	// token := "nb98xra8HqsxlymV3M4vFcus8FJvUsSnrGD8kRt09fg6CAi5OpuJAkODKTiN1W5k"
	token := "K8R1LPIiqHkxkrb-cA5Kwe3TUi86sgxnVrp1VNpfl04-p2iawOOOlAA6TnzmJvGk"

	result := sst.VerifyToken(token)
	t.Log(result.OK)
	t.Log(result)
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
	err = sst.InsertIntoRevokeList(token, userId, time.Now().Unix())
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole)
	aesKey := "^ct9<.yT3CK*MQ6j/V"
	sqlFile := "/tmp/test.db"
	sst, err := NewSSTokenOption(aesKey, sqlFile, logger)
	if err != nil {
		b.Error(err)
	}

	userId := "cdi-service"

	for n := 0; n < b.N; n++ {
		_, err = sst.GenerateToken(userId)
		if err != nil {
			b.Error(err)
		}
		// b.Log(token)
	}
}

func getSSTOpt() *SSTokenOption {
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole)
	aesKey := "^ct9<.yT3CK*MQ6j/V"
	sqlFile := "/tmp/test.db"
	sst, err := NewSSTokenOption(aesKey, sqlFile, logger)
	if err != nil {
		panic(err)
	}
	return sst
}

func BenchmarkRevokeToken(b *testing.B) {
	sst := getSSTOpt()

	for n := 0; n < b.N; n++ {
		// create token then revoke it
		userId := strconv.FormatInt(time.Now().Unix(), 10)

		token, _ := sst.GenerateToken(userId)
		sst.RevokeToken(token)
	}
}
