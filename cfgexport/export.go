package cfgexport

import (
	"helloserver/log"
	"fmt"
	"regexp"
	"strings"
	"xtool/msoffice"
	// "caiqimin.tech/basic/xlog"
	"io/ioutil"
	"path"
	"flag"
)

type ExportConfig struct {
	Js bool
	Jsdir string
	Py bool
	Pydir string
	Lua bool
	Luadir string
	Cs bool
	Csdir string
	Ignore string
	UseSheet bool

	SrcDir string
}

type ExportResult struct {
	Total int
	IgnoreCnt int
	SuccCnt int
	FailedCnt int
	Err error
}

func (this *ExportConfig)needIgnoreThisFile(fileName string)(bool) {
	if this.Ignore == "" {
		return false
	}
	ignore := strings.Replace(this.Ignore, ".", "\\.", -1)
	ignore = strings.Replace(ignore, "*", ".*", -1)
	reg := regexp.MustCompile(ignore)
	strs := reg.FindAllString(fileName, -1)
	if len(strs) > 0 {
		return true
	}
	return false
}

func needSkipThisRow(sheet *msoffice.XlsSheet, line int) bool {
	if line >= len(sheet.Cells) {
		return true
	}
	row := sheet.Cells[line]
	if len(row) <= 0 {
		return true
	}
	firstCell := row[0]
	if strings.HasPrefix(firstCell, "#") {
		return true
	}
	empty := true
	colsNum := sheet.ColsNum
	if colsNum > len(row) {
		colsNum = len(row)
	}
	for i := 0; i < colsNum; i++ {
		val := row[i]
		if strings.TrimSpace(val) == "" {
			continue
		}
		if needSkipThisCol(sheet, i) {
			continue
		}
		empty = false
		break
	}
	return empty
}

func needSkipThisCol(sheet *msoffice.XlsSheet, col int) bool {
	linesRequire := 3
	// fieldNameLineNO := 1
	// typeLineNO := 2
	// commLineNO := 3

	if len(sheet.Cells) < linesRequire {
		return true
	}

	//check if all are empty
	empty := true
	for i := 0; i < linesRequire; i++ {
		if len(sheet.Cells[i]) <= col {
			continue
		}
		val := sheet.Cells[i][col]
		val = strings.TrimSpace(val)
		if val != "" {
			empty = false
			break
		}
	}
	if empty {
		return true
	}

	//check if the first row is start with '#'
	firstRow := sheet.Cells[0]
	if len(firstRow) <= col {
		return false
	}
	val := firstRow[col]
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "#") {
		return true
	}
	return false
}

func validate(sheets []*msoffice.XlsSheet) (map[string]*ExportTable, error) {
	linesRequire := 3
	fieldNameLineNO := 1
	typeLineNO := 2
	commLineNO := 3

	nameExp := regexp.MustCompile(`\w+`)

	exportTables := map[string]*ExportTable{}

	for _, sheet := range sheets {
		//check if lines is sufficient
		if sheet.LinesNum < linesRequire {
			return nil, fmt.Errorf("[sheet: %v] require more lines", sheet.Name)
		}
		sheet.Name = strings.TrimSpace(sheet.Name)
		if strings.HasPrefix(sheet.Name, "#") {
			continue
		}
		tbl := &ExportTable{
			Name:             sheet.Name,
			Sheet:            sheet,
			Fields:           []*Field{},
			FieldsByColIndex: map[int]*Field{},
			FieldsByName:     map[string]*Field{},
		}

		succ := nameExp.MatchString(sheet.Name)
		if !succ {
			return nil, fmt.Errorf("[sheet: %v] invalid sheet name",
				sheet.Name)
		}
		//check if field name is valid
		for c := 0; c < sheet.ColsNum; c++ {
			needSkip := needSkipThisCol(sheet, c)
			if needSkip {
				continue
			}
			colName := msoffice.GetColIndexName(c)

			fieldTypeName := sheet.Cells[typeLineNO-1][c]
			fieldName := sheet.Cells[fieldNameLineNO-1][c]
			comment := sheet.Cells[commLineNO-1][c]

			fieldTypeName = strings.TrimSpace(fieldTypeName)
			fieldName = strings.TrimSpace(fieldName)

			fieldName, isPrimary, isMultiPrimary := checkPrimaryKey(fieldName)
			fieldType := parseFieldType(fieldTypeName)

			if fieldType == nil {
				return nil, fmt.Errorf("[sheet: %v, col: %v, fieldTypeName: %v] parse field type failed",
					sheet.Name, colName, fieldTypeName)
			}

			field := &Field{
				Name:      fieldName,
				Typ:       fieldType,
				Comment:   comment,
				ColIndex:  c,
				IsPrimary: isPrimary,
				IsMultiPrimary: isMultiPrimary,
			}

			if !nameExp.MatchString(fieldName) {
				return nil, fmt.Errorf("[sheet: %v, col: %v] invalid field name",
					sheet.Name, colName)
			}

			if _, succ := tbl.FieldsByName[fieldName]; succ {
				return nil, fmt.Errorf("[sheet: %v, col: %v, fieldTypeName: %v] duplicated field",
					sheet.Name, colName, fieldTypeName)
			}

			if isPrimary && !isValidPrimaryKeyFieldType(fieldType.Name) {
				return nil, fmt.Errorf("[sheet: %v, col: %v, fieldTypeName: %v] this kind of field type cannot be primary",
					sheet.Name, colName, fieldTypeName)
			}

			tbl.Fields = append(tbl.Fields, field)
			tbl.FieldsByColIndex[c] = field
			tbl.FieldsByName[field.Name] = field
			if isPrimary {
				tbl.PrimaryKeyFields = append(tbl.PrimaryKeyFields, field)
			}
			if isMultiPrimary {
				tbl.MultiPrimaryKeyFields = append(tbl.MultiPrimaryKeyFields, field)
			}
		}
		if len(tbl.PrimaryKeyFields) > 0 && len(tbl.MultiPrimaryKeyFields) > 0 {
			return nil, fmt.Errorf("[sheet: %v] we still not support mutiple primary fields $ mix with primary fields *",
				sheet.Name)
		}
		if len(tbl.MultiPrimaryKeyFields) == 0 && len(tbl.PrimaryKeyFields) == 0 {
			for _, field := range tbl.Fields {
				if isValidPrimaryKeyFieldType(field.Typ.Name) {
					tbl.PrimaryKeyFields = append(tbl.PrimaryKeyFields, field)
					field.IsPrimary = true
					break
				}
			}
		}
		if len(tbl.MultiPrimaryKeyFields) == 0 && len(tbl.PrimaryKeyFields) == 0 {
			return nil, fmt.Errorf("[sheet: %v] no proper primary key can be found",
				sheet.Name)
		}
		exportTables[tbl.Name] = tbl
	}

	return exportTables, nil
}

