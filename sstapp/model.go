package sstapp

import (
	"database/sql"
	"fmt"
	"github.com/leyle/crud-server-side-token/internal"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	ServerSideTokenHeaderName = "X-Server-Side-Token"
)

var (
	SQLiteCfgPath    = ".config/sst"
	SQLiteDbFilename = "sst.db"
)

var singleSST *SSTokenOption
var lock = &sync.RWMutex{}

type SSTokenOption struct {
	aesKey     []byte
	sqliteFile string
	logger     *zerolog.Logger

	tokenPrefix string

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

func NewSSTokenOption(serviceName, aesKey string) (*SSTokenOption, error) {
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole)

	if singleSST != nil {
		logger.Info().Msg("sst token instance already created")
		return singleSST, nil
	}

	lock.Lock()
	defer lock.Unlock()

	if singleSST != nil {
		logger.Info().Msg("sst token instance already created")
		return singleSST, nil
	}

	if serviceName == "" {
		logger.Error().Msg("service name or aes key can't be empty string")
		return nil, ErrServiceNameEmpty
	}

	if aesKey == "" {
		logger.Error().Msg("aes key can't be empty string")
		return nil, ErrAesKeyEmpty
	}

	sst := &SSTokenOption{
		aesKey:     []byte(aesKey),
		logger:     &logger,
		revokeList: make([]*revokedToken, 0),
	}
	sst.tokenPrefix = sst.getTokenPrefix(serviceName)

	err := sst.insureSqliteFile()
	if err != nil {
		return nil, err
	}

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

	singleSST = sst

	logger.Info().Msg("create new server side token instance successfully")

	return sst, nil
}

func (sst *SSTokenOption) insureSqliteFile() error {
	sstDbPath, err := sst.getDBPath()
	if err != nil {
		sst.logger.Error().Err(err).Msg("create sqlite path failed")
		return err
	}

	sqliteDbPath := fmt.Sprintf("%s/%s", sstDbPath, SQLiteDbFilename)
	sst.sqliteFile = sqliteDbPath

	return nil
}

func (sst *SSTokenOption) getDBPath() (string, error) {
	dbPath := SQLiteCfgPath
	if !filepath.IsAbs(SQLiteCfgPath) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dbPath = fmt.Sprintf("%s/%s", home, SQLiteCfgPath)
	}

	// double insure path exist
	err := os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dbPath, nil
}

func (sst *SSTokenOption) SqliteFilePath() string {
	return sst.sqliteFile
}

func (sst *SSTokenOption) getTokenPrefix(service string) string {
	upperCase := strings.ToUpper(service)
	prefix := fmt.Sprintf("%s%s-", "SST-", upperCase)
	return prefix
}

func (sst *SSTokenOption) encrypt(userId []byte) (string, error) {
	cipherText, err := internal.GcmEncrypt(sst.aesKey, userId)
	if err != nil {
		sst.logger.Error().Err(err).Msg("encrypt userId failed")
		return "", err
	}

	b64str := internal.HexEncodeCipherText(cipherText)

	return b64str, nil
}

func (sst *SSTokenOption) decrypt(b64CipherText string) (string, error) {
	cipherText, err := internal.HexDecodeCipherString(b64CipherText)
	if err != nil {
		sst.logger.Warn().Err(err).Msg("decode base64 cipher text failed")
		return "", err
	}

	plainBytes, err := internal.GcmDecrypt(sst.aesKey, cipherText)
	if err != nil {
		sst.logger.Warn().Err(err).Msg("decrypt cipher text failed")
		return "", err
	}

	return string(plainBytes), nil
}

func (sst *SSTokenOption) packSSToken(cipher string) string {
	return fmt.Sprintf("%s%s", sst.tokenPrefix, cipher)
}

func (sst *SSTokenOption) unpackSSToken(msg string) (string, error) {
	if !strings.HasPrefix(msg, sst.tokenPrefix) {
		return "", ErrInvalidTokenFormat
	}

	cutMsg := strings.SplitN(msg, sst.tokenPrefix, 2)

	return cutMsg[1], nil
}

func createTokenOK(token, userId string) *OperateResult {
	result := &OperateResult{
		Token: token,
		OK:    true,
		Msg:   userId,
		T:     time.Now().Unix(),
	}
	return result
}

func createTokenFailed(userId, reason string, err error) *OperateResult {
	result := &OperateResult{
		Token: userId,
		OK:    false,
		Msg:   reason,
		T:     time.Now().Unix(),
		Err:   err,
	}
	return result
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

func revokeTokenOK(token string, t int64) *OperateResult {
	// t means token revoke time
	result := &OperateResult{
		Token: token,
		OK:    true,
		T:     t,
	}
	return result
}
