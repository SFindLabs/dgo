package utils

import (
	kruntime "dgo/framework/tools/runtime"
	kinit "dgo/work/base/initialize"
	"bytes"
	"fmt"
	"sort"
	"strings"
)

var GenerateSql *Generate

type Generate struct {
}

func (ts *Generate) Run(dbStr string, tmpDbId int64, tableSchema, tableName string, isSplit, isDivide, isRead bool) string {

	tableMap, err := getTableStructure(dbStr, tmpDbId, tableSchema, tableName)
	if err != nil {
		kinit.LogError.Println("get table structure fail :", err)
		return err.Error()
	}

	createStr := createTable(tableName, tableMap)

	insertStr := insertTable(tableName, tableMap, isSplit, isDivide)

	getAllStr := getAllTable(tableName, isSplit, isDivide, isRead)

	tableIndexArr, err := getTableIndex(dbStr, tmpDbId, tableSchema, tableName)
	if err != nil {
		kinit.LogError.Println("get table structure fail :", err)
		return err.Error()
	}

	getStr := getTable(tableName, tableIndexArr, isSplit, isDivide, isRead)

	updateStr := updateTable(tableName, tableIndexArr, isSplit, isDivide)

	updateMustStr := updateMustTable(tableName, tableIndexArr, isSplit, isDivide)

	deleteStr := deleteTable(tableName, tableIndexArr, isSplit, isDivide)

	deleteMustStr := deleteMustTable(tableName, tableIndexArr, isSplit, isDivide)

	return fmt.Sprintf("\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s", createStr, insertStr, getAllStr, getStr, updateStr, updateMustStr, deleteStr, deleteMustStr)
}

//---------------------------------------------------------------------------

func createTable(tableName string, rows map[string]tableType) string {
	var sql string
	upperTableName := convertUpper(tableName, true)
	path, err := kruntime.GetCurrentPath()
	if err != nil {
		kinit.LogError.Println("get path fail :", err)
		return err.Error()
	}
	slicePath := strings.Split(path, "/")
	dirName := slicePath[len(slicePath)-1]
	sql += "package model\n\nimport (\n\t\"errors\"\n\tkinit \"" + dirName + "/work/base/initialize\"\n\tjgorm \"github.com/jinzhu/gorm\"\n\t//\"time\"\n)"
	sql += "\n\nvar " + upperTableName + "Obj " + upperTableName + "\n\nvar " + convertUpper(tableName, false) + "DB = \"\"\n\ntype " + upperTableName + " struct {\n"

	for rowKey, rowVal := range rows {
		if "ID" == strings.ToUpper(rowKey) {
			sql += "\tID     int64  `gorm:\"primary_key\" json:\"-\"`\n"
		} else {
			sql += "\t" + convertUpper(rowKey, true) + " " + getType(rowVal.dataType, rowKey) + " `gorm:\"column:" + rowKey + "\" json:\"" + rowKey + "\"`\n"
		}
	}

	sql += "}\n\n"
	sql += "func (" + upperTableName + ") TableName() string {\n"
	sql += "\treturn \"" + tableName + "\"\n}\n"

	return sql
}

//---------------------------------------------------------------------------

func insertTable(tableName string, rows map[string]tableType, isSplit, isDivide bool) string {
	var bf bytes.Buffer
	var bt bytes.Buffer
	dbName := convertUpper(tableName, false)
	upperTableName := convertUpper(tableName, true)
	if isSplit {
		_, _ = fmt.Fprintf(&bf, "func (%s) Insert(tx *jgorm.DB, dbId int64", upperTableName)
		if isDivide {
			_, _ = fmt.Fprintf(&bt, ` {
	if tx == nil {
		tx, _ = kinit.GetMysqlDivideConnect(`+dbName+`DB, dbId)
	}`)
		} else {
			_, _ = fmt.Fprintf(&bt, ` {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect(`+dbName+`DB, dbId)
	}`)
		}
	} else {
		_, _ = fmt.Fprintf(&bf, "func (%s) Insert(tx *jgorm.DB", upperTableName)
		_, _ = fmt.Fprintf(&bt, ` { 
	if tx == nil {  
		tx, _ = kinit.GetMysqlConnect(`+dbName+`DB)
	}`)
	}
	_, _ = fmt.Fprintf(&bt, "\n\t//timeStr := time.Now().Format(\"2006-01-02 15:04:05\")\n\tobj := %s{\n", upperTableName)

	for rowKey, rowVal := range rows {
		if rowVal.columnKey != "PRI" {
			_, _ = fmt.Fprintf(&bf, ", %s %s", convertUpper(rowKey, false), getType(rowVal.dataType, rowKey))
			_, _ = fmt.Fprintf(&bt, "\t\t%s:%s,\n", convertUpper(rowKey, true), convertUpper(rowKey, false))
		}
	}

	_, _ = fmt.Fprintf(&bf, ")(%s, error)", upperTableName)
	_, _ = fmt.Fprintf(&bt, `	}
	if err := tx.Create(&obj).Error; err != nil {
		kinit.LogError.Println(err)
		return obj, err
	}
	return obj, nil
}`)

	return bf.String() + bt.String()
}

