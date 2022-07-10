package sstapp

import (
	"time"
)

func (sst *SSTokenOption) GenerateToken(userId string) (string, error) {
	srcUserId := encodeUserId(userId)
	// cipher, err := internal.Encrypt(sst.aesKey, srcUserId)
	cipher, err := sst.encrypt([]byte(srcUserId))
	if err != nil {
		sst.logger.Error().Err(err).Str("userId", userId).Msg("GenerateToken failed")
		return "", err
	}

	token := sst.packSSToken(cipher)

	sst.logger.Trace().Str("userId", userId).Str("token", token).Msg("GenerateToken succeed")
	sst.logger.Info().Str("userId", userId).Msg("GenerateToken succeed")
	return token, nil
}

func (sst *SSTokenOption) VerifyToken(token string) *OperateResult {

	// 0. check token format
	// 1. token has been revoked
	// 2. token cannot be decrypted

	// 0. check token format
	cipher, err := sst.unpackSSToken(token)
	if err != nil {
		sst.logger.Warn().Err(err).Msg("verify token failed")
		return checkTokenInvalid(token, err.Error(), 0, err)
	}

	// 1. check if token has been revoked
	for _, rt := range sst.revokeList {
		if rt.token == token {
			sst.logger.Warn().Str("token", token).Msgf("token has been revoked at[%d]", rt.t)
			return checkTokenInvalid(token, "token has been revoked", rt.t, ErrRevokeAlreadyRevoked)
		}
	}

	// 2. check if token can be decrypted
	// text, err := internal.Decrypt(sst.aesKey, cipher)
	text, err := sst.decrypt(cipher)
	if err != nil {
		sst.logger.Warn().Err(err).Str("token", token).Msg("decrypt aes token failed")
		return checkTokenInvalid(token, "invalid token format", 0, ErrDecryptMsgFailed)
	}

	userId, createdAt, err := decodeUserId(text)
	if err != nil {
		sst.logger.Warn().Err(err).Str("token", token).Msg("decode user id failed")
		return checkTokenInvalid(token, "invalid token format, maybe old version", 0, ErrInvalidRawMsgFormat)
	}

	sst.logger.Debug().Str("userId", userId).Int64("t", createdAt).Msg("verify token, decode succeed")
	return checkTokenValid(token, userId, createdAt)
}

func (sst *SSTokenOption) RevokeToken(token string) *OperateResult {
	// 1. check token is valid
	// 2. add token to revokeList
	// 3. add token to revoke list sqlite3

	// 1. check token is valid
	result := sst.VerifyToken(token)
	if !result.OK {
		sst.logger.Warn().Str("token", token).Msgf("revoke token, but token is invalid[%s|%d]", result.Msg, result.T)
		return revokeTokenFailed(token, result.Msg, result.Err)
	}

	// 2. add token to revokeList
	sst.mutex.Lock()
	defer sst.mutex.Unlock()
	rv := &revokedToken{
		token:  token,
		userId: result.Msg,
		t:      time.Now().Unix(),
	}
	sst.revokeList = append(sst.revokeList, rv)

	// 3. add token into db's revocation list
	err := sst.insertIntoRevocationList(rv.token, rv.userId, rv.t)
	if err != nil {
		sst.logger.Error().Err(err).Msg("revoke token failed")
		return revokeTokenFailed(token, err.Error(), ErrSaveDBFailed)
	}

	return revokeTokenSucceed(token, rv.t)
}
