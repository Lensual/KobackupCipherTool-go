package infoxml

import (
	"encoding/xml"
	"errors"
	"os"
)

// InfoXml 结构体用于解析 info.xml
type InfoXml struct {
	XMLName xml.Name `xml:"info.xml"`
	Rows    []Row    `xml:"row"`
}

// Row 表示 info.xml 中的每一行
type Row struct {
	Table   string   `xml:"table,attr"`
	Columns []Column `xml:"column"`
}

// Column 表示 row 中的每一列
type Column struct {
	Name  string `xml:"name,attr"`
	Value *Value `xml:"value"`
}

// Value 表示 column 的值
type Value struct {
	String  string `xml:"String,attr"`
	Integer string `xml:"Integer,attr"`
	Long    string `xml:"Long,attr"`
	Boolean string `xml:"Boolean,attr"`
	Null    string `xml:"Null,attr"`
}

// GetString 获取 String 类型的值
func (v *Value) GetString() string {
	return v.String
}

// GetInteger 获取 Integer 类型的值
func (v *Value) GetInteger() int {
	var n int
	if v.Integer != "" {
		// 简单解析整数
		for _, c := range v.Integer {
			if c >= '0' && c <= '9' {
				n = n*10 + int(c-'0')
			}
		}
	}
	return n
}

// GetLong 获取 Long 类型的值
func (v *Value) GetLong() int64 {
	var n int64
	if v.Long != "" {
		for _, c := range v.Long {
			if c >= '0' && c <= '9' {
				n = n*10 + int64(c-'0')
			}
		}
	}
	return n
}

// GetBoolean 获取 Boolean 类型的值
func (v *Value) GetBoolean() bool {
	return v.Boolean == "true"
}

// IsNull 判断是否为 null
func (v *Value) IsNull() bool {
	return v.Null == "null"
}

// Parse 从 info.xml 文件中解析
func Parse(xmlPath string) (*InfoXml, error) {
	content, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, err
	}

	var infoXml InfoXml
	err = xml.Unmarshal(content, &infoXml)
	if err != nil {
		return nil, err
	}

	return &infoXml, nil
}

// GetRowsByTable 根据表名获取所有行
func (ix *InfoXml) GetRowsByTable(tableName string) []Row {
	var result []Row
	for _, row := range ix.Rows {
		if row.Table == tableName {
			result = append(result, row)
		}
	}
	return result
}

// GetFirstRowByTable 根据表名获取第一行
func (ix *InfoXml) GetFirstRowByTable(tableName string) *Row {
	for _, row := range ix.Rows {
		if row.Table == tableName {
			return &row
		}
	}
	return nil
}

// GetColumnValue 根据列名获取值
func (r *Row) GetColumnValue(columnName string) *Value {
	for _, col := range r.Columns {
		if col.Name == columnName {
			return col.Value
		}
	}
	return nil
}

// GetColumnString 根据列名获取字符串值
func (r *Row) GetColumnString(columnName string) string {
	val := r.GetColumnValue(columnName)
	if val == nil {
		return ""
	}
	return val.GetString()
}

// GetColumnInteger 根据列名获取整数值
func (r *Row) GetColumnInteger(columnName string) int {
	val := r.GetColumnValue(columnName)
	if val == nil {
		return 0
	}
	return val.GetInteger()
}

// GetColumnLong 根据列名获取长整数值
func (r *Row) GetColumnLong(columnName string) int64 {
	val := r.GetColumnValue(columnName)
	if val == nil {
		return 0
	}
	return val.GetLong()
}

// GetColumnBoolean 根据列名获取布尔值
func (r *Row) GetColumnBoolean(columnName string) bool {
	val := r.GetColumnValue(columnName)
	if val == nil {
		return false
	}
	return val.GetBoolean()
}

// GetColumnNull 根据列名判断是否为null
func (r *Row) GetColumnNull(columnName string) bool {
	val := r.GetColumnValue(columnName)
	if val == nil {
		return true
	}
	return val.IsNull()
}

// GetHeaderInfo 获取 HeaderInfo 表的数据
func (ix *InfoXml) GetHeaderInfo() (*HeaderInfo, error) {
	row := ix.GetFirstRowByTable("HeaderInfo")
	if row == nil {
		return nil, errors.New("HeaderInfo not found")
	}
	return &HeaderInfo{
		BackupVersion:    row.GetColumnInteger("backupVersion"),
		SelectDataSize:   row.GetColumnLong("selectDataSize"),
		AutoBackup:       row.GetColumnBoolean("autoBackup"),
		IsbackupBOPD:     row.GetColumnBoolean("isbackupBOPD"),
		Version:          row.GetColumnInteger("version"),
		AutoBackupRandom: row.GetColumnNull("autoBackupRandom"),
		MiniVersion:      row.GetColumnInteger("miniVersion"),
		DateTime:         row.GetColumnLong("dateTime"),
	}, nil
}

