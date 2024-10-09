package main

import (
	"crypto/sha256"
	"flag"
	"log"
	"os"

	"github.com/Lensual/KobackupCipherTool-go/internal"
	"github.com/Lensual/KobackupCipherTool-go/internal/utils"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	argPassword := flag.String("password", "", "")
	argEncMsgV3 := flag.String("encMsgV3", "", "")
	argInput := flag.String("input", "", "")
	argOutput := flag.String("output", "", "")
	flag.Parse()

	// 32 bytes key is aes-256
	encMsgV3, err := internal.ParseEncMsgV3(*argPassword, *argEncMsgV3)
	if err != nil {
		log.Printf("ParseEncMsgV3 Failed: %v", err)
		os.Exit(1)
	}

	log.Printf("encMsgV3.Salt: %X", encMsgV3.Salt)
	log.Printf("encMsgV3.Iv: %X", encMsgV3.Iv)

	key := pbkdf2.Key([]byte(*argPassword), encMsgV3.Salt, 5000, 32, sha256.New)
	log.Printf("key: %X", key)

	err = utils.DecryptFile(*argInput, *argOutput, key, encMsgV3.Iv, utils.ALGO_AES_GCM)
	if err != nil {
		log.Printf("DecryptFile Failed: %v", err)
		os.Exit(1)
	}

	log.Printf("Success")

	os.Exit(0)
}
