package infoxml

// BackupFilePhoneInfo 表示手机设备信息
type BackupFilePhoneInfo struct {
	ProductManufacturer string // 产品制造商
	SnHash              string // 序列号哈希
	VersionSdk          int    // SDK版本
	DisplayId           string // 显示ID
	BOPD_reason         string // BOPD原因
	BOPD_running_mode   string // BOPD运行模式
	BoardPlatform       string // 主板平台
	ProductBrand        string // 产品品牌
	ProductModel        string // 产品型号
	VersionRelease      string // 系统版本
	BOPD_info           string // BOPD信息
	ProductDeviceId     string // 产品设备ID
}
