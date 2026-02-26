package infoxml

// HeaderInfo 表示备份文件的头部信息
type HeaderInfo struct {
	BackupVersion    int   // 备份版本
	SelectDataSize   int64 // 选择的数据大小
	AutoBackup       bool  // 是否自动备份
	IsbackupBOPD     bool  // 是否备份BOPD
	Version          int   // 版本号
	AutoBackupRandom bool  // 自动备份随机值(null)
	MiniVersion      int   // 最小版本
	DateTime         int64 // 日期时间
}
