package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/Lensual/KobackupCipherTool-go/internal"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	argPassword := flag.String("password", "", "Decryption password used to generate HMAC key")
	argCheckMsgV3 := flag.String("checkMsgV3", "", "CheckMsgV3 string containing expected HMAC, salt and filename info")
	argInput := flag.String("input", "", "Input file path to verify hash")
	flag.Parse()

	checkMsgV3Items, err := internal.ParseCheckMsgV3(*argCheckMsgV3)
	if err != nil {
		log.Fatalf("ParseCheckMsgV3 Failed: %v", err)
	}

	// find the input file in checkMsgV3
	var checkMsgV3Item internal.CheckMsgV3Item
	for _, item := range checkMsgV3Items {
		if item.FileName == filepath.Base(*argInput) {
			checkMsgV3Item = item
		}
	}

	log.Printf("checkMsgV3Item.ExpectedHmac: %X", checkMsgV3Item.ExpectedHmac)
	log.Printf("checkMsgV3Item.Salt: %X", checkMsgV3Item.Salt)
	log.Printf("checkMsgV3Item.FileName: %s", checkMsgV3Item.FileName)

	inputFile, err := os.OpenFile(*argInput, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("Input File can't open: %v", err)
	}
	defer inputFile.Close()

	fileHash, err := hmacFile(*argPassword, checkMsgV3Item.Salt, inputFile)
	if err != nil {
		log.Fatalf("hmacFile Failed: %v", err)
	}

	log.Printf("File Hash: %X", fileHash)

	if !hmac.Equal(fileHash, checkMsgV3Item.ExpectedHmac) {
		log.Fatalf("Hash Dismatch: %v", err)
	}

	log.Printf("Success")

	os.Exit(0)
}

// hmacFile
func hmacFile(password string, salt []byte, fileReader io.Reader) ([]byte, error) {
	// parse hmacKey
	pbkdf2Key := pbkdf2.Key([]byte(password), salt, 5000, 32, sha256.New)
	hmacKey := []byte(hex.EncodeToString(pbkdf2Key))

	// hmacSha256
	hmacHash := hmac.New(sha256.New, hmacKey)
	_, err := io.Copy(hmacHash, fileReader)
	if err != nil {
		return nil, err
	}
	return hmacHash.Sum(nil), nil
}
