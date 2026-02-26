package infoxml

// BackupFilesTypeInfo 表示备份类型信息
type BackupFilesTypeInfo struct {
	EncryptType   int  // 加密类型
	TypeAttch     int  // 类型附件
	PromptMsg     bool // 提示消息(null)
	EPerbackupkey bool // 加密备份密钥(null)
	PwkeySalt     bool // 密码盐(null)
	Type          int  // 类型
}