func exportToMap(tbl *ExportTable) (*ExportTreeNode, error) {
	linesRequire := 3
	// fieldNameLineNO := 1
	// typeLineNO := 2
	// commLineNO := 3

	sheet := tbl.Sheet

	mm := make(map[string]interface{})
	contentRoot := &ExportTreeNode{
		FieldRef: nil,
		FieldValue: nil,
		Children: map[string]*ExportTreeNode{},
	}

	//check if field name is valid
	for l := linesRequire; l < sheet.LinesNum; l++ {
		rowCells := sheet.Cells[l]
		row := map[string]interface{}{}

		if needSkipThisRow(sheet, l) {
			continue
		}


		for c := 0; c < sheet.ColsNum; c++ {
			field, exists := tbl.FieldsByColIndex[c]
			if !exists {
				continue
			}
			fieldName := field.Name

			colName := msoffice.GetColIndexName(c)
			fieldData := ""
			if c >= len(rowCells) {
				fieldData = ""
			} else {
				fieldData = rowCells[c]
			}
			val, err := field.Typ.TryParse(fieldData)
			if err != nil {
				return nil, fmt.Errorf("[sheet: %v, col: %v, line: %v] invalid data",
					sheet.Name, colName, l + 1)
			}

			row[fieldName] = val
		}
		if len(tbl.MultiPrimaryKeyFields) > 0 {
			//setup tree
			node := contentRoot
			for _, pkeyField := range tbl.MultiPrimaryKeyFields {
				val, _ := row[pkeyField.Name]
				valStr := fmt.Sprintf("%v", val)
	
				_, ok := node.Children[valStr]
				if !ok {
					node.Children[valStr] = &ExportTreeNode{
						FieldRef: pkeyField,
						FieldValue: val,
						Children: map[string]*ExportTreeNode{},
					}
					node.ChildrenKeySeq = append(node.ChildrenKeySeq, valStr)
				} else {
					return nil, fmt.Errorf("[sheet: %v, line: %v] duplicated primary key found",
						sheet.Name, l + 1)
				}
	
				node.Children[valStr].FieldValue = row
			}
		} else {
			//setup primary data structure
			tmp := mm
			for i, pkeyField := range tbl.PrimaryKeyFields {
				val, _ := row[pkeyField.Name]
				valStr := fmt.Sprintf("%v", val)
				if i < len(tbl.PrimaryKeyFields)-1 {
					if _, ok := tmp[valStr]; !ok {
						tmp[valStr] = make(map[string]interface{})
					}
				} else {
					if _, ok := tmp[valStr]; !ok {
						tmp[valStr] = row
					} else {
						return nil, fmt.Errorf("[sheet: %v, line: %v] duplicated primary key found",
							sheet.Name, l + 1)
					}
					break
				}
				tmp = tmp[valStr].(map[string]interface{})
			}
			//setup tree
			node := contentRoot
			for i, pkeyField := range tbl.PrimaryKeyFields {
				val, _ := row[pkeyField.Name]
				valStr := fmt.Sprintf("%v", val)
	
				_, ok := node.Children[valStr]
				if !ok {
					node.Children[valStr] = &ExportTreeNode{
						FieldRef: pkeyField,
						FieldValue: val,
						Children: map[string]*ExportTreeNode{},
					}
					node.ChildrenKeySeq = append(node.ChildrenKeySeq, valStr)
				}
	
				if i == len(tbl.PrimaryKeyFields)-1 {
					node.Children[valStr].FieldValue = row
				}
				node = node.Children[valStr]
			}
		}
	}
	return contentRoot, nil
}

