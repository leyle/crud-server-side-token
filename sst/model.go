package sst

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"github.com/rs/zerolog"
	"sync"
)

type SSTokenOption struct {
	aesKey     []byte
	sqliteFile string
	logger     zerolog.Logger

	db         *sql.DB
	revokeList []*revokedToken

	sync.RWMutex
}

type OperateResult struct {
	Token  string
	OK     bool
	Reason string
	T      int64 // T means token last updated time
}

type revokedToken struct {
	token string
	t     int64
}

func NewSSTokenOption(aesKey, sqliteFile string, logger zerolog.Logger) (*SSTokenOption, error) {
	sst := &SSTokenOption{
		sqliteFile: sqliteFile,
		logger:     logger,
		revokeList: make([]*revokedToken, 0),
	}

	sst.aesKey = []byte(sst.getAesKey())

	// initial sqlite3 connection and create db?
	err := sst.GetDb()
	if err != nil {
		sst.logger.Error().Err(err).Msg("create new sst option failed")
		return nil, err
	}

	err = sst.CreateTable()
	if err != nil {
		return nil, err
	}

	// load revoke list data into memory
	err = sst.LoadRevokeList()
	if err != nil {
		return nil, err
	}

	return sst, nil
}

func (sst *SSTokenOption) getAesKey() string {
	m := md5.New()
	m.Write([]byte(sst.aesKey))
	return hex.EncodeToString(m.Sum(nil))
}

func checkTokenValid(token string, t int64) *OperateResult {
	// t is token creation time
	result := &OperateResult{
		Token: token,
		OK:    true,
		T:     t,
	}
	return result
}

func checkTokenInvalid(token, reason string, t int64) *OperateResult {
	// t is token invalid time or current time
	result := &OperateResult{
		Token:  token,
		OK:     false,
		Reason: reason,
		T:      t,
	}
	return result
}

func revokeTokenFailed(token, reason string) *OperateResult {
	result := &OperateResult{
		Token:  token,
		OK:     false,
		Reason: reason,
	}
	return result
}

func revokeTokenSucceed(token string, t int64) *OperateResult {
	// t means token revoke time
	result := &OperateResult{
		Token: token,
		OK:    true,
		T:     t,
	}
	return result
}
