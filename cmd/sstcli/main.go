package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/leyle/server-side-token/sstapp"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var logger = logmiddleware.GetLogger(logmiddleware.LogTargetConsole)

func main() {
	// cli commands
	// create ase key ./sstcli -createAesKey
	// create token: ./sstcli -secretFile /path/to/key.yaml -createToken userid
	// verify token: ./sstcli -secretFile /path/to/key.yaml -verifyToken token
	// revoke token: ./sstcli -secretFile /path/to/key.yaml -revokeToken token

	var aesKeyLen int
	var userId string
	var vToken string
	var xToken string

	var secretFile string

	flag.IntVar(&aesKeyLen, "createAesKey", -1, "-createAesKey 20")
	flag.StringVar(&userId, "createToken", "", "-createToken $USER_ID")
	flag.StringVar(&vToken, "verifyToken", "", "-verifyToken $TOKEN")
	flag.StringVar(&xToken, "revokeToken", "", "-revokeToken $TOKEN")
	flag.StringVar(&secretFile, "secretFile", "", "-secretFile /some/yaml/file/path")

	flag.Parse()

	if aesKeyLen > 0 {
		generateAesKey(aesKeyLen, &logger)
		return
	}

	// below commands needs secret key
	if secretFile == "" {
		logger.Error().Msg("no secret yaml file")
		flag.PrintDefaults()
		os.Exit(1)
	}

	secretKey, err := readSecretFile(secretFile, &logger)
	if err != nil {
		os.Exit(1)
	}

	if userId != "" {
		generateToken(userId, secretKey, &logger)
		return
	}

	if vToken != "" {
		verifyToken(vToken, secretKey)
		return
	}

	if xToken != "" {
		revokeToken(xToken, secretKey)
		return
	}

	// default behavior, print usage msg
	flag.PrintDefaults()
}

func generateAesKey(keyLen int, logger *zerolog.Logger) string {
	aesKey := generateRandomKey(keyLen)
	logger.Info().Str("aesKey", aesKey).Msg("Generate aesKey succeeded")
	return aesKey
}

func generateToken(userId, secretKey string, logger *zerolog.Logger) string {
	sst, err := sstapp.NewSSTokenOption(secretKey)
	if err != nil {
		logger.Error().Err(err).Send()
		os.Exit(1)
	}
	token, err := sst.GenerateToken(getContext(), userId)
	if err != nil {
		logger.Error().Err(err).Send()
		os.Exit(1)
	}
	logger.Info().Str("token", token).Send()

	return token
}

func verifyToken(token, secretKey string) bool {
	sst, err := sstapp.NewSSTokenOption(secretKey)
	if err != nil {
		logger.Error().Err(err).Send()
		os.Exit(1)
	}

	result := sst.VerifyToken(getContext(), token)
	if !result.OK {
		logger.Warn().Str("token", token).Msgf("invalid token, %s, maybe wrong token, or maybe wrong aes key", result.Msg)
		return false
	}

	logger.Info().Str("token", token).Str("userId", result.Msg).Msg("token is valid")
	return true
}

func revokeToken(token, secretKey string) bool {
	sst, err := sstapp.NewSSTokenOption(secretKey)
	if err != nil {
		logger.Error().Err(err).Send()
		os.Exit(1)
	}

	result := sst.RevokeToken(getContext(), token)
	if !result.OK {
		if result.Err == sstapp.ErrRevokeAlreadyRevoked {
			logger.Info().Str("token", token).Msg("revoke token succeeded")
			return true
		}
		logger.Warn().Str("token", token).Msgf("revoke token failed, %s", result.Msg)
		return false
	}

	logger.Info().Str("token", token).Str("userId", result.Msg).Msg("revoke token succeeded")
	return true
}

func generateRandomKey(length int) string {
	base := "0123456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ#$%&*+,-./:;<=>?@^{}"
	baseLen := len(base)

	data := make([]byte, length)
	rand.Read(data)
	for i := 0; i < length; i++ {
		data[i] = base[int(data[i])%baseLen]
	}

	return string(data)
}

func readStdinAesKey(logger *zerolog.Logger) string {
	var userAesKey string
	for {
		logger.Info().Msg("input aes key:")
		n, err := fmt.Scanln(&userAesKey)
		if n <= 0 {
			logger.Error().Msg("no aes key")
			continue
		}
		if err != nil {
			logger.Error().Err(err).Send()
			continue
		}
		break
	}

	return userAesKey
}

func readSecretFile(filepath string, logger *zerolog.Logger) (string, error) {
	// read secret key yaml file, response is secret string

	type KeyCfg struct {
		Key string `yaml:"key"`
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Error().Err(err).Msg("read secret key yaml file failed")
		return "", err
	}

	cfg := KeyCfg{}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logger.Error().Err(err).Msg("parse secret file failed")
		return "", err
	}

	return cfg.Key, nil
}

func getContext() context.Context {
	reqId := logmiddleware.GenerateReqId()
	logger = logger.With().Str("reqId", reqId).Logger()
	lctx := logger.WithContext(context.Background())
	return lctx
}
