package main

import (
	"crypto/sha256"
	"flag"
	"log"
	"os"
	"path/filepath"

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

	// 检查 input 参数
	inputPath := *argInput
	outputPath := *argOutput

	// 使用 os.Stat 获取文件信息
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Printf("Failed to stat input path: %v", err)
		os.Exit(1)
	}

	// 判断是文件还是目录
	if fileInfo.IsDir() {
		// 如果是目录，遍历目录中所有文件
		log.Printf("Input is a directory, starting folder traversal decryption...")

		// 计算输出目录路径：在原目录名后添加 "_decrypted"
		outputDir := inputPath + "_decrypted"

		// 创建输出目录
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			log.Printf("Failed to create output directory: %v", err)
			os.Exit(1)
		}

		// 使用 filepath.WalkDir 遍历目录
		err = filepath.WalkDir(inputPath, func(path string, d os.DirEntry, err error) error {
			// 忽略目录遍历中的错误，继续处理其他文件
			if err != nil {
				log.Printf("Walk error on %s: %v, skipping...", path, err)
				return nil
			}

			// 跳过目录本身，只处理文件
			if d.IsDir() {
				return nil
			}

			// 计算相对路径
			relPath, err := filepath.Rel(inputPath, path)
			if err != nil {
				log.Printf("Failed to get relative path for %s: %v, skipping...", path, err)
				return nil
			}

			// 构建输出文件路径
			outputFilePath := filepath.Join(outputDir, relPath)

			// 确保输出文件的父目录存在
			outputDirPath := filepath.Dir(outputFilePath)
			err = os.MkdirAll(outputDirPath, 0755)
			if err != nil {
				log.Printf("Failed to create output subdirectory %s: %v, skipping...", outputDirPath, err)
				return nil
			}

			// 解密文件
			log.Printf("Decrypting: %s -> %s", path, outputFilePath)
			err = utils.DecryptFile(path, outputFilePath, key, encMsgV3.Iv, utils.ALGO_AES_GCM)
			if err != nil {
				log.Printf("DecryptFile Failed for %s: %v, skipping...", path, err)
				return nil
			}

			log.Printf("Success: %s", outputFilePath)
			return nil
		})

		if err != nil {
			log.Printf("WalkDir Failed: %v", err)
			os.Exit(1)
		}

		log.Printf("Folder decryption completed")
	} else {
		// 如果是文件，直接解密（现有逻辑）
		err = utils.DecryptFile(inputPath, outputPath, key, encMsgV3.Iv, utils.ALGO_AES_GCM)
		if err != nil {
			log.Printf("DecryptFile Failed: %v", err)
			os.Exit(1)
		}

		log.Printf("Success")
	}

	os.Exit(0)
}
