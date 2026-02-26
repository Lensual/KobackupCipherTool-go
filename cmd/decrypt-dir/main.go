package main

import (
	"crypto/sha256"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lensual/KobackupCipherTool-go/internal"
	"github.com/Lensual/KobackupCipherTool-go/internal/utils"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	argPassword := flag.String("password", "", "Decryption password used to generate AES key")
	argInput := flag.String("input", "", "Input directory path")
	flag.Parse()

	// 检查 input 参数
	inputPath := *argInput

	// 使用 os.Stat 获取文件信息
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Fatalf("Failed to stat input path: %v", err)
	}

	// 必须是目录
	if !fileInfo.IsDir() {
		log.Fatalf("Input is not a directory: %s", inputPath)
	}

	// 从 info.xml 解析 encMsgV3
	infoXmlPath := filepath.Join(inputPath, "info.xml")
	encMsgV3Str, err := internal.ParseInfoXml(infoXmlPath)
	if err != nil {
		log.Fatalf("Failed to parse info.xml: %v", err)
	}
	log.Printf("encMsgV3 from info.xml: %s", encMsgV3Str)

	// 32 bytes key is aes-256
	encMsgV3, err := internal.ParseEncMsgV3(*argPassword, encMsgV3Str)
	if err != nil {
		log.Fatalf("ParseEncMsgV3 Failed: %v", err)
	}

	log.Printf("encMsgV3.Salt: %X", encMsgV3.Salt)
	log.Printf("encMsgV3.Iv: %X", encMsgV3.Iv)

	key := pbkdf2.Key([]byte(*argPassword), encMsgV3.Salt, 5000, 32, sha256.New)
	log.Printf("key: %X", key)

	// 计算输出目录路径：在原目录名后添加 "_decrypted"
	outputDir := inputPath + "_decrypted"

	// 创建输出目录
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// 尝试从 backupinfo.ini 获取包名
	backupInfoIniPath := filepath.Join(inputPath, "backupinfo.ini")
	var targetTarDirs []string
	if _, err := os.Stat(backupInfoIniPath); err != nil {
		log.Printf("Failed to find backupinfo.ini: %v", err)
		os.Exit(1)
	}

	// backupinfo.ini 存在，解析包名
	backupInfo, err := internal.ParseBackupInfo(backupInfoIniPath)
	if err != nil {
		log.Fatalf("Failed to parse backupinfo.ini: %v", err)
	}
	// 遍历所有包名，构建 appDataTar 目录路径
	for _, packageName := range backupInfo.PackageNames {
		targetTarDir := filepath.Join(inputPath, packageName+"_appDataTar")
		if _, err := os.Stat(targetTarDir); err == nil {
			targetTarDirs = append(targetTarDirs, targetTarDir)
			log.Printf("Found package name: %s, targeting directory: %s", packageName, targetTarDir)
		} else {
			log.Printf("Target directory %s does not exist, skipping...", targetTarDir)
		}
	}

	// 使用 filepath.WalkDir 遍历所有目标目录
	for _, targetTarDir := range targetTarDirs {
		log.Printf("Walking directory: %s", targetTarDir)
		err = filepath.WalkDir(targetTarDir, func(path string, d os.DirEntry, err error) error {
			// 忽略目录遍历中的错误，继续处理其他文件
			if err != nil {
				log.Printf("Walk error on %s: %v, skipping...", path, err)
				return nil
			}

			// 跳过目录本身，只处理文件
			if d.IsDir() {
				return nil
			}

			// 只处理 .tar 文件
			if !strings.HasSuffix(d.Name(), ".tar") {
				return nil
			}

			// 计算相对路径（相对于原始 inputPath）
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
			log.Fatalf("WalkDir Failed: %v", err)
		}
	}

	log.Printf("Folder decryption completed")
	os.Exit(0)
}