// GetBackupFilePhoneInfo 获取 BackupFilePhoneInfo 表的数据
func (ix *InfoXml) GetBackupFilePhoneInfo() (*BackupFilePhoneInfo, error) {
	row := ix.GetFirstRowByTable("BackupFilePhoneInfo")
	if row == nil {
		return nil, errors.New("BackupFilePhoneInfo not found")
	}
	return &BackupFilePhoneInfo{
		ProductManufacturer: row.GetColumnString("productManufacturer"),
		SnHash:              row.GetColumnString("snHash"),
		VersionSdk:          row.GetColumnInteger("versionSdk"),
		DisplayId:           row.GetColumnString("displayId"),
		BOPD_reason:         row.GetColumnString("BOPD_reason"),
		BOPD_running_mode:   row.GetColumnString("BOPD_running_mode"),
		BoardPlatform:       row.GetColumnString("boardPlatform"),
		ProductBrand:        row.GetColumnString("productBrand"),
		ProductModel:        row.GetColumnString("productModel"),
		VersionRelease:      row.GetColumnString("versionRelease"),
		BOPD_info:           row.GetColumnString("BOPD_info"),
		ProductDeviceId:     row.GetColumnString("productDeviceId"),
	}, nil
}

// GetBackupFileVersionInfo 获取 BackupFileVersionInfo 表的数据
func (ix *InfoXml) GetBackupFileVersionInfo() (*BackupFileVersionInfo, error) {
	row := ix.GetFirstRowByTable("BackupFileVersionInfo")
	if row == nil {
		return nil, errors.New("BackupFileVersionInfo not found")
	}
	return &BackupFileVersionInfo{
		DbVersion:         row.GetColumnInteger("dbVersion"),
		SoftVersion:       row.GetColumnInteger("softVersion"),
		BackupVersionName: row.GetColumnString("backupVersionName"),
	}, nil
}

// GetBackupFilesTypeInfo 获取 BackupFilesTypeInfo 表的数据
func (ix *InfoXml) GetBackupFilesTypeInfo() (*BackupFilesTypeInfo, error) {
	row := ix.GetFirstRowByTable("BackupFilesTypeInfo")
	if row == nil {
		return nil, errors.New("BackupFilesTypeInfo not found")
	}
	return &BackupFilesTypeInfo{
		EncryptType:   row.GetColumnInteger("encrypt_type"),
		TypeAttch:     row.GetColumnInteger("type_attch"),
		PromptMsg:     row.GetColumnNull("promptMsg"),
		EPerbackupkey: row.GetColumnNull("e_perbackupkey"),
		PwkeySalt:     row.GetColumnNull("pwkey_salt"),
		Type:          row.GetColumnInteger("type"),
	}, nil
}

// GetBackupFileModuleInfo 获取所有 BackupFileModuleInfo 表的数据
func (ix *InfoXml) GetBackupFileModuleInfo() ([]BackupFileModuleInfo, error) {
	rows := ix.GetRowsByTable("BackupFileModuleInfo")
	if len(rows) == 0 {
		return nil, errors.New("BackupFileModuleInfo not found")
	}

	var result []BackupFileModuleInfo
	for _, row := range rows {
		info := BackupFileModuleInfo{
			DeviceAllLanguages: row.GetColumnNull("deviceAllLanguages"),
			CheckInfoType:      row.GetColumnNull("checkInfoType"),
			DeviceDensityDpi:   row.GetColumnInteger("deviceDensityDpi"),
			Tables:             row.GetColumnNull("tables"),
			CheckMsgV3:         row.GetColumnString("checkMsgV3"),
			IsBundleApp:        row.GetColumnBoolean("isBundleApp"),
			Name:               row.GetColumnString("name"),
			Type:               row.GetColumnInteger("type"),
			SdkSupport:         row.GetColumnInteger("sdkSupport"),
			DeviceCpuArchType:  row.GetColumnNull("deviceCpuArchType"),
			CheckInfo:          row.GetColumnNull("checkInfo"),
			AppSignatures:      row.GetColumnString("appSignatures"),
			ArkBcVersion:       row.GetColumnLong("arkBcVersion"),
			CheckComplexMsgV3:  row.GetColumnNull("checkComplexMsgV3"),
			RecordTotal:        row.GetColumnInteger("recordTotal"),
			IsCopyFileEncrypt:  row.GetColumnBoolean("isCopyFileEncrypt"),
			CopyFilePath:       row.GetColumnNull("copyFilePath"),
			CheckMsg:           row.GetColumnNull("checkMsg"),
			EncMsgV3:           row.GetColumnString("encMsgV3"),
		}
		result = append(result, info)
	}

	return result, nil
}
