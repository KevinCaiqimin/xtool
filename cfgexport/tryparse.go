package cfgexport

import (
	"fmt"
	"strings"

	"github.com/KevinCaiqimin/go-basic/mathext"
	// "github.com/KevinCaiqimin/go-basic/xlog"
)

func selectAString(str string, idx int) (string, int, error) {
	lenS := len(str)

	strStart, strEnd := -1, -1
	for i := idx; i < lenS; i++ {
		if strStart >= 0 { // == 0?
			if str[i] == '"' && str[i-1] != '\\' { // "\"key\"", '\'key\''
				strEnd = i //  strStart[0] == "\"" and true or false -> -1 index
				break
			}
		} else if str[i] == '"' {
			strStart = i
		} else if str[i] != ' ' {
			return "", -1, fmt.Errorf("invalid string value found")
		}
	}
	if strStart >= lenS || strEnd >= lenS || strStart < 0 || strEnd < 0 {
		return "", 0, fmt.Errorf("invalid string value found")
	}

	return str[strStart+1 : strEnd], strEnd + 1, nil
}

func selectValue(valType, strVal string, idx int, terminator byte) (interface{}, int, error) {
	lenS := len(strVal)
	keyStr := strVal[idx:]
	newIdx := lenS
	if strings.TrimSpace(keyStr) == "" {
		return nil, lenS, nil
	}
	switch valType {
	case FIELD_TYPE_INT:
		for i := idx; i < lenS; i++ {
			if strVal[i] == terminator {
				keyStr = strVal[idx:i]
				newIdx = i + 1
				break
			}
		}
		val, err := mathext.TryParseIntVal(keyStr)
		if err != nil {
			return nil, -1, err
		}
		return val, newIdx, nil
	case FIELD_TYPE_FLOAT:
		for i := idx; i < lenS; i++ {
			if strVal[i] == terminator {
				keyStr = strVal[idx:i]
				newIdx = i + 1
				break
			}
		}
		val, err := mathext.TryParseFloatVal(keyStr)
		if err != nil {
			return nil, -1, err
		}
		return val, newIdx, nil
	case FIELD_TYPE_STRING:
		str, idx, err := selectAString(strVal, idx)
		if err != nil {
			return nil, -1, err
		}
		for i := idx; i < lenS; i++ {
			if strVal[i] == terminator {
				newIdx = i + 1
				break
			} else if strVal[i] != ' ' {
				return nil, -1, fmt.Errorf("invalid dictionary key")
			}
		}
		return str, newIdx, nil
	case FIELD_TYPE_BOOL:
		for i := idx; i < lenS; i++ {
			if strVal[i] == terminator {
				keyStr = strVal[idx:i]
				newIdx = i + 1
				break
			}
		}
		val, err := mathext.TryParseBoolVal(keyStr)
		if err != nil {
			return nil, -1, err
		}
		return val, newIdx, nil
	default:
		return nil, -1, fmt.Errorf("invalid value type: %v", valType)
	}
}

func selectDicKey(keyType, strVal string, idx int) (interface{}, int, error) {
	return selectValue(keyType, strVal, idx, '=')
}

func selectDicVal(valType, strVal string, idx int) (interface{}, int, error) {
	return selectValue(valType, strVal, idx, ',')
}

func tryParseDicData(kType, vType, strVal string) (keys []interface{}, vals []interface{}, err error) {
	lenS := len(strVal)
	idx := 0
	for idx < lenS {
		k, newIdx, err := selectValue(kType, strVal, idx, '=')
		if err != nil {
			return nil, nil, err
		}
		idx = newIdx
		if k == nil {
			break
		}
		v, newIdx, err := selectValue(vType, strVal, idx, ',')
		if err != nil {
			return nil, nil, err
		}
		if v == nil {
			return nil, nil, fmt.Errorf("empty dic value")
		}
		idx = newIdx
		keys = append(keys, k)
		vals = append(vals, v)
	}

	return keys, vals, nil
}

func tryParseArrayData(fieldType string, strVal string) (vals []interface{}, err error) {
	lenS := len(strVal)

	idx := 0
	for idx < lenS {
		v, newIdx, err := selectValue(fieldType, strVal, idx, ',')
		if err != nil {
			return nil, err
		}
		if v == nil {
			break
		}
		idx = newIdx
		vals = append(vals, v)
	}

	return vals, nil
}

func tryParseSimpleVal(valType string, strVal string) (interface{}, error) {
	switch valType {
	case FIELD_TYPE_INT:
		val, err := mathext.TryParseIntVal(strVal)
		return val, err
	case FIELD_TYPE_FLOAT:
		val, err := mathext.TryParseFloatVal(strVal)
		return val, err
	case FIELD_TYPE_STRING:
		val, err := mathext.TryParseStringVal(strVal)
		val = strings.Replace(val, "\\", "\\\\", -1)
		val = strings.Replace(val, "\"", "\\\"", -1)
		val = strings.Replace(val, "'", "\\'", -1)
		val = strings.Replace(val, "\n", "\\n", -1)
		return val, err
	case FIELD_TYPE_BOOL:
		val, err := mathext.TryParseBoolVal(strVal)
		return val, err
	default:
		return nil, fmt.Errorf("invalid field type: %v", valType)
	}
}
