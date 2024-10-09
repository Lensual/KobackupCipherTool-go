package internal

import (
	"encoding/hex"
	"errors"
)

type EncMsgV3 struct {
	Salt []byte
	Iv   []byte
}

// ParseEncMsgV3 parse the encMsgV3
//
//	password string the password
//	encMsgV3 string from info.xml
//	r1 []byte aesKey
//	r2 []byte iv
//	err error
func ParseEncMsgV3(password string, encMsgV3 string) (r EncMsgV3, err error) {
	if len(encMsgV3) != 96 {
		return r, errors.New("encMsgV3 must be 96 characters")
	}

	r.Salt, err = hex.DecodeString(encMsgV3[0:64])
	if err != nil {
		return r, err
	}
	r.Iv, err = hex.DecodeString(encMsgV3[64:])
	if err != nil {
		return r, err
	}

	return r, nil
}
