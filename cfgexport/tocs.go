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

type CSVars struct {
	ClassName     string
	ClassItemName string
	DataVarName   string
}

var csVars *CSVars

func tocsGenStruct(file *ExportFile, tblName string, indent string) (string, error) {
	f1 := "template/tocs/conf_class_item.templ"
	f2 := "template/tocs/conf_class_item_field.templ"
	f3 := "template/tocs/conf_class_item_value.templ"
	//类模板
	structTempl, e1 := utils.ReadStringFromFile(f1)
	//类字段
	fieldTempl, e2 := utils.ReadStringFromFile(f2)
	//构造函数赋值
	valueTempl, e3 := utils.ReadStringFromFile(f3)
	if e1 != nil {
		return "", fmt.Errorf("Read file %v error: %v", f1, e1)
	}
	if e2 != nil {
		return "", fmt.Errorf("Read file %v error: %v", f2, e2)
	}
	if e3 != nil {
		return "", fmt.Errorf("Read file %v error: %v", f3, e3)
	}
	structTempl = utils.FillIndent(structTempl, indent)
	fieldIndent := utils.GetIndentOf(structTempl, "$FIELDS")
	fieldTempl = utils.FillIndent(fieldTempl, fieldIndent)
	valueIndent := utils.GetIndentOf(structTempl, "$CONSTRUCT_VALUES")
	valueTempl = utils.FillIndent(valueTempl, valueIndent)

	tbl, _ := file.Tables[tblName]

	fieldsBuf := &bytes.Buffer{}
	constructParmBuf := &bytes.Buffer{}
	valuesBuf := &bytes.Buffer{}

	for i := 0; i < len(tbl.Fields); i++ {
		field := tbl.Fields[i]
		fieldType := ""
		switch field.Typ.Name {
		case FIELD_TYPE_INT:
			fieldType = "int"
		case FIELD_TYPE_FLOAT:
			fieldType = "float"
		case FIELD_TYPE_STRING:
			fieldType = "string"
		case FIELD_TYPE_BOOL:
			fieldType = "bool"
		case FIELD_TYPE_LUA:
			fieldType = "string"
		case FIELD_TYPE_DIC:
			switch field.Typ.Params[0] {
			case FIELD_TYPE_INT:
				fieldType = "Dictionary<int, "
			case FIELD_TYPE_FLOAT:
				fieldType = "Dictionary<float, "
			case FIELD_TYPE_STRING:
				fieldType = "Dictionary<string, "
			case FIELD_TYPE_BOOL:
				fieldType = "Dictionary<bool, "
			default:
				panic("invalid dic key type")
			}
			switch field.Typ.Params[1] {
			case FIELD_TYPE_INT:
				fieldType += "int>"
			case FIELD_TYPE_FLOAT:
				fieldType += "float>"
			case FIELD_TYPE_STRING:
				fieldType += "string>"
			case FIELD_TYPE_BOOL:
				fieldType += "bool>"
			default:
				panic("invalid dic value type")
			}
		case FIELD_TYPE_ARRAY:
			switch field.Typ.Params[0] {
			case FIELD_TYPE_INT:
				fieldType = "int[]"
			case FIELD_TYPE_FLOAT:
				fieldType = "float[]"
			case FIELD_TYPE_STRING:
				fieldType = "string[]"
			case FIELD_TYPE_BOOL:
				fieldType = "bool[]"
			default:
				panic("invalid array key type")
			}
		default:
			panic("Invalid primary key type")
		}
		fieldStat := fieldTempl
		fieldStat = strings.Replace(fieldStat, "$COMMENT", field.Comment, -1)
		fieldStat = strings.Replace(fieldStat, "$FIELD_TYPE", fieldType, -1)
		fieldStat = strings.Replace(fieldStat, "$FIELD_NAME", field.Name, -1)

		fieldsBuf.WriteString(fieldStat)

		parmStr := ""
		if i == 0 {
			parmStr = fmt.Sprintf("%v %v", fieldType, field.Name)
		} else {
			parmStr = fmt.Sprintf(",%v %v", fieldType, field.Name)
		}
		constructParmBuf.WriteString(parmStr)

		valueStr := valueTempl
		valueStr = strings.Replace(valueStr, "$FIELD_TYPE", field.Comment, -1)
		valueStr = strings.Replace(valueStr, "$FIELD_NAME", fieldType, -1)
		valuesBuf.WriteString(valueStr)
	}
	structStr := structTempl
	structStr = strings.Replace(structStr, "$CLASS_ITEM_NAME", csVars.ClassItemName, -1)
	structStr = strings.Replace(structStr, "$CLASS_NAME", csVars.ClassName, -1)
	structStr = strings.Replace(structStr, "$FIELDS", utils.UnindentLines(fieldsBuf.String(), 1), -1)
	structStr = strings.Replace(structStr, "$CONSTRUCT_PARMS", constructParmBuf.String(), -1)
	structStr = strings.Replace(structStr, "$CONSTRUCT_VALUES", valuesBuf.String(), -1)

	return structStr, nil
}

