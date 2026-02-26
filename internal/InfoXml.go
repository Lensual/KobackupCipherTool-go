package internal

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

// ParseInfoXml 从 info.xml 文件中解析 encMsgV3
func ParseInfoXml(xmlPath string) (string, error) {
	content, err := os.ReadFile(xmlPath)
	if err != nil {
		return "", err
	}

	var infoXml InfoXml
	err = xml.Unmarshal(content, &infoXml)
	if err != nil {
		return "", err
	}

	// 遍历查找 BackupFileModuleInfo 表中的 encMsgV3
	for _, row := range infoXml.Rows {
		if row.Table == "BackupFileModuleInfo" {
			for _, column := range row.Columns {
				if column.Name == "encMsgV3" && column.Value != nil {
					if column.Value.String != "" {
						return column.Value.String, nil
					}
				}
			}
		}
	}

	return "", errors.New("encMsgV3 not found in info.xml")
}
