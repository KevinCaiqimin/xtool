package cfgexport

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/KevinCaiqimin/go-basic/utils"
)

func tojsWriteNodeToBuf(tbl *ExportTable, node *ExportTreeNode, buf *bytes.Buffer, indent string) {
	stepIndent := getIndent()
	line := ""
	itsRoot := node.FieldRef == nil

	if !itsRoot {
		var k interface{}
		if len(node.Children) == 0 {
			leafVal := node.FieldValue.(map[string]interface{})
			k, _ = leafVal[node.FieldRef.Name]
		} else {
			k = node.FieldValue
		}
		switch node.FieldRef.Typ.Name {
		case FIELD_TYPE_INT, FIELD_TYPE_FLOAT:
			line = fmt.Sprintf("%v%v : {\n", indent, k)
		case FIELD_TYPE_STRING:
			line = fmt.Sprintf("%v\"%v\" : {\n", indent, k)
		default:
			panic("Invalid primary key type")
		}
		buf.WriteString(line)
	}
	if len(node.Children) > 0 {
		newIndent := indent + stepIndent
		if itsRoot {
			newIndent = indent
		}
		for _, k := range node.ChildrenKeySeq {
			child, _ := node.Children[k]
			tojsWriteNodeToBuf(tbl, child, buf, newIndent)
		}
	} else if node.FieldRef != nil {
		//write leaf node
		leafIndent := indent + stepIndent
		leafVal := node.FieldValue.(map[string]interface{})
		for _, field := range tbl.Fields {
			fieldVal, _ := leafVal[field.Name]
			line := ""
			switch field.Typ.Name {
			case FIELD_TYPE_INT, FIELD_TYPE_FLOAT:
				line = fmt.Sprintf("%v\"%v\" : %v,\n", leafIndent, field.Name, fieldVal)
			case FIELD_TYPE_STRING:
				line = fmt.Sprintf("%v\"%v\" : \"%v\",\n", leafIndent, field.Name, fieldVal)
			case FIELD_TYPE_BOOL:
				val := fieldVal.(bool)
				valStr := ""
				if val {
					valStr = "true"
				} else {
					valStr = "false"
				}
				line = fmt.Sprintf("%v\"%v\" : %v,\n", leafIndent, field.Name, valStr)
			case FIELD_TYPE_LUA:
				line = fmt.Sprintf("%v\"%v\" : \"%v\",\n", leafIndent, field.Name, fieldVal)
			case FIELD_TYPE_DIC:
				val := fieldVal.(map[interface{}]interface{})
				line = fmt.Sprintf("%v\"%v\" : {\n", leafIndent, field.Name)
				tmpInde := leafIndent + stepIndent
				for k, v := range val {
					switch field.Typ.Params[0] {
					case FIELD_TYPE_INT, FIELD_TYPE_FLOAT:
						line += fmt.Sprintf("%v%v : ", tmpInde, k)
					case FIELD_TYPE_STRING:
						line += fmt.Sprintf("%v\"%v\" : ", tmpInde, k)
					default:
						panic("invalid dic key type")
					}
					switch field.Typ.Params[1] {
					case FIELD_TYPE_INT, FIELD_TYPE_FLOAT:
						line += fmt.Sprintf("%v,\n", v)
					case FIELD_TYPE_STRING:
						line += fmt.Sprintf("\"%v\",\n", v)
					case FIELD_TYPE_BOOL:
						tmpV := v.(bool)
						tmpStr := ""
						if tmpV {
							tmpStr = "true"
						} else {
							tmpStr = "false"
						}
						line += fmt.Sprintf("\"%v\",\n", tmpStr)
					default:
						panic("invalid dic value type")
					}
				}
				line += fmt.Sprintf("%v},\n", leafIndent)
			case FIELD_TYPE_ARRAY:
				val := fieldVal.([]interface{})
				line = fmt.Sprintf("%v\"%v\" : [", leafIndent, field.Name)
				for _, v := range val {
					switch field.Typ.Params[0] {
					case FIELD_TYPE_INT, FIELD_TYPE_FLOAT:
						line += fmt.Sprintf("%v, ", v)
					case FIELD_TYPE_STRING:
						line += fmt.Sprintf("\"%v\", ", v)
					case FIELD_TYPE_BOOL:
						tmpV := v.(bool)
						tmpStr := ""
						if tmpV {
							tmpStr = "true"
						} else {
							tmpStr = "false"
						}
						line += fmt.Sprintf("\"%v\", ", tmpStr)
					default:
						panic("invalid array value type")
					}
				}
				line += fmt.Sprintf("],\n")
			default:
				panic("Invalid primary key type")
			}
			buf.WriteString(line)
		}
	}
	if !itsRoot {
		buf.WriteString(fmt.Sprintf("%v},\n", indent))
	}
}