func getFieldCSType(field *Field) string {
	fieldType := ""
	switch field.Typ.Name {
	case FIELD_TYPE_INT:
		fieldType = "int"
	case FIELD_TYPE_FLOAT:
		fieldType = "float"
	case FIELD_TYPE_STRING:
		fieldType = "string"
	case FIELD_TYPE_BOOL:
		fieldType = "bool"
	case FIELD_TYPE_LUA:
		fieldType = "string"
	case FIELD_TYPE_DIC:
		switch field.Typ.Params[0] {
		case FIELD_TYPE_INT:
			fieldType = "Dictionary<int, "
		case FIELD_TYPE_FLOAT:
			fieldType = "Dictionary<float, "
		case FIELD_TYPE_STRING:
			fieldType = "Dictionary<string, "
		case FIELD_TYPE_BOOL:
			fieldType = "Dictionary<bool, "
		default:
			panic("invalid dic key type")
		}
		switch field.Typ.Params[1] {
		case FIELD_TYPE_INT:
			fieldType += "int>"
		case FIELD_TYPE_FLOAT:
			fieldType += "float>"
		case FIELD_TYPE_STRING:
			fieldType += "string>"
		case FIELD_TYPE_BOOL:
			fieldType += "bool>"
		default:
			panic("invalid dic value type")
		}
	case FIELD_TYPE_ARRAY:
		switch field.Typ.Params[0] {
		case FIELD_TYPE_INT:
			fieldType = "int[]"
		case FIELD_TYPE_FLOAT:
			fieldType = "float[]"
		case FIELD_TYPE_STRING:
			fieldType = "string[]"
		case FIELD_TYPE_BOOL:
			fieldType = "bool[]"
		default:
			panic("invalid array key type")
		}
	default:
		panic("Invalid primary key type")
	}
	return fieldType
}

func genNodeCSDicType(tbl *ExportTable, targetField *Field) string {
	itsRoot := targetField == nil
	buf := &bytes.Buffer{}
	if targetField != nil && !targetField.IsPrimary {
		return ""
	}
	priNum := 0
	for i := 0; i < len(tbl.Fields); i++ {
		field := tbl.Fields[i]
		if !field.IsPrimary {
			continue
		}
		priNum++
	}
	if itsRoot {
		leftPriNum := priNum
		buf.WriteString("Dictionary<")
		for i := 0; i < len(tbl.Fields); i++ {
			field := tbl.Fields[i]
			if !field.IsPrimary {
				break
			}
			leftPriNum--
			// write key type
			buf.WriteString(getFieldCSType(field) + ",")
			// write value type
			if leftPriNum <= 0 {
				buf.WriteString(csVars.ClassItemName)
			} else {
				buf.WriteString("Dictionary<")
			}
		}
		for i := 0; i < priNum; i++ {
			buf.WriteString(">")
		}
		return buf.String()
	} else {
		leftPriNum := priNum
		startGen := false
		deeps := 0
		for i := 0; i < len(tbl.Fields); i++ {
			field := tbl.Fields[i]
			if !field.IsPrimary {
				break
			}
			leftPriNum--
			if startGen == false && field.Name == targetField.Name {
				startGen = true
				if leftPriNum <= 0 {
					buf.WriteString(csVars.ClassItemName)
					break
				} else {
					buf.WriteString("Dictionary<")
				}
				continue
			}
			if !startGen {
				continue
			}
			deeps++
			// write key type
			buf.WriteString(getFieldCSType(field) + ",")
			// write value type
			if leftPriNum <= 0 {
				buf.WriteString(csVars.ClassItemName)
			} else {
				buf.WriteString("Dictionary<")
			}
		}
		for i := 0; i < deeps; i++ {
			buf.WriteString(">")
		}
		return buf.String()
	}
}

