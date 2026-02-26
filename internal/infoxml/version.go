package infoxml

// BackupFileVersionInfo 表示备份版本信息
type BackupFileVersionInfo struct {
	DbVersion         int    // 数据库版本
	SoftVersion       int    // 软件版本
	BackupVersionName string // 备份版本名称
}
