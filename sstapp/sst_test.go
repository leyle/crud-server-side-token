package sstapp

import (
	"github.com/leyle/go-api-starter/logmiddleware"
	"strconv"
	"testing"
	"time"
)

func getSSTOpt() *SSTokenOption {
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole)
	aesKey := "^ct9<.yT3CK*MQ6j/V"
	sst, err := NewSSTokenOption(aesKey, &logger)
	if err != nil {
		panic(err)
	}
	return sst
}

func TestGenerateToken(t *testing.T) {
	sst := getSSTOpt()

	userId := "cdi-service"

	token, err := sst.GenerateToken(userId)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
}

func TestVerifyToken(t *testing.T) {
	sst := getSSTOpt()

	// token := "nb98xra8HqsxlymV3M4vFcus8FJvUsSnrGD8kRt09fg6CAi5OpuJAkODKTiN1W5k"
	// token := "SST-8R1LPIiqHkxkrb-cA5Kwe3TUi86sgxnVrp1VNpfl04-p2iawOOOlAA6TnzmJvGk"
	token, err := sst.GenerateToken("hello-cdi")
	if err != nil {
		t.Error(err)
	}

	result := sst.VerifyToken(token)
	t.Log(result.OK)
	t.Log(result)
}

func TestVerifyInvalidToken(t *testing.T) {
	sst := getSSTOpt()

	// token := "abcinvalidtoken"
	// token := "4JqjLbYWaM2Fos0Tg-PgRYBiAm3rNS2WSnLvThKWvdC034JSkprv7rOhwyocIKnx"
	token := "SST-SST-Ukuw_3pFVHZGzoWh2-2PM6_wUE2ZwEmXqhOhfHkQH1Nl8TzjiI3pkFREU0H2zbkJ"

	result := sst.VerifyToken(token)

	t.Log(result)

	if result.OK {
		t.FailNow()
	}
}

func TestRevokeToken(t *testing.T) {
	sst := getSSTOpt()

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
	sst := getSSTOpt()

	userId := "cdi-service"

	token, err := sst.GenerateToken(userId)
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
	sst := getSSTOpt()

	userId := "cdi-service"

	for n := 0; n < b.N; n++ {
		_, err = sst.GenerateToken(userId)
		if err != nil {
			b.Error(err)
		}
		// b.Log(token)
	}
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
