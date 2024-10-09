package internal

import (
	"encoding/hex"
	"errors"
	"strings"
)

type CheckMsgV3Item struct {
	ExpectedHmac []byte
	Salt         []byte
	FileName     string
}

// ParseCheckMsgV3 parse the CheckMsgV3
//
//	checkMsgV3 string from info.xml
//	r1 []CheckMsgV3Item CheckMsgV3 Items
//	r2 error
func ParseCheckMsgV3(checkMsgV3 string) ([]CheckMsgV3Item, error) {
	if len(checkMsgV3) < 128 {
		return nil, errors.New("checkMsgV3 is less than 128 characters")
	}

	pendingItems := strings.Split(checkMsgV3, "**")
	items := make([]CheckMsgV3Item, 0, len(pendingItems))
	for _, pendingItem := range pendingItems {
		// pendingItem e.g. e56ac33a0eb3e97e501ded79eecc16496feb009a3ec46911186881f3dd73f3b7cec932efa6414914304a7e024f96686c38c7137bd734a407ba0a40d24696f813_com.tencent.mm514.tar
		strs := strings.Split(pendingItem, "_")
		if len(strs) != 2 {
			return nil, errors.New("assert split string length failed")
		}
		checkMsgV3PrefixStr := strs[0]
		filename := strs[1]

		expectedHmac, salt, err := parseCheckMsgV3ItemPrefixStr(checkMsgV3PrefixStr)
		if err != nil {
			return nil, err
		}

		items = append(items, CheckMsgV3Item{
			ExpectedHmac: expectedHmac,
			Salt:         salt,
			FileName:     filename,
		})
	}

	return items, nil
}

// parseCheckMsgV3ItemPrefixStr
//
//	r1 []byte expectedHmac
//	r2 []byte salt
//	err error
func parseCheckMsgV3ItemPrefixStr(checkMsgV3PrefixStr string) ([]byte, []byte, error) {
	if len(checkMsgV3PrefixStr) != 128 {
		return nil, nil, errors.New("checkMsgV3 prefix string is not 128 characters")
	}

	// split checkMsgV3PrefixStr
	expectedHmac, err := hex.DecodeString(checkMsgV3PrefixStr[0:64])
	if err != nil {
		return nil, nil, err
	}
	salt, err := hex.DecodeString(checkMsgV3PrefixStr[64:128])
	if err != nil {
		return nil, nil, err
	}

	return expectedHmac, salt, nil
}