func exportToFile(exportFile *ExportFile, exportTo string) error {
	switch exportTo {
	case "lua":
		return tolua(exportFile)
	case "py":
		return topy(exportFile)
	case "js":
		break
	default:
		return nil //TODO
	}
	return nil
}

func itsValidExportFile(fileName string) bool {
	ext := path.Ext(fileName)
	if ext == ".xlsx" {
		return true
	}
	return false
}

func exportDir(srcdir, offsetDir string, expCfg *ExportConfig) *ExportResult {
	result := &ExportResult{}
	files, err := ioutil.ReadDir(srcdir)
	if err != nil {
		result.Err = err
		return result
	}

	total := 0
	succCnt := 0
	ignoredCnt := 0
	var failedFiles []string
	for _, file := range files {
		filePath := srcdir + "/" + file.Name()
		if file.IsDir() {
			exportDir(filePath, offsetDir + "/" + file.Name(), expCfg)
		} else {
			if !itsValidExportFile(file.Name()) {
				continue
			}
			total++
			if expCfg.needIgnoreThisFile(filePath) {
				ignoredCnt++
				log.Warn("File will be ignored: %v", filePath)
				continue
			}
			err = exportFile(filePath, offsetDir, expCfg)
			if err == nil {
				log.Info("Export file %v successfully", file.Name())
				succCnt++
			} else {
				log.Error("[%v] %v", filePath, err)
				failedFiles = append(failedFiles, filePath)
			}
		}
	}
	result.Total = total
	result.IgnoreCnt = ignoredCnt
	result.SuccCnt = succCnt
	result.FailedCnt = total - succCnt

	return result
}

func exportFile(fileName, offsetDir string, expCfg *ExportConfig) error {
	sheetsData, err := msoffice.ReadAllExcelData(fileName)
	if err != nil {
		return err
	}
	mtbls, err := validate(sheetsData)
	if err != nil {
		return err
	}
	exportFile := &ExportFile{
		FileName: fileName,
		Tables: mtbls,
		ExportToDir: "",
	}
	for _, tbl := range mtbls {
		contentRoot, err := exportToMap(tbl)
		if err != nil {
			return err
		}
		tbl.ContentRoot = contentRoot
	}
	if expCfg.Js {
		exportFile.ExportToDir = expCfg.Jsdir + offsetDir
		if expCfg.UseSheet {
			err = tojsUseSheet(exportFile)
		} else {
			err = tojs(exportFile)
		}
		if err != nil {
			return err
		}
	}
	if expCfg.Lua {
		exportFile.ExportToDir = expCfg.Luadir + offsetDir
		if expCfg.UseSheet {
			err = toluaUseSheet(exportFile)
		} else {
			err = tolua(exportFile)
		}
		if err != nil {
			return err
		}
	}
	if expCfg.Py {
		exportFile.ExportToDir = expCfg.Pydir + offsetDir
		if expCfg.UseSheet {
			err = topyUseSheet(exportFile)
		} else {
			err = topy(exportFile)
		}
		if err != nil {
			return err
		}
	}
	if expCfg.Cs {
		exportFile.ExportToDir = expCfg.Csdir + offsetDir
		if expCfg.UseSheet {
			err = tocsUseSheet(exportFile)
		} else {
			err = tocs(exportFile)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func ExportCfg(expCfg *ExportConfig) error {
	log.Info("Export start")
	if !(expCfg.Js || expCfg.Lua || expCfg.Py || expCfg.Cs) {
		log.Error("at least one type that config will be export to be need to be enabled")
		flag.Usage()
		return nil
	}

	result := exportDir(expCfg.SrcDir, "", expCfg)

	if result.Err != nil || result.Total != result.SuccCnt {
		if result.Err != nil {
			log.Error("Export config failed with error: %v", result.Err)
		} else {
			log.Warn("Total: %v, Success: %v, Failed: %v, Ignored: %v", result.Total, result.SuccCnt, result.FailedCnt, result.IgnoreCnt)
		}
	} else {
		log.Info("Export completly and succefully")
	}

	return nil
}
