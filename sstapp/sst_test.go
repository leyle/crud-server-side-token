package sstapp

import (
	"context"
	"github.com/leyle/go-api-starter/logmiddleware"
	"strconv"
	"sync"
	"testing"
	"time"
)

const CtxLoggerName = "x-req-id"

const serviceName = "library"

func getSSTOpt() *SSTokenOption {
	aesKey := "Jzw8C%E/y,FSK4<}n?"
	sst, err := NewSSTokenOption(serviceName, aesKey)
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
	SQLiteCfgPath = ".config/fabric-state/sst"
	ctx := getContext()
	sst := getSSTOpt()

	userId := "cdi-service"

	result := sst.GenerateToken(ctx, userId)
	if result.Err != nil {
		t.Error(result.Err)
	}
	t.Log(result.Token)
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
	cresult := sst.GenerateToken(ctx, "hello-cdi")
	if cresult.Err != nil {
		t.Error(cresult.Err)
	}

	result := sst.VerifyToken(ctx, cresult.Token)
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

	cresult := sst.GenerateToken(ctx, userId)
	if cresult.Err != nil {
		t.Error(cresult.Err)
	}
	t.Log(cresult.Token)

	result := sst.RevokeToken(ctx, cresult.Token)
	t.Log(result)

	for _, rv := range sst.revokeList {
		t.Log(rv)
	}

}

func TestSQLiteOpt(t *testing.T) {
	ctx := getContext()
	sst := getSSTOpt()

	userId := "cdi-service"

	cresult := sst.GenerateToken(ctx, userId)
	if cresult.Err != nil {
		t.Error(cresult.Err)
	}
	t.Log(cresult.Token)

	// token = "_xMqRdnhyzBgYbnXeAuEa8CN5sMN8O_zIlIcDYLqAZQvp84zG93SsHGIILQb5wgU"
	err := sst.insertIntoRevocationList(cresult.Token, userId, time.Now().Unix())
	if err != nil {
		t.Error(err)
	}
}

func TestSingleTon(t *testing.T) {
	count := 5
	var wg sync.WaitGroup

	createSST := func() {
		defer wg.Done()
		getSSTOpt()
	}

	for i := 0; i < count; i++ {
		wg.Add(1)
		go createSST()
	}

	wg.Wait()
}

func BenchmarkGenerateToken(b *testing.B) {
	ctx := getContext()
	sst := getSSTOpt()

	userId := "cdi-service"

	for n := 0; n < b.N; n++ {
		result := sst.GenerateToken(ctx, userId)
		if result.Err != nil {
			b.Error(result.Err)
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

		result := sst.GenerateToken(ctx, userId)
		sst.RevokeToken(ctx, result.Token)
	}
}