func genGetDataFuncs(tbl *ExportTable, indent string) (string, error) {
	f1 := "template/tocs/conf_get_func.templ"
	f2 := "template/tocs/conf_get_func_loopdic.templ"
	funcTempl, e1 := utils.ReadStringFromFile(f1)
	loopDicTempl, e2 := utils.ReadStringFromFile(f2)
	if e1 != nil {
		return "", fmt.Errorf("Read file %v error: %v", f1, e1)
	}
	if e2 != nil {
		return "", fmt.Errorf("Read file %v error: %v", f2, e2)
	}
	funcTempl = utils.FillIndent(funcTempl, indent)
	loopIndent := utils.GetIndentOf(funcTempl, "$LOOP_GET")
	loopDicTempl = utils.FillIndent(loopDicTempl, loopIndent)

	buf := &bytes.Buffer{}

	loopDicCheck := ""
	priNO := 0

	loopStat := ""
	keysStat := ""

	for i := 0; i < len(tbl.Fields); i++ {
		field := tbl.Fields[i]
		if !field.IsPrimary {
			continue
		}
		priNO++

		loopDicCheck = loopDicTempl
		key := field.Name
		preKey := ""
		if priNO > 1 {
			preKey = tbl.Fields[i-1].Name
		}
		retType := genNodeCSDicType(tbl, field)
		dicName := fmt.Sprintf("dic%v", priNO)
		fieldType := getFieldCSType(field)

		loopDicCheck = strings.Replace(loopDicCheck, "$DIC_NAME", dicName, -1)
		if priNO == 1 {
			loopDicCheck = strings.Replace(loopDicCheck, "$PRE_DIC", csVars.DataVarName, -1)
		} else {
			loopDicCheck = strings.Replace(loopDicCheck, "$PRE_DIC", fmt.Sprintf("dic%d[%v]", priNO-1, preKey), -1)
		}
		loopDicCheck = strings.Replace(loopDicCheck, "$KEY_NAME", key, -1)

		if priNO == 1 {
			keysStat += fmt.Sprintf("%v %v", fieldType, key)
		} else {
			keysStat += fmt.Sprintf(", %v %v", fieldType, key)
		}
		loopStat += loopDicCheck

		funcStat := funcTempl
		funcStat = strings.Replace(funcStat, "$RETURN_TYPE", retType, -1)
		funcStat = strings.Replace(funcStat, "$KEYS", keysStat, -1)
		funcStat = strings.Replace(funcStat, "$LOOP_GET", utils.UnindentLines(loopStat, 1), -1)
		funcStat = strings.Replace(funcStat, "$DIC_NAME", dicName, -1)
		funcStat = strings.Replace(funcStat, "$KEY_NAME", key, -1)

		buf.WriteString(funcStat)
	}

	return buf.String(), nil
}

func correctString(str string) string {
	// str = strings.Replace(str, "\"", "\"", -1)
	// str = strings.Replace(str, "'", "\\'", -1)
	return str
}

