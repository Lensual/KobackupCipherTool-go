package infoxml_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Lensual/KobackupCipherTool-go/internal/infoxml"
)

// TestParseFromWorkspace 从工作区目录读取 info.xml 并打印解析结果
func TestParseFromWorkspace(t *testing.T) {
	// 从工作区目录读取 info.xml（相对于 internal/infoxml 目录）
	xmlPath := "../../info.xml"

	// 检查文件是否存在
	if _, err := os.Stat(xmlPath); os.IsNotExist(err) {
		t.Skipf("info.xml not found at path: %s", xmlPath)
		return
	}

	// 解析 info.xml
	infoXml, err := infoxml.Parse(xmlPath)
	if err != nil {
		t.Fatalf("Failed to parse info.xml: %v", err)
	}

	// 打印解析结果
	printParseResults(infoXml)
}

// printParseResults 打印所有解析结果
func printParseResults(infoXml *infoxml.InfoXml) {
	fmt.Println("=== Info.xml 解析结果 ===")
	fmt.Printf("共解析到 %d 个表\n\n", len(infoXml.Rows))

	// 打印 HeaderInfo
	if headerInfo, err := infoXml.GetHeaderInfo(); err == nil {
		fmt.Println("--- HeaderInfo ---")
		fmt.Printf("  BackupVersion:    %d\n", headerInfo.BackupVersion)
		fmt.Printf("  SelectDataSize:   %d\n", headerInfo.SelectDataSize)
		fmt.Printf("  AutoBackup:       %v\n", headerInfo.AutoBackup)
		fmt.Printf("  IsbackupBOPD:     %v\n", headerInfo.IsbackupBOPD)
		fmt.Printf("  Version:          %d\n", headerInfo.Version)
		fmt.Printf("  AutoBackupRandom: %v\n", headerInfo.AutoBackupRandom)
		fmt.Printf("  MiniVersion:      %d\n", headerInfo.MiniVersion)
		fmt.Printf("  DateTime:         %d\n", headerInfo.DateTime)
		fmt.Println()
	}

	// 打印 BackupFilePhoneInfo
	if phoneInfo, err := infoXml.GetBackupFilePhoneInfo(); err == nil {
		fmt.Println("--- BackupFilePhoneInfo ---")
		fmt.Printf("  ProductManufacturer: %s\n", phoneInfo.ProductManufacturer)
		fmt.Printf("  SnHash:              %s\n", phoneInfo.SnHash)
		fmt.Printf("  VersionSdk:          %d\n", phoneInfo.VersionSdk)
		fmt.Printf("  DisplayId:           %s\n", phoneInfo.DisplayId)
		fmt.Printf("  BOPD_reason:        %s\n", phoneInfo.BOPD_reason)
		fmt.Printf("  BOPD_running_mode:  %s\n", phoneInfo.BOPD_running_mode)
		fmt.Printf("  BoardPlatform:      %s\n", phoneInfo.BoardPlatform)
		fmt.Printf("  ProductBrand:        %s\n", phoneInfo.ProductBrand)
		fmt.Printf("  ProductModel:        %s\n", phoneInfo.ProductModel)
		fmt.Printf("  VersionRelease:      %s\n", phoneInfo.VersionRelease)
		fmt.Printf("  BOPD_info:           %s\n", phoneInfo.BOPD_info)
		fmt.Printf("  ProductDeviceId:     %s\n", phoneInfo.ProductDeviceId)
		fmt.Println()
	}

	// 打印 BackupFileVersionInfo
	if versionInfo, err := infoXml.GetBackupFileVersionInfo(); err == nil {
		fmt.Println("--- BackupFileVersionInfo ---")
		fmt.Printf("  DbVersion:         %d\n", versionInfo.DbVersion)
		fmt.Printf("  SoftVersion:       %d\n", versionInfo.SoftVersion)
		fmt.Printf("  BackupVersionName: %s\n", versionInfo.BackupVersionName)
		fmt.Println()
	}

	// 打印 BackupFilesTypeInfo
	if typeInfo, err := infoXml.GetBackupFilesTypeInfo(); err == nil {
		fmt.Println("--- BackupFilesTypeInfo ---")
		fmt.Printf("  EncryptType:   %d\n", typeInfo.EncryptType)
		fmt.Printf("  TypeAttch:     %d\n", typeInfo.TypeAttch)
		fmt.Printf("  PromptMsg:     %v\n", typeInfo.PromptMsg)
		fmt.Printf("  EPerbackupkey: %v\n", typeInfo.EPerbackupkey)
		fmt.Printf("  PwkeySalt:     %v\n", typeInfo.PwkeySalt)
		fmt.Printf("  Type:          %d\n", typeInfo.Type)
		fmt.Println()
	}

	// 打印 BackupFileModuleInfo
	if moduleInfos, err := infoXml.GetBackupFileModuleInfo(); err == nil {
		fmt.Println("--- BackupFileModuleInfo ---")
		for i, module := range moduleInfos {
			fmt.Printf("  [%d] Name: %s\n", i, module.Name)
			fmt.Printf("      Type:               %d\n", module.Type)
			fmt.Printf("      SdkSupport:         %d\n", module.SdkSupport)
			fmt.Printf("      IsBundleApp:        %v\n", module.IsBundleApp)
			fmt.Printf("      IsCopyFileEncrypt:  %v\n", module.IsCopyFileEncrypt)
			fmt.Printf("      EncMsgV3:           %s\n", module.EncMsgV3)
			fmt.Printf("      CheckMsgV3:         %s\n", module.CheckMsgV3[:min(50, len(module.CheckMsgV3))]+"...")
			fmt.Printf("      AppSignatures:      %s\n", module.AppSignatures[:min(50, len(module.AppSignatures))]+"...")
			fmt.Println()
		}
	}
}
