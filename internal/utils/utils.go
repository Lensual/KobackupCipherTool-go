package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"
)

type ALGO int

const (
	ALGO_AES_CTR ALGO = iota
	ALGO_AES_GCM
)

func DecryptFile(in string, out string, key []byte, iv []byte, algo ALGO) error {
	inFile, err := os.OpenFile(in, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer outFile.Close()

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	switch algo {
	case ALGO_AES_CTR:
	case ALGO_AES_GCM:
		return GcmDecrypt(inFile, outFile, blockCipher, key, iv)
	}

	return nil
}

func CtrDecrypt(in io.Reader, out io.Writer, blockCipher cipher.Block, key []byte, iv []byte) error {
	aesCtr := cipher.NewCTR(blockCipher, iv)
	cipherStreamReader := cipher.StreamReader{
		S: aesCtr,
		R: in,
	}

	_, err := io.Copy(out, cipherStreamReader)
	if err != nil {
		return err
	}

	return nil
}

func GcmDecrypt(in io.Reader, out io.Writer, blockCipher cipher.Block, key []byte, iv []byte) error {
	// aesGcm, err := cipher.NewGCM(blockCipher)
	aesGcm, err := cipher.NewGCMWithNonceSize(blockCipher, 16) // compatible kobackup
	if err != nil {
		return err
	}

	// golang is not support streaming AEAD
	inBuf, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	outBuf, err := aesGcm.Open(nil, iv, inBuf, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, bytes.NewReader(outBuf))
	if err != nil {
		return err
	}

	return nil
}
