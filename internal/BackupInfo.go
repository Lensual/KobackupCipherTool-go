package internal

import (
	"os"
	"strconv"
	"strings"
	"unicode/utf16"
)

// ReadUTF16LEFile 读取 UTF-16LE 编码的文件
func ReadUTF16LEFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// 检查 UTF-16LE BOM
	if len(content) >= 2 && content[0] == 0xFF && content[1] == 0xFE {
		// UTF-16LE BOM
		utf16Content := make([]uint16, (len(content)-2)/2)
		for i := 0; i < len(utf16Content); i++ {
			utf16Content[i] = uint16(content[i*2+2]) | uint16(content[i*2+3])<<8
		}
		return string(utf16.Decode(utf16Content)), nil
	}

	// 尝试 UTF-16LE 解码
	if len(content) >= 2 && len(content)%2 == 0 {
		// 检查是否可能是 UTF-16LE（通常包含大量空字节）
		nullCount := 0
		for i := 1; i < len(content); i += 2 {
			if content[i] == 0 {
				nullCount++
			}
		}
		// 如果超过一半的偶数位置是0，认为是 UTF-16LE
		if nullCount > len(content)/4 {
			utf16Content := make([]uint16, len(content)/2)
			for i := 0; i < len(utf16Content); i++ {
				utf16Content[i] = uint16(content[i*2]) | uint16(content[i*2+1])<<8
			}
			result := string(utf16.Decode(utf16Content))
			if !strings.Contains(result, "\ufffd") {
				return result, nil
			}
		}
	}

	// 回退到直接返回原始内容（可能是普通文本或 UTF-8）
	return string(content), nil
}

// BackupInfo 存储从 backupinfo.ini 解析出的信息
type BackupInfo struct {
	PackageNames []string  // 包名列表
	Apps         []AppInfo // 应用信息列表
}

// AppInfo 存储单个应用的信息
type AppInfo struct {
	PackageName string
	AppName     string
	VersionCode int
	VersionName string
	IsHaveDb    int
	IsHapApp    int
	ApkSize     int64
	DbSize      int64
}

// ParseBackupInfo 从 backupinfo.ini 文件中解析包名
func ParseBackupInfo(iniPath string) (*BackupInfo, error) {
	content, err := ReadUTF16LEFile(iniPath)
	if err != nil {
		return nil, err
	}

	backupInfo := &BackupInfo{}
	lines := strings.Split(content, "\n")

	var currentApp string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 解析 app_info 字段（可能包含多个包名，用逗号分隔）
		// 格式: app_info=com.tencent.mm,
		if strings.HasPrefix(line, "app_info=") {
			value := strings.TrimPrefix(line, "app_info=")
			// 使用逗号分割多个包名
			parts := strings.Split(value, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part != "" {
					backupInfo.PackageNames = append(backupInfo.PackageNames, part)
					currentApp = part
				}
			}
			break
		}
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 解析应用详细信息（需要先确定当前是哪个应用）
		// 格式: [com.tencent.mm]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentApp = strings.Trim(line, "[]")
			if currentApp == "headerinfo" {
				//skip
				currentApp = ""
			}
			if currentApp == "overview" {
				//skip
				currentApp = ""
			}
			continue
		}

		// 解析应用字段
		if currentApp != "" {
			if strings.HasPrefix(line, "app_name=") {
				value := strings.TrimPrefix(line, "app_name=")
				backupInfo.Apps = append(backupInfo.Apps, AppInfo{PackageName: currentApp, AppName: value})
			} else if strings.HasPrefix(line, "version_code=") {
				value := strings.TrimPrefix(line, "version_code=")
				if v, err := strconv.Atoi(value); err == nil {
					for i := range backupInfo.Apps {
						if backupInfo.Apps[i].PackageName == currentApp {
							backupInfo.Apps[i].VersionCode = v
						}
					}
				}
			} else if strings.HasPrefix(line, "version_name=") {
				value := strings.TrimPrefix(line, "version_name=")
				for i := range backupInfo.Apps {
					if backupInfo.Apps[i].PackageName == currentApp {
						backupInfo.Apps[i].VersionName = value
					}
				}
			} else if strings.HasPrefix(line, "is_have_db=") {
				value := strings.TrimPrefix(line, "is_have_db=")
				if v, err := strconv.Atoi(value); err == nil {
					for i := range backupInfo.Apps {
						if backupInfo.Apps[i].PackageName == currentApp {
							backupInfo.Apps[i].IsHaveDb = v
						}
					}
				}
			} else if strings.HasPrefix(line, "is_hap_app=") {
				value := strings.TrimPrefix(line, "is_hap_app=")
				if v, err := strconv.Atoi(value); err == nil {
					for i := range backupInfo.Apps {
						if backupInfo.Apps[i].PackageName == currentApp {
							backupInfo.Apps[i].IsHapApp = v
						}
					}
				}
			} else if strings.HasPrefix(line, "apk_size=") {
				value := strings.TrimPrefix(line, "apk_size=")
				if v, err := strconv.ParseInt(value, 10, 64); err == nil {
					for i := range backupInfo.Apps {
						if backupInfo.Apps[i].PackageName == currentApp {
							backupInfo.Apps[i].ApkSize = v
						}
					}
				}
			} else if strings.HasPrefix(line, "db_size=") {
				value := strings.TrimPrefix(line, "db_size=")
				if v, err := strconv.ParseInt(value, 10, 64); err == nil {
					for i := range backupInfo.Apps {
						if backupInfo.Apps[i].PackageName == currentApp {
							backupInfo.Apps[i].DbSize = v
						}
					}
				}
			}
		}
	}

	if len(backupInfo.PackageNames) == 0 {
		return nil, nil
	}

	return backupInfo, nil
}
