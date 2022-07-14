package sstapp

import "errors"

var (
	ErrServiceNameEmpty     = errors.New("service name is empty string")
	ErrAesKeyEmpty          = errors.New("aes key is empty string")
	ErrRevokeAlreadyRevoked = errors.New("token has been revoked")
	ErrDecryptMsgFailed     = errors.New("decrypt msg failed")
	ErrInvalidRawMsgFormat  = errors.New("invalid raw msg format")
	ErrInvalidTokenFormat   = errors.New("invalid token format")
	ErrSaveDBFailed         = errors.New("save revoke token into db failed")
)
