package sst

import (
	"github.com/leyle/server-side-token/internal"
	"time"
)

func (sst *SSTokenOption) GenerateToken(userId string) (string, error) {
	srcUserId := encodeUserId(userId)
	token, err := internal.Encrypt(sst.aesKey, srcUserId)
	if err != nil {
		sst.logger.Error().Err(err).Str("userId", userId).Msg("GenerateToken failed")
		return "", err
	}

	sst.logger.Trace().Str("userId", userId).Str("token", token).Msg("GenerateToken succeed")
	sst.logger.Info().Str("userId", userId).Msg("GenerateToken succeed")
	return token, nil
}

func (sst *SSTokenOption) VerifyToken(token string) *OperateResult {
	// 1. token has been revoked
	// 2. token cannot be decrypted

	// 1. check if token has been revoked
	for _, rt := range sst.revokeList {
		if rt.token == token {
			sst.logger.Warn().Str("token", token).Msgf("token has been revoked at[%d]", rt.t)
			return checkTokenInvalid(token, "token has been revoked", rt.t)
		}
	}

	// 2. check if token can be decrypted
	text, err := internal.Decrypt(sst.aesKey, token)
	if err != nil {
		sst.logger.Error().Err(err).Str("token", token).Msg("decrypt aes token failed")
	}

	userId, createdAt, err := decodeUserId(text)
	if err != nil {
		sst.logger.Error().Err(err).Str("token", token).Msg("decode user id failed")
	}

	sst.logger.Debug().Str("userId", userId).Int64("t", createdAt).Msg("verify token, decode succeed")
	return checkTokenValid(token, createdAt)
}

func (sst *SSTokenOption) RevokeToken(token string) *OperateResult {
	// 1. check token is valid
	// 2. add token to revokeList
	// 3. add token to revoke list sqlite3

	// 1. check token is valid
	result := sst.VerifyToken(token)
	if !result.OK {
		sst.logger.Warn().Str("token", token).Msgf("revoke token, but token is invalid[%s|%d]", result.Reason, result.T)
		return revokeTokenFailed(token, result.Reason)
	}

	// 2. add token to revokeList
	sst.Lock()
	defer sst.Unlock()
	rv := &revokedToken{
		token: token,
		t:     time.Now().Unix(),
	}
	sst.revokeList = append(sst.revokeList, rv)

	// add token into db's revoke list
	err := sst.InsertIntoRevokeList(rv.token, rv.t)
	if err != nil {
		sst.logger.Error().Err(err).Msg("revoke token failed")
		return revokeTokenFailed(token, err.Error())
	}

	return revokeTokenSucceed(token, rv.t)
}
