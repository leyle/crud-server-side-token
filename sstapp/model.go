package sstapp

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"github.com/rs/zerolog"
	"sync"
)

const (
	ServerSideTokenHeaderName = "X-Server-Side-Token"
)

type SSTokenOption struct {
	aesKey     []byte
	sqliteFile string
	logger     *zerolog.Logger

	db         *sql.DB
	revokeList []*revokedToken

	mutex sync.RWMutex
}

type OperateResult struct {
	Token string
	OK    bool
	Msg   string
	T     int64 // T means token last updated time
}

type revokedToken struct {
	token  string
	userId string
	t      int64
}

func NewSSTokenOption(aesKey, sqliteFile string, logger *zerolog.Logger) (*SSTokenOption, error) {
	sst := &SSTokenOption{
		sqliteFile: sqliteFile,
		logger:     logger,
		revokeList: make([]*revokedToken, 0),
	}

	sst.aesKey = []byte(sst.getAesKey(aesKey))

	// initial sqlite3 connection and create db?
	err := sst.getDb()
	if err != nil {
		sst.logger.Error().Err(err).Msg("create new sstapp option failed")
		return nil, err
	}

	err = sst.createTable()
	if err != nil {
		return nil, err
	}

	// load revoke list data into memory
	err = sst.loadRevokeList()
	if err != nil {
		return nil, err
	}

	return sst, nil
}

func (sst *SSTokenOption) New(logger *zerolog.Logger) *SSTokenOption {
	netSST := &SSTokenOption{
		aesKey:     sst.aesKey,
		sqliteFile: sst.sqliteFile,
		logger:     logger,
		revokeList: sst.revokeList,
	}
	return netSST
}

func (sst *SSTokenOption) getAesKey(key string) string {
	m := md5.New()
	m.Write([]byte(key))
	return hex.EncodeToString(m.Sum(nil))
}

func checkTokenValid(token, userId string, t int64) *OperateResult {
	// t is token creation time
	result := &OperateResult{
		Token: token,
		OK:    true,
		Msg:   userId,
		T:     t,
	}
	return result
}

func checkTokenInvalid(token, reason string, t int64) *OperateResult {
	// t is token invalid time or current time
	result := &OperateResult{
		Token: token,
		OK:    false,
		Msg:   reason,
		T:     t,
	}
	return result
}

func revokeTokenFailed(token, reason string) *OperateResult {
	result := &OperateResult{
		Token: token,
		OK:    false,
		Msg:   reason,
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
