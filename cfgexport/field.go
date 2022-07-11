package cfgexport

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/KevinCaiqimin/xtool/msoffice"

	"github.com/KevinCaiqimin/go-basic/mathext"
	// "github.com/KevinCaiqimin/go-basic/xlog"
)

type FieldType struct {
	Name   string
	Params []string
}

type Field struct {
	Name           string
	Typ            *FieldType
	Comment        string
	ColIndex       int
	IsPrimary      bool
	EnableDefault  bool
	IsMultiPrimary bool //多主键
}

type ExportTreeNode struct {
	FieldRef       *Field
	FieldValue     interface{}
	Children       map[string]*ExportTreeNode
	ChildrenKeySeq []string
}

type ExportTable struct {
	Name                  string
	Fields                []*Field
	FieldsByColIndex      map[int]*Field
	FieldsByName          map[string]*Field
	PrimaryKeyFields      []*Field
	MultiPrimaryKeyFields []*Field
	Sheet                 *msoffice.XlsSheet
	ContentRoot           *ExportTreeNode
	Content               map[string]interface{}
}

type ExportFile struct {
	FileName    string
	ExportToDir string
	Tables      map[string]*ExportTable
}

func (this *FieldType) IsValidVal(strVal string) bool {
	switch this.Name {
	case FIELD_TYPE_INT, FIELD_TYPE_FLOAT, FIELD_TYPE_STRING, FIELD_TYPE_BOOL:
		_, err := tryParseSimpleVal(this.Name, strVal)
		return err == nil
	case FIELD_TYPE_LUA:
		_, err := mathext.TryParseStringVal(strVal)
		return err == nil
	case FIELD_TYPE_DIC:
		_, _, err := tryParseDicData(this.Params[0], this.Params[1], strVal)
		return err == nil
	case FIELD_TYPE_ARRAY:
		_, err := tryParseArrayData(this.Params[0], strVal)
		return err == nil
	}
	return false
}

func (this *FieldType) TryParse(strVal string) (interface{}, error) {
	switch this.Name {
	case FIELD_TYPE_INT, FIELD_TYPE_FLOAT, FIELD_TYPE_STRING, FIELD_TYPE_BOOL:
		return tryParseSimpleVal(this.Name, strVal)
	case FIELD_TYPE_LUA:
		return mathext.TryParseStringVal(strVal)
	case FIELD_TYPE_DIC:
		keys, vals, err := tryParseDicData(this.Params[0], this.Params[1], strVal)
		if err != nil {
			return nil, err
		}
		m := make(map[interface{}]interface{})
		for i, k := range keys {
			v := vals[i]
			m[k] = v
		}
		return m, nil
	case FIELD_TYPE_ARRAY:
		ary, err := tryParseArrayData(this.Params[0], strVal)
		if err != nil {
			return nil, err
		}
		return ary, nil
	default:
		return nil, fmt.Errorf("invalid field type: %s", this.Name)
	}
	return nil, nil
}

func isFieldTypeSupported(typ string) bool {
	_, ok := fieldTypeDic[typ]
	return ok
}

func tryParseArrayField(fieldType string) (bool, *FieldType) {
	param := "(int|float|string)"
	typeExp, err := regexp.Compile(" *" + param + "\\[ *\\] *")
	if err != nil {
		return false, nil
	}
	paramExp, err := regexp.Compile(param)
	if err != nil {
		return false, nil
	}
	matched := typeExp.MatchString(fieldType)
	if !matched {
		return false, nil
	}
	paramsStr := paramExp.FindString(fieldType)
	paramsStr = strings.TrimSpace(paramsStr)

	typ := &FieldType{
		Name:   FIELD_TYPE_ARRAY,
		Params: []string{},
	}
	typ.Params = append(typ.Params, paramsStr)

	return true, typ
}

func tryParseDicField(fieldType string) (bool, *FieldType) {
	param := " *(int|float|string) *\\, *(int|float|string) *"
	typeExp, err := regexp.Compile(" *dic\\<" + param + "\\> *")
	if err != nil {
		return false, nil
	}
	paramExp, err := regexp.Compile(param)
	if err != nil {
		return false, nil
	}
	matched := typeExp.MatchString(fieldType)
	if !matched {
		return false, nil
	}
	paramsStr := paramExp.FindString(fieldType)

	typ := &FieldType{
		Name:   FIELD_TYPE_DIC,
		Params: []string{},
	}
	tmp := strings.Split(paramsStr, ",")
	for _, paramStr := range tmp {
		paramStr := strings.TrimSpace(paramStr)
		if paramStr == "" {
			continue
		}
		typ.Params = append(typ.Params, paramStr)
	}

	return true, typ
}

func parseFieldType(fieldType string) *FieldType {
	typ := &FieldType{
		Name:   "",
		Params: []string{},
	}
	//check simple type
	switch fieldType {
	case FIELD_TYPE_INT, FIELD_TYPE_FLOAT, FIELD_TYPE_STRING, FIELD_TYPE_BOOL, FIELD_TYPE_LUA:
		typ.Name = fieldType
		break
	default:
		succ, typ := tryParseDicField(fieldType)
		if succ {
			return typ
		}
		succ, typ = tryParseArrayField(fieldType)
		if succ {
			return typ
		}
		return nil
	}
	return typ
}

func isValidPrimaryKeyFieldType(fieldTypeName string) bool {
	switch fieldTypeName {
	case FIELD_TYPE_INT, FIELD_TYPE_FLOAT, FIELD_TYPE_STRING:
		return true
	}
	return false
}

func checkPrimaryKey(fieldName string) (string, bool, bool) {
	if strings.HasPrefix(fieldName, "*") {
		return strings.TrimLeft(fieldName, "*"), true, false
	}
	if strings.HasPrefix(fieldName, "$") {
		return strings.TrimLeft(fieldName, "$"), false, true
	}
	return fieldName, false, false
}
