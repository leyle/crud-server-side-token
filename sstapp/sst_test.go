package sstapp

import (
	"context"
	"github.com/leyle/go-api-starter/logmiddleware"
	"strconv"
	"testing"
	"time"
)

const CtxLoggerName = "x-req-id"

func getSSTOpt() *SSTokenOption {
	aesKey := "Jzw8C%E/y,FSK4<}n?"
	sst, err := NewSSTokenOption(aesKey)
	if err != nil {
		panic(err)
	}
	return sst
}

func getContext() context.Context {
	id := logmiddleware.GenerateReqId()
	ctx := context.Background()

	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole).With().Str(CtxLoggerName, id).Logger()

	lctx := logger.WithContext(ctx)

	return lctx
}

func TestGenerateToken(t *testing.T) {
	ctx := getContext()
	sst := getSSTOpt()

	userId := "cdi-service"

	token, err := sst.GenerateToken(ctx, userId)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
}

func TestRunTimes(t *testing.T) {
	n := 0
	count := 10
	for n < count {
		TestGenerateToken(t)
		n++
	}
}

func TestVerifyToken(t *testing.T) {
	ctx := getContext()
	sst := getSSTOpt()

	// token := "nb98xra8HqsxlymV3M4vFcus8FJvUsSnrGD8kRt09fg6CAi5OpuJAkODKTiN1W5k"
	// token := "SST-8R1LPIiqHkxkrb-cA5Kwe3TUi86sgxnVrp1VNpfl04-p2iawOOOlAA6TnzmJvGk"
	token, err := sst.GenerateToken(ctx, "hello-cdi")
	if err != nil {
		t.Error(err)
	}

	result := sst.VerifyToken(ctx, token)
	t.Log(result.OK)
	t.Log(result)
}

func TestVerifyInvalidToken(t *testing.T) {
	ctx := getContext()
	sst := getSSTOpt()

	token := "SST-8hg9zF/npLEuXuHFlO/drzM0O5e8iHM8RUNhCepgnyARkwCVM8UIa26IEdbnU8Lx"

	result := sst.VerifyToken(ctx, token)

	t.Log(result)

	if result.OK {
		t.FailNow()
	}
}

func TestRevokeToken(t *testing.T) {
	ctx := getContext()
	sst := getSSTOpt()

	userId := "cdi-service"

	token, err := sst.GenerateToken(ctx, userId)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)

	result := sst.RevokeToken(ctx, token)
	t.Log(result)

	for _, rv := range sst.revokeList {
		t.Log(rv)
	}

}

func TestSQLiteOpt(t *testing.T) {
	ctx := getContext()
	sst := getSSTOpt()

	userId := "cdi-service"

	token, err := sst.GenerateToken(ctx, userId)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)

	// token = "_xMqRdnhyzBgYbnXeAuEa8CN5sMN8O_zIlIcDYLqAZQvp84zG93SsHGIILQb5wgU"
	err = sst.insertIntoRevocationList(token, userId, time.Now().Unix())
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	var err error
	ctx := getContext()
	sst := getSSTOpt()

	userId := "cdi-service"

	for n := 0; n < b.N; n++ {
		_, err = sst.GenerateToken(ctx, userId)
		if err != nil {
			b.Error(err)
		}
		// b.Log(token)
	}
}

func BenchmarkRevokeToken(b *testing.B) {
	ctx := getContext()
	sst := getSSTOpt()

	for n := 0; n < b.N; n++ {
		// create token then revoke it
		userId := strconv.FormatInt(time.Now().Unix(), 10)

		token, _ := sst.GenerateToken(ctx, userId)
		sst.RevokeToken(ctx, token)
	}
}
