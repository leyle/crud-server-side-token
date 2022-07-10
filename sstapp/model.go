package sstapp

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"sync"
)

const (
	ServerSideTokenHeaderName = "X-Server-Side-Token"
)

const (
	sqliteCfgPath    = ".config/sst"
	sqliteDbFilename = "sst.db"
)

const sstPrefix = "SST-"

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
	T     int64
	Err   error
}

type revokedToken struct {
	token  string
	userId string
	t      int64
}

func NewSSTokenOption(aesKey string, logger *zerolog.Logger) (*SSTokenOption, error) {
	sst := &SSTokenOption{
		logger:     logger,
		revokeList: make([]*revokedToken, 0),
	}

	err := sst.insureSqliteFile()
	if err != nil {
		return nil, err
	}

	sst.aesKey = []byte(sst.getAesKey(aesKey))

	// initial sqlite3 connection and create db?
	err = sst.getDb()
	if err != nil {
		sst.logger.Error().Err(err).Msg("create new sstapp option failed")
		return nil, err
	}

	err = sst.createTable()
	if err != nil {
		return nil, err
	}

	// load revoke list data into memory
	err = sst.loadRevocationList()
	if err != nil {
		return nil, err
	}

	return sst, nil
}

func (sst *SSTokenOption) insureSqliteFile() error {
	home, err := os.UserHomeDir()
	if err != nil {
		sst.logger.Error().Err(err).Msg("create sqlite path, get user home dir failed")
		return err
	}

	sstDbPath := fmt.Sprintf("%s/%s", home, sqliteCfgPath)
	err = os.MkdirAll(sstDbPath, os.ModePerm)
	if err != nil {
		sst.logger.Error().Err(err).Msg("create sqlite path failed")
		return err
	}

	sqliteDbPath := fmt.Sprintf("%s/%s", sstDbPath, sqliteDbFilename)
	sst.sqliteFile = sqliteDbPath

	return nil
}

func (sst *SSTokenOption) SqliteFilePath() string {
	return sst.sqliteFile
}

func (sst *SSTokenOption) Copy(logger *zerolog.Logger) *SSTokenOption {
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

func (sst *SSTokenOption) packSSToken(cipher string) string {
	return fmt.Sprintf("%s%s", sstPrefix, cipher)
}

func (sst *SSTokenOption) unpackSSToken(msg string) (string, error) {
	if !strings.HasPrefix(msg, sstPrefix) {
		return "", ErrInvalidTokenFormat
	}

	cutMsg := strings.SplitN(msg, sstPrefix, 2)

	return cutMsg[1], nil
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

func checkTokenInvalid(token, reason string, t int64, err error) *OperateResult {
	// t is token invalid time or current time
	result := &OperateResult{
		Token: token,
		OK:    false,
		Msg:   reason,
		T:     t,
		Err:   err,
	}
	return result
}

func revokeTokenFailed(token, reason string, err error) *OperateResult {
	result := &OperateResult{
		Token: token,
		OK:    false,
		Msg:   reason,
		Err:   err,
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