func tocsWriteNodeToBuf(tbl *ExportTable, node *ExportTreeNode, buf *bytes.Buffer, indent string) {
	stepIndent := getIndent()
	itsRoot := node.FieldRef == nil

	if !itsRoot {
		buf.WriteString(fmt.Sprintf("%v{", indent))
		var keyVal interface{}
		if len(node.Children) == 0 {
			leafVal := node.FieldValue.(map[string]interface{})
			keyVal, _ = leafVal[node.FieldRef.Name]
		} else {
			keyVal = node.FieldValue
		}
		switch node.FieldRef.Typ.Name {
		case FIELD_TYPE_INT, FIELD_TYPE_FLOAT:
			buf.WriteString(fmt.Sprintf("%v, ", keyVal))
		case FIELD_TYPE_STRING:
			buf.WriteString(fmt.Sprintf("\"%v\", ", correctString(keyVal.(string))))
		default:
			panic("Invalid primary key type")
		}
	}
	dicType := genNodeCSDicType(tbl, node.FieldRef)
	buf.WriteString(fmt.Sprintf("new %v()\n", dicType))
	buf.WriteString(fmt.Sprintf("%v{\n", indent+stepIndent))

	if len(node.Children) > 0 {
		newIndent := indent + stepIndent + stepIndent
		for _, k := range node.ChildrenKeySeq {
			child, _ := node.Children[k]
			tocsWriteNodeToBuf(tbl, child, buf, newIndent)
		}
	} else if node.FieldRef != nil {
		//write leaf node
		leafIndent := indent + stepIndent + stepIndent
		leafVal := node.FieldValue.(map[string]interface{})
		for i := 0; i < len(tbl.Fields); i++ {
			field := tbl.Fields[i]
			fieldVal, _ := leafVal[field.Name]
			line := ""
			switch field.Typ.Name {
			case FIELD_TYPE_INT, FIELD_TYPE_FLOAT:
				line = fmt.Sprintf("%v%v = %v,\n", leafIndent, field.Name, fieldVal)
			case FIELD_TYPE_STRING:
				line = fmt.Sprintf("%v%v = \"%v\",\n", leafIndent, field.Name, correctString(fieldVal.(string)))
			case FIELD_TYPE_BOOL:
				val := fieldVal.(bool)
				valStr := ""
				if val {
					valStr = "true"
				} else {
					valStr = "false"
				}
				line = fmt.Sprintf("%v%v = %v,\n", leafIndent, field.Name, valStr)
			case FIELD_TYPE_LUA:
				line = fmt.Sprintf("%v%v = \"%v\",\n", leafIndent, field.Name, correctString(fieldVal.(string)))
			case FIELD_TYPE_DIC:
				val := fieldVal.(map[interface{}]interface{})
				dicDataType := "Dictionary<"
				switch field.Typ.Params[0] {
				case FIELD_TYPE_INT:
					dicDataType += "int, "
				case FIELD_TYPE_FLOAT:
					dicDataType += "float, "
				case FIELD_TYPE_STRING:
					dicDataType += "string, "
				default:
					panic("invalid dic key type")
				}
				switch field.Typ.Params[1] {
				case FIELD_TYPE_INT:
					dicDataType += "int>"
				case FIELD_TYPE_FLOAT:
					dicDataType += "float>"
				case FIELD_TYPE_STRING:
					dicDataType += "string>"
				case FIELD_TYPE_BOOL:
					dicDataType += "bool>"
				default:
					panic("invalid dic value type")
				}
				dicStat := fmt.Sprintf("%v{\n", leafIndent+stepIndent)
				tmpInde := leafIndent + stepIndent + stepIndent
				for k, v := range val {
					switch field.Typ.Params[0] {
					case FIELD_TYPE_INT, FIELD_TYPE_FLOAT:
						dicStat += fmt.Sprintf("%v{%v, ", tmpInde, k)
					case FIELD_TYPE_STRING:
						dicStat += fmt.Sprintf("%v{\"%v\", ", tmpInde, k)
					default:
						panic("invalid dic key type")
					}
					switch field.Typ.Params[1] {
					case FIELD_TYPE_INT, FIELD_TYPE_FLOAT:
						dicStat += fmt.Sprintf("%v},\n", v)
					case FIELD_TYPE_STRING:
						dicStat += fmt.Sprintf("\"%v\"},\n", v)
					case FIELD_TYPE_BOOL:
						tmpV := v.(bool)
						tmpStr := ""
						if tmpV {
							tmpStr = "true"
						} else {
							tmpStr = "false"
						}
						dicStat += fmt.Sprintf("%v},\n", tmpStr)
					default:
						panic("invalid dic value type")
					}
				}
				dicStat += fmt.Sprintf("%v},\n", leafIndent+stepIndent)
				line = fmt.Sprintf("%v%v = new %v()\n%v", leafIndent, field.Name, dicDataType, dicStat)
			case FIELD_TYPE_ARRAY:
				val := fieldVal.([]interface{})
				aryDataType := ""
				switch field.Typ.Params[0] {
				case FIELD_TYPE_INT:
					aryDataType += "int[]"
				case FIELD_TYPE_FLOAT:
					aryDataType += "float[]"
				case FIELD_TYPE_STRING:
					aryDataType += "string[]"
				case FIELD_TYPE_BOOL:
					aryDataType += "bool[]"
				default:
					panic("invalid dic key type")
				}
				aryStat := ""
				for _, v := range val {
					switch field.Typ.Params[0] {
					case FIELD_TYPE_INT, FIELD_TYPE_FLOAT:
						aryStat += fmt.Sprintf("%v,", v)
					case FIELD_TYPE_STRING:
						aryStat += fmt.Sprintf("\"%v\",", v)
					case FIELD_TYPE_BOOL:
						tmpV := v.(bool)
						tmpStr := ""
						if tmpV {
							tmpStr = "true"
						} else {
							tmpStr = "false"
						}
						aryStat += fmt.Sprintf("%v,", tmpStr)
					default:
						panic("invalid dic value type")
					}
				}
				line = fmt.Sprintf("%v%v = new %v {%v},\n", leafIndent, field.Name, aryDataType, aryStat)
			default:
				panic("Invalid primary key type")
			}
			buf.WriteString(line)
		}
	}
	buf.WriteString(fmt.Sprintf("%v}\n", indent+stepIndent))
	if !itsRoot {
		buf.WriteString(fmt.Sprintf("%v},\n", indent))
	}
}

