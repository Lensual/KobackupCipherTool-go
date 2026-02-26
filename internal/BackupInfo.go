package internal

import (
	"os"
	"strconv"
	"strings"
)

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
	content, err := os.ReadFile(iniPath)
	if err != nil {
		return nil, err
	}

	backupInfo := &BackupInfo{}
	lines := strings.Split(string(content), "\n")

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