//---------------------------------------------------------------------------

func getAllTable(tableName string, isSplit, isDivide, isRead bool) string {
	var sql string
	dbName := convertUpper(tableName, false)
	db := "tx *jgorm.DB"
	dbConn := `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB)
	}`
	if isRead {
		dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlReadConnect(` + dbName + `DB)
	}`
	}
	if isSplit {
		db = "tx *jgorm.DB, dbId int64"
		if !isRead {
			if isDivide {
				dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlDivideConnect(` + dbName + `DB, dbId)
	}`
			} else {
				dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB, dbId)
	}`
			}
		} else {
			if isDivide {
				dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlReadDivideConnect(` + dbName + `DB, dbId)
	}`
			} else {
				dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlReadConnect(` + dbName + `DB, dbId)
	}`
			}
		}
	}

	upperTableName := convertUpper(tableName, true)
	sql += fmt.Sprintf(`
func (%s) GetAll(%s) (objs []%s) {
	%s
	tx.Find(&objs)
	return
}%s`, upperTableName, db, upperTableName, dbConn, "\n")

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------

func getTable(tableName string, rows []tableIndexParam, isSplit, isDivide, isRead bool) string {
	var sql string
	dbName := convertUpper(tableName, false)
	db := "tx *jgorm.DB"
	dbConn := `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB)
	}`

	if isRead {
		dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlReadConnect(` + dbName + `DB)
	}`
	}
	if isSplit {
		db = "tx *jgorm.DB, dbId int64"
		if !isRead {
			if isDivide {
				dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlDivideConnect(` + dbName + `DB, dbId)
	}`
			} else {
				dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB, dbId)
	}`
			}
		} else {
			if isDivide {
				dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlReadDivideConnect(` + dbName + `DB, dbId)
	}`
			} else {
				dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlReadConnect(` + dbName + `DB, dbId)
	}`
			}
		}
	}

	upperTableName := convertUpper(tableName, true)
	sql += ""

	for _, row := range rows {
		sql += fmt.Sprintf(`
func (%s) GetBy%s(%s, %s) (obj %s) {
	%s
	tx.Where("%s", %s).First(&obj)
	return
}%s`, upperTableName, row.byStr, db, row.paramStr, upperTableName, dbConn, row.whereStr, row.whereParam, "\n")
	}

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------

func updateTable(tableName string, rows []tableIndexParam, isSplit, isDivide bool) string {
	var sql string
	dbName := convertUpper(tableName, false)
	upperTableName := convertUpper(tableName, true)
	db := "tx *jgorm.DB"
	dbConn := `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB)
	}`
	if isSplit {
		db = "tx *jgorm.DB, dbId int64"
		dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB,dbId)
	}`
		if isDivide {
			dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlDivideConnect(` + dbName + `DB,dbId)
	}`
		}
	}

	sql += ""
	for _, row := range rows {
		sql += fmt.Sprintf(`
func (%s) UpdateBy%s(%s, %s, updateMap map[string]interface{}) error {
	%s
	if err := tx.Model(%s{}).Where("%s", %s).Updates(updateMap).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}%s`, upperTableName, row.byStr, db, row.paramStr, dbConn, upperTableName, row.whereStr, row.whereParam, "\n")
	}

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------

func updateMustTable(tableName string, rows []tableIndexParam, isSplit, isDivide bool) string {
	var sql string
	dbName := convertUpper(tableName, false)
	upperTableName := convertUpper(tableName, true)
	db := "tx *jgorm.DB"
	dbConn := `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB)
	}`
	if isSplit {
		db = "tx *jgorm.DB, dbId int64"
		dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB,dbId)
	}`
		if isDivide {
			dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlDivideConnect(` + dbName + `DB,dbId)
	}`
		}
	}

	sql += ""
	for _, row := range rows {
		sql += fmt.Sprintf(`
func (%s) UpdateMustBy%s(%s, %s, updateMap map[string]interface{}) error {
	%s
	result := tx.Model(%s{}).Where("%s", %s).Updates(updateMap)
	if result.RowsAffected == 0 {
		errMsg := errors.New("%s UpdateBy%s failed, rows is zero")
		if result.Error != nil {
			errMsg = result.Error
		}
		kinit.LogError.Println(errMsg)
		return errMsg
	}
	return nil
}%s`, upperTableName, row.byStr, db, row.paramStr, dbConn, upperTableName, row.whereStr, row.whereParam, upperTableName, row.byStr, "\n")
	}

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------