func convertToCSStr(file *ExportFile, tblName, indent string, templFile string) (string, error) {
	// linesRequire := 3
	sheetTemplFilePath := templFile

	sheetTemplBuf, err := utils.ReadStringFromFile(sheetTemplFilePath)
	if err != nil {
		return "", fmt.Errorf("Read file %v error: %v", sheetTemplFilePath, err)
	}

	sheetTempl := sheetTemplBuf
	sheetTempl = utils.FillIndent(sheetTempl, indent)

	tbl, _ := file.Tables[tblName]

	confDataIndent := utils.GetIndentOf(sheetTempl, "$CONF_DATA")
	buf := &bytes.Buffer{}
	tocsWriteNodeToBuf(tbl, tbl.ContentRoot, buf, confDataIndent)

	now := time.Now()
	exportTime := now.Format("2006-01-02 15:04:05")

	classItemIndent := utils.GetIndentOf(sheetTempl, "$CONF_CLASS_ITEM")
	getFuncsIndent := utils.GetIndentOf(sheetTempl, "$GET_CONF_FUNCS")

	confClassItem, e1 := tocsGenStruct(file, tblName, classItemIndent)
	getFuncsStat, e2 := genGetDataFuncs(tbl, getFuncsIndent)
	confDataType := genNodeCSDicType(tbl, nil)
	confData := buf.String()
	dataVarName := csVars.DataVarName

	if e1 != nil {
		return "", e1
	}
	if e2 != nil {
		return "", e2
	}

	sheetTempl = strings.Replace(sheetTempl, "$CLASS_NAME", csVars.ClassName, -1)
	sheetTempl = strings.Replace(sheetTempl, "$CLASS_ITEM_NAME", csVars.ClassItemName, -1)
	sheetTempl = strings.Replace(sheetTempl, "$exportTime", exportTime, -1)
	sheetTempl = strings.Replace(sheetTempl, "$DATA_VAR", dataVarName, -1)
	sheetTempl = strings.Replace(sheetTempl, "$CONF_DATA_TYPE", confDataType, -1)
	sheetTempl = strings.Replace(sheetTempl, "$GET_CONF_FUNCS", utils.UnindentLines(getFuncsStat, 1), -1)
	sheetTempl = strings.Replace(sheetTempl, "$CONF_CLASS_ITEM", utils.UnindentLines(confClassItem, 1), -1)
	sheetTempl = strings.Replace(sheetTempl, "$CONF_DATA", utils.UnindentLines(confData, 1), -1)

	return sheetTempl, nil
}

