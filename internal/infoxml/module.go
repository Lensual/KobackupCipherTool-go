package infoxml

// BackupFileModuleInfo 表示备份模块信息
type BackupFileModuleInfo struct {
	DeviceAllLanguages bool   // 设备所有语言(null)
	CheckInfoType      bool   // 检查信息类型(null)
	DeviceDensityDpi   int    // 设备DPI密度
	Tables             bool   // 表(null)
	CheckMsgV3         string // 检查消息V3
	IsBundleApp        bool   // 是否捆绑应用
	Name               string // 包名
	Type               int    // 类型
	SdkSupport         int    // SDK支持
	DeviceCpuArchType  bool   // 设备CPU架构类型(null)
	CheckInfo          bool   // 检查信息(null)
	AppSignatures      string // 应用签名
	ArkBcVersion       int64  // ARK BC版本
	CheckComplexMsgV3  bool   // 检查复杂消息V3(null)
	RecordTotal        int    // 记录总数
	IsCopyFileEncrypt  bool   // 是否复制文件加密
	CopyFilePath       bool   // 复制文件路径(null)
	CheckMsg           bool   // 检查消息(null)
	EncMsgV3           string // 加密消息V3
}