func deleteTable(tableName string, rows []tableIndexParam, isSplit, isDivide bool) string {
	var sql string
	upperTableName := convertUpper(tableName, true)
	dbName := convertUpper(tableName, false)
	db := "tx *jgorm.DB"
	dbConn := `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB)
	}`
	if isSplit {
		db = "tx *jgorm.DB, dbId int64"
		dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB,dbId)
	}`
		if isDivide {
			dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlDivideConnect(` + dbName + `DB,dbId)
	}`
		}
	}

	sql += ""
	for _, row := range rows {
		sql += fmt.Sprintf(`
func (%s) DeleteBy%s(%s, %s) error {
	%s
	var objs %s
	if err := tx.Where("%s", %s).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}%s`, upperTableName, row.byStr, db, row.paramStr, dbConn, upperTableName, row.whereStr, row.whereParam, "\n")
	}

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------

func deleteMustTable(tableName string, rows []tableIndexParam, isSplit, isDivide bool) string {
	var sql string
	upperTableName := convertUpper(tableName, true)
	dbName := convertUpper(tableName, false)
	db := "tx *jgorm.DB"
	dbConn := `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB)
	}`
	if isSplit {
		db = "tx *jgorm.DB, dbId int64"
		dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlConnect(` + dbName + `DB,dbId)
	}`
		if isDivide {
			dbConn = `if tx == nil {
		tx, _ = kinit.GetMysqlDivideConnect(` + dbName + `DB,dbId)
	}`
		}
	}

	sql += ""
	for _, row := range rows {
		sql += fmt.Sprintf(`
func (%s) DeleteMustBy%s(%s, %s) error {
	%s
	var objs %s
	result := tx.Where("%s", %s).Delete(objs)
	if result.RowsAffected == 0 {
		errMsg := errors.New("%s DeleteBy%s failed, rows is zero")
		if result.Error != nil {
			errMsg = result.Error
		}
		kinit.LogError.Println(errMsg)
		return errMsg
	}
	return nil
}%s`, upperTableName, row.byStr, db, row.paramStr, dbConn, upperTableName, row.whereStr, row.whereParam, upperTableName, row.byStr, "\n")
	}

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------
//获取表结构
type tableType struct {
	dataType  string
	columnKey string
}

func getTableStructure(dbStr string, tmpDbId int64, tableSchema, tableName string) (map[string]tableType, error) {
	tableMap := make(map[string]tableType)
	sqlStr := "select column_name,data_type,column_key from information_schema.columns where table_schema=? and table_name=?;"
	model, err := kinit.GetMysqlConnect(dbStr, tmpDbId)
	if err != nil {
		return tableMap, err
	}
	rows, err := model.Raw(sqlStr, tableSchema, tableName).Rows()
	if err != nil {
		return tableMap, err
	}

	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		columnName, dataType, columnKey := "", "", ""
		err := rows.Scan(&columnName, &dataType, &columnKey)
		if err != nil {
			return tableMap, err
		}
		tableMap[columnName] = tableType{
			dataType:  dataType,
			columnKey: columnKey,
		}
	}

	return tableMap, nil
}

//---------------------------------------------------------------------------
//获取表索引
type tableIndexParam struct {
	byStr      string
	paramStr   string
	whereStr   string
	whereParam string
}

func getTableIndex(dbStr string, tmpDbId int64, tableSchema, tableName string) ([]tableIndexParam, error) {
	tmpIndexSlice := make([]tableIndexParam, 0)
	//获取列名索引以及索引顺序
	sqlStr := "select column_name, index_name, seq_in_index from information_schema.statistics where table_schema=? and table_name=?;"
	model, err := kinit.GetMysqlConnect(dbStr, tmpDbId)
	if err != nil {
		return tmpIndexSlice, err
	}
	rows, err := model.Raw(sqlStr, tableSchema, tableName).Rows()
	if err != nil {
		return tmpIndexSlice, err
	}
	defer func() {
		_ = rows.Close()
	}()

	tmpUnion := make(map[string]map[int]string)
	columnArr := make([]string, 0)

	for rows.Next() {
		columnName, indexName, indexSeq := "", "", 0
		err := rows.Scan(&columnName, &indexName, &indexSeq)
		if err != nil {
			return tmpIndexSlice, err
		}
		if _, ok := tmpUnion[indexName]; ok {
			tmpUnion[indexName][indexSeq] = columnName
		} else {
			tmpUnion[indexName] = map[int]string{indexSeq: columnName}
		}
		columnArr = append(columnArr, columnName)
	}

	//获取列名以及数据类型
	sqlStr = "select column_name, data_type from information_schema.columns where table_schema=? and table_name=? and column_name in (?);"
	columnRows, err := model.Raw(sqlStr, tableSchema, tableName, columnArr).Rows()
	if err != nil {
		return tmpIndexSlice, err
	}
	defer func() {
		_ = columnRows.Close()
	}()

	tmpColumnMap := make(map[string]string)
	for columnRows.Next() {
		columnName, dataType := "", ""
		err := columnRows.Scan(&columnName, &dataType)
		if err != nil {
			return tmpIndexSlice, err
		}
		tmpColumnMap[columnName] = dataType
	}

	//拼接列数据
	for _, v := range tmpUnion {
		byStr, paramStr, whereStr, whereParam := "", "", "", ""
		var tmpKeys []int
		for tmpKey := range v {
			tmpKeys = append(tmpKeys, tmpKey)
		}
		sort.Ints(tmpKeys)
		for _, key := range tmpKeys {
			tmpVal := v[key]
			if byStr == "" {
				byStr += fmt.Sprintf("%s", convertUpper(tmpVal, true))
			} else {
				byStr += fmt.Sprintf("And%s", convertUpper(tmpVal, true))
			}
			if paramStr == "" {
				paramStr += fmt.Sprintf("%s %s", convertUpper(tmpVal, false), getType(tmpColumnMap[tmpVal], tmpVal))
			} else {
				paramStr += fmt.Sprintf(", %s %s", convertUpper(tmpVal, false), getType(tmpColumnMap[tmpVal], tmpVal))
			}
			if whereStr == "" {
				whereStr += fmt.Sprintf("%s = ?", tmpVal)
			} else {
				whereStr += fmt.Sprintf(" and %s = ?", tmpVal)
			}
			if whereParam == "" {
				whereParam += fmt.Sprintf("%s", convertUpper(tmpVal, false))
			} else {
				whereParam += fmt.Sprintf(", %s", convertUpper(tmpVal, false))
			}
		}
		tmpIndexSlice = append(tmpIndexSlice, tableIndexParam{
			byStr,
			paramStr,
			whereStr,
			whereParam,
		})
	}

	return tmpIndexSlice, nil
}

//---------------------------------------------------------------------------

func getType(s, name string) string {
	if s == "smallint" {
		return "int64"
	}
	if s == "varchar" {
		return "string"
	}
	if s == "tinyint" {
		return "int64"
	}
	if s == "mediumint" {
		return "int64"
	}
	if s == "int" {
		if strings.Index(name, "time") != -1 {
			return "int64"
		}
		return "int64"
	}
	if s == "text" {
		return "string"
	}
	if s == "mediumtext" {
		return "string"
	}
	if s == "char" {
		return "string"
	}
	if s == "mediumblob" {
		return "string"
	}
	if s == "enum" {
		return "string"
	}
	if s == "float" {
		return "string"
	}
	if s == "date" {
		return "string"
	}
	if s == "decimal" {
		return "float64"
	}
	if s == "double" {
		return "float64"
	}
	if s == "longtext" {
		return "string"
	}
	if s == "bigint" {
		return "int64"
	}
	if s == "datetime" {
		return "string"
	}
	if s == "blob" {
		return "string"
	}
	if s == "varbinary" {
		return "string"
	}
	if s == "timestamp" {
		return "string"
	}
	if s == "set" {
		return "string"
	}
	if s == "longblob" {
		return "string"
	}
	if s == "time" {
		return "string"
	}

	return ""
}

func convertUpper(s string, isAll bool) string {
	tmp := strings.Split(s, "_")
	var res string
	for i := 0; i < len(tmp); i++ {
		if !isAll && i == 0 {
			res += tmp[i]
		} else {
			v := []rune(tmp[i])
			for y := 0; y < len(v); y++ {
				if y == 0 {
					if v[y] >= 97 && v[y] <= 122 {
						v[y] -= 32
					}
					res += string(v[y])
				} else {
					res += string(v[y])
				}
			}
		}
	}
	return res
}