func tocs(file *ExportFile) error {
	fileTemplFilePath := "template/tocs/file_tocs.templ"

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

	f0 := "template/tocs/conf_class_datavar_name.templ"
	f1 := "template/tocs/conf_class_name.templ"
	f2 := "template/tocs/conf_class_item_name.templ"
	classDataVarNameTempl, e0 := utils.ReadStringFromFile(f0)
	classNameTempl, e1 := utils.ReadStringFromFile(f1)
	classItemNameTempl, e2 := utils.ReadStringFromFile(f2)
	if e0 != nil {
		return fmt.Errorf("Read file %v error: %v", f0, e0)
	}
	if e1 != nil {
		return fmt.Errorf("Read file %v error: %v", f1, e1)
	}
	if e2 != nil {
		return fmt.Errorf("Read file %v error: %v", f2, e2)
	}

	indent := utils.GetIndentOf(fileTempl, "$CONFIGS")

	var buf bytes.Buffer
	for _, tbl := range file.Tables {
		csVars = &CSVars{}
		csVars.ClassName = strings.Replace(classNameTempl, "$sheetName", tbl.Name, -1)
		csVars.ClassItemName = strings.Replace(classItemNameTempl, "$sheetName", tbl.Name, -1)
		csVars.DataVarName = classDataVarNameTempl

		tblContent, err := convertToCSStr(file, tbl.Name, indent, "template/tocs/sheet_tocs.templ")
		if err != nil {
			return err
		}
		buf.WriteString(tblContent)
		buf.WriteString("\n")
	}
	tablesContent := buf.String()
	fileTempl = strings.Replace(fileTempl, "$CONFIGS", tablesContent, -1)

	os.MkdirAll(file.ExportToDir, os.ModePerm)
	err = ioutil.WriteFile(file.ExportToDir+"/"+pureFileName+".cs", []byte(fileTempl), os.ModePerm)

	return err
}

func tocsUseSheet(file *ExportFile) error {
	f0 := "template/tocs/conf_class_datavar_name.templ"
	f1 := "template/tocs/conf_class_name.templ"
	f2 := "template/tocs/conf_class_item_name.templ"
	classDataVarNameTempl, e0 := utils.ReadStringFromFile(f0)
	classNameTempl, e1 := utils.ReadStringFromFile(f1)
	classItemNameTempl, e2 := utils.ReadStringFromFile(f2)
	if e0 != nil {
		return fmt.Errorf("Read file %v error: %v", f0, e0)
	}
	if e1 != nil {
		return fmt.Errorf("Read file %v error: %v", f1, e1)
	}
	if e2 != nil {
		return fmt.Errorf("Read file %v error: %v", f2, e2)
	}

	for _, tbl := range file.Tables {
		csVars = &CSVars{}
		csVars.ClassName = strings.Replace(classNameTempl, "$sheetName", tbl.Name, -1)
		csVars.ClassItemName = strings.Replace(classItemNameTempl, "$sheetName", tbl.Name, -1)
		csVars.DataVarName = classDataVarNameTempl

		tblContent, err := convertToCSStr(file, tbl.Name, "", "template/tocs/file_tocs_usesheet.templ")
		if err != nil {
			return err
		}
		os.MkdirAll(file.ExportToDir, os.ModePerm)
		err = ioutil.WriteFile(file.ExportToDir+"/"+tbl.Name+".cs", []byte(tblContent), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
