package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func main() {
	argPassword := flag.String("password", "", "")
	argCheckMsgV3 := flag.String("checkMsgV3", "", "")
	argInput := flag.String("input", "", "")
	flag.Parse()

	inputFile, err := os.OpenFile(*argInput, os.O_RDONLY, 0)
	if err != nil {
		log.Printf("Input File can't open: %e", err)
		return
	}

	result, err := Verify(*argPassword, *argCheckMsgV3, inputFile)
	if err != nil {
		log.Printf("Verify Failed: %e", err)
	}

	log.Printf("Verify: %t", result)
}

func Verify(password string, checkMsgV3 string, inputReader io.Reader) (bool, error) {
	if len(checkMsgV3) != 128 {
		return false, errors.New("checkMsgV3 length is not 128 characters")
	}

	// split checkMsgV3
	expectedHmac, err := hex.DecodeString(checkMsgV3[0:64])
	if err != nil {
		return false, err
	}
	salt, err := hex.DecodeString(checkMsgV3[64:128])
	if err != nil {
		return false, err
	}

	// parse hmacKey
	hmacKey := pbkdf2.Key([]byte(password), salt, 5000, 32, sha256.New)

	// hmacSha256(input, hexEncode(hmacKey))
	hmacHash := hmac.New(sha256.New, []byte(hex.EncodeToString(hmacKey)))
	_, err = io.Copy(hmacHash, inputReader)
	if err != nil {
		return false, err
	}
	calculatedHmac := hmacHash.Sum(nil)

	return hmac.Equal(calculatedHmac, expectedHmac), nil
}
