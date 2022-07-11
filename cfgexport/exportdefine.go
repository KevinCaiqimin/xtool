package cfgexport

var FIELD_TYPE_INT = "int"
var FIELD_TYPE_FLOAT = "float"
var FIELD_TYPE_STRING = "string"
var FIELD_TYPE_BOOL = "bool"
var FIELD_TYPE_DIC = "dic"
var FIELD_TYPE_ARRAY = "array"
var FIELD_TYPE_LUA = "lua"

var fieldTypeDic = map[string]string{
	FIELD_TYPE_INT:    FIELD_TYPE_INT,
	FIELD_TYPE_FLOAT:  FIELD_TYPE_FLOAT,
	FIELD_TYPE_STRING: FIELD_TYPE_STRING,
	FIELD_TYPE_BOOL: FIELD_TYPE_BOOL,
	FIELD_TYPE_DIC:    FIELD_TYPE_DIC,
	FIELD_TYPE_ARRAY:  FIELD_TYPE_ARRAY,
	FIELD_TYPE_LUA:    FIELD_TYPE_LUA,
}

func getIndent() string {
	return "    "
}