package sstapp

import "errors"

var (
	ErrRevokeAlreadyRevoked = errors.New("token has been revoked")
	ErrDecryptMsgFailed     = errors.New("decrypt msg failed")
	ErrInvalidRawMsgFormat  = errors.New("invalid raw msg format")
	ErrSaveDBFailed         = errors.New("save revoke token into db failed")
)
