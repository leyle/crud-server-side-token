package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/leyle/server-side-token/sstapp"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	// cli commands
	// create ase key ./sstcli -createAesKey
	// create token: ./sstcli -createToken userid
	// verify token: ./sstcli -verifyToken token
	// revoke token: ./sstcli -revokeToken token

	var aesKeyLen int
	var userId string
	var vToken string
	var xToken string

	flag.IntVar(&aesKeyLen, "createAesKey", -1, "-createAesKey 20")
	flag.StringVar(&userId, "createToken", "", "-createToken $USER_ID")
	flag.StringVar(&vToken, "verifyToken", "", "-verifyToken $TOKEN")
	flag.StringVar(&xToken, "revokeToken", "", "-revokeToken $TOKEN")

	flag.Parse()

	logger := logmiddleware.GetLogger(logmiddleware.LogTargetConsole)

	if aesKeyLen > 0 {
		generateAesKey(aesKeyLen, &logger)
		return
	}

	if userId != "" {
		generateToken(userId, &logger)
		return
	}

	if vToken != "" {
		verifyToken(vToken, &logger)
		return
	}

	if xToken != "" {
		revokeToken(xToken, &logger)
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

func generateToken(userId string, logger *zerolog.Logger) string {
	userAesKey := readStdinAesKey(logger)
	sst, err := sstapp.NewSSTokenOption(userAesKey, logger)
	if err != nil {
		logger.Error().Err(err).Send()
		os.Exit(1)
	}
	token, err := sst.GenerateToken(userId)
	if err != nil {
		logger.Error().Err(err).Send()
		os.Exit(1)
	}
	logger.Info().Str("token", token).Send()

	return token
}

func verifyToken(token string, logger *zerolog.Logger) bool {
	userAesKey := readStdinAesKey(logger)
	sst, err := sstapp.NewSSTokenOption(userAesKey, logger)
	if err != nil {
		logger.Error().Err(err).Send()
		os.Exit(1)
	}

	result := sst.VerifyToken(token)
	if !result.OK {
		logger.Warn().Str("token", token).Msgf("invalid token, %s, maybe wrong token, or maybe wrong aes key", result.Msg)
		return false
	}

	logger.Info().Str("token", token).Str("userId", result.Msg).Msg("token is valid")
	return true
}

func revokeToken(token string, logger *zerolog.Logger) bool {
	userAesKey := readStdinAesKey(logger)
	sst, err := sstapp.NewSSTokenOption(userAesKey, logger)
	if err != nil {
		logger.Error().Err(err).Send()
		os.Exit(1)
	}

	result := sst.RevokeToken(token)
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