func convertToJsStr(file *ExportFile, tblName, indent string, templFile string) (string, error) {
	// linesRequire := 3
	sheetTemplFilePath := templFile

	stepIndent := getIndent()

	sheetTemplBuf, err := ioutil.ReadFile(sheetTemplFilePath)
	if err != nil {
		return "", err
	}
	pureFileName := utils.GetFilePathShortName(file.FileName)

	sheetTempl := string(sheetTemplBuf[:])

	tbl, _ := file.Tables[tblName]
	var buf *bytes.Buffer = &bytes.Buffer{}

	tojsWriteNodeToBuf(tbl, tbl.ContentRoot, buf, indent+stepIndent)
	content := buf.String()

	buf.Reset()
	for _, field := range tbl.Fields {
		buf.WriteString(fmt.Sprintf("%v%v: %v\n", indent+stepIndent, field.Name, field.Comment))
	}
	comment := buf.String()

	now := time.Now()
	exportTime := now.Format("2006-01-02 15:04:05")

	sheetTempl = strings.Replace(sheetTempl, "$fileName", pureFileName, -1)
	sheetTempl = strings.Replace(sheetTempl, "$sheetName", tbl.Name, -1)
	sheetTempl = strings.Replace(sheetTempl, "$exportTime", exportTime, -1)
	sheetTempl = strings.Replace(sheetTempl, "$content", content, -1)
	sheetTempl = strings.Replace(sheetTempl, "$comment", comment, -1)

	return sheetTempl, nil
}

func tojs(file *ExportFile) error {
	fileTemplFilePath := "template/tojs/file_tojs.templ"

	fileTemplBuf, err := ioutil.ReadFile(fileTemplFilePath)
	if err != nil {
		return err
	}
	pureFileName := utils.GetFilePathShortName(file.FileName)

	now := time.Now()
	exportTime := now.Format("2006-01-02 15:04:05")

	fileTempl := string(fileTemplBuf[:])

	fileTempl = strings.Replace(fileTempl, "$fileName", pureFileName, -1)
	fileTempl = strings.Replace(fileTempl, "$exportTime", exportTime, -1)

	indent := getIndent()

	var buf bytes.Buffer
	for _, tbl := range file.Tables {
		tblContent, err := convertToJsStr(file, tbl.Name,
			indent, "template/tojs/sheet_tojs.templ")
		if err != nil {
			return err
		}
		buf.WriteString(tblContent)
		buf.WriteString("\n")
	}
	tablesContent := buf.String()
	fileTempl = strings.Replace(fileTempl, "$tables", tablesContent, -1)

	os.MkdirAll(file.ExportToDir, os.ModePerm)
	err = ioutil.WriteFile(file.ExportToDir+"/"+pureFileName+".js", []byte(fileTempl), os.ModePerm)

	return err
}

func tojsUseSheet(file *ExportFile) error {
	for _, tbl := range file.Tables {
		tblContent, err := convertToJsStr(file, tbl.Name, "", "template/tojs/file_tojs_usesheet.templ")
		if err != nil {
			return err
		}
		os.MkdirAll(file.ExportToDir, os.ModePerm)
		err = ioutil.WriteFile(file.ExportToDir+"/"+tbl.Name+".js", []byte(tblContent), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
