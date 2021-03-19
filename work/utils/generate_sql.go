package utils

import (
	kruntime "dgo/framework/tools/runtime"
	kinit "dgo/work/base/initialize"
	"bytes"
	ksql "database/sql"
	"fmt"
	"strings"
)

var GenerateSql *Generate

type Generate struct {
}

func (ts *Generate) Run(tableSchema, tableName string, isSplit, isDivide, isRead bool) string {
	sqlStr := "select column_name,data_type,column_key from information_schema.columns where table_schema=? and table_name=?;"
	model, _ := kinit.GetMysqlConnect("")
	rows, err := model.Raw(sqlStr, tableSchema, tableName).Rows()
	if err != nil {
		kinit.LogError.Println("get  fail :", err)
		return err.Error()
	}
	createStr := createTable(tableName, rows)

	rows, err = model.Raw(sqlStr, tableSchema, tableName).Rows()
	if err != nil {
		kinit.LogError.Println("get  fail :", err)
		return err.Error()
	}
	insertStr := insertTable(tableName, rows, isSplit, isDivide)

	sqlStr = "select column_name from information_schema.statistics where table_schema=? and table_name=?;"

	getAllStr := getAllTable(tableName, isSplit, isDivide, isRead)

	rows, err = model.Raw(sqlStr, tableSchema, tableName).Rows()
	if err != nil {
		kinit.LogError.Println("get  fail :", err)
		return err.Error()
	}
	getStr := getTable(tableSchema, tableName, rows, isSplit, isDivide, isRead)

	rows, err = model.Raw(sqlStr, tableSchema, tableName).Rows()
	if err != nil {
		kinit.LogError.Println("get  fail :", err)
		return err.Error()
	}
	updateStr := updateTable(tableSchema, tableName, rows, isSplit, isDivide)

	rows, err = model.Raw(sqlStr, tableSchema, tableName).Rows()
	if err != nil {
		kinit.LogError.Println("get  fail :", err)
		return err.Error()
	}
	updateMustStr := updateMustTable(tableSchema, tableName, rows, isSplit, isDivide)

	rows, err = model.Raw(sqlStr, tableSchema, tableName).Rows()
	if err != nil {
		kinit.LogError.Println("get  fail :", err)
		return err.Error()
	}
	deleteStr := deleteTable(tableSchema, tableName, rows, isSplit, isDivide)

	rows, err = model.Raw(sqlStr, tableSchema, tableName).Rows()
	if err != nil {
		kinit.LogError.Println("get  fail :", err)
		return err.Error()
	}
	deleteMustStr := deleteMustTable(tableSchema, tableName, rows, isSplit, isDivide)

	return fmt.Sprintf("\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s", createStr, insertStr, getAllStr, getStr, updateStr, updateMustStr, deleteStr, deleteMustStr)
}

//---------------------------------------------------------------------------

func createTable(tableName string, rows *ksql.Rows) string {
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
	for rows.Next() {
		columnName, dataType, columnKey := "", "", ""
		err := rows.Scan(&columnName, &dataType, &columnKey)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}
		if "ID" == strings.ToUpper(columnName) {
			sql += "\tID     int64  `gorm:\"primary_key\" json:\"-\"`\n"
		} else {
			sql += "\t" + convertUpper(columnName, true) + " " + getType(dataType, columnName) + " `gorm:\"column:" + columnName + "\" json:\"" + columnName + "\"`\n"
		}
	}
	sql += "}\n\n"
	sql += "func (" + upperTableName + ") TableName() string {\n"
	sql += "\treturn \"" + tableName + "\"\n}\n"

	return sql
}

//---------------------------------------------------------------------------

func insertTable(tableName string, rows *ksql.Rows, isSplit, isDivide bool) string {
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
	for rows.Next() {
		columnName, dataType, columnKey := "", "", ""
		_ = rows.Scan(&columnName, &dataType, &columnKey)

		if columnKey != "PRI" {
			_, _ = fmt.Fprintf(&bf, ", %s %s", convertUpper(columnName, false), getType(dataType, columnName))
			_, _ = fmt.Fprintf(&bt, "\t\t%s:%s,\n", convertUpper(columnName, true), convertUpper(columnName, false))
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
	return objs
}%s`, upperTableName, db, upperTableName, dbConn, "\n")

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------

func getTable(tableSchema, tableName string, rows *ksql.Rows, isSplit, isDivide, isRead bool) string {
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
	for rows.Next() {
		columnName, dataType := "", ""
		err := rows.Scan(&columnName)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}

		sqlStr := "select data_type from information_schema.columns where table_schema=? and table_name=? and column_name = ?;"
		model, _ := kinit.GetMysqlConnect("")
		row, err := model.Raw(sqlStr, tableSchema, tableName, columnName).Rows()
		if err != nil {
			kinit.LogError.Println("get  fail :", err)
			return err.Error()
		}
		if !row.Next() {
			kinit.LogError.Println("no next body row")
			return ""
		}
		err = row.Scan(&dataType)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}

		sql += fmt.Sprintf(`
func (%s) GetBy%s(%s, %s %s) (objs %s) {
	%s
	tx.Where("%s=? ", %s).First(&objs)
	return objs
}%s`, upperTableName, convertUpper(columnName, true), db, convertUpper(columnName, false), getType(dataType, columnName), upperTableName, dbConn, columnName, convertUpper(columnName, false), "\n")
	}

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------

func updateTable(tableSchema, tableName string, rows *ksql.Rows, isSplit, isDivide bool) string {
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
	for rows.Next() {
		columnName, dataType := "", ""
		err := rows.Scan(&columnName)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}

		sqlStr := "select data_type from information_schema.columns where table_schema=? and table_name=? and column_name = ?;"
		model, _ := kinit.GetMysqlConnect("")
		row, err := model.Raw(sqlStr, tableSchema, tableName, columnName).Rows()
		if err != nil {
			kinit.LogError.Println("get  fail :", err)
			return err.Error()
		}
		if !row.Next() {
			kinit.LogError.Println("no next body row")
			return ""
		}
		err = row.Scan(&dataType)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}

		sql += fmt.Sprintf(`
func (%s) UpdateBy%s(%s, %s %s, updateMap map[string]interface{}) error {
	%s
	if err := tx.Model(%s{}).Where("%s=?", %s).Updates(updateMap).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}%s`, upperTableName, convertUpper(columnName, true), db, convertUpper(columnName, false), getType(dataType, columnName), dbConn, upperTableName, columnName, convertUpper(columnName, false), "\n")
	}

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------

func updateMustTable(tableSchema, tableName string, rows *ksql.Rows, isSplit, isDivide bool) string {
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
	for rows.Next() {
		columnName, dataType := "", ""
		err := rows.Scan(&columnName)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}

		sqlStr := "select data_type from information_schema.columns where table_schema=? and table_name=? and column_name = ?;"
		model, _ := kinit.GetMysqlConnect("")
		row, err := model.Raw(sqlStr, tableSchema, tableName, columnName).Rows()
		if err != nil {
			kinit.LogError.Println("get  fail :", err)
			return err.Error()
		}
		if !row.Next() {
			kinit.LogError.Println("no next body row")
			return ""
		}
		err = row.Scan(&dataType)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}

		sql += fmt.Sprintf(`
func (%s) UpdateMustBy%s(%s, %s %s, updateMap map[string]interface{}) error {
	%s
	result := tx.Model(%s{}).Where("%s=?", %s).Updates(updateMap)
	if result.RowsAffected == 0 {
		errMsg := errors.New("%s UpdateBy%s failed, rows is zero")
		if result.Error != nil {
			errMsg = result.Error
		}
		kinit.LogError.Println(errMsg)
		return errMsg
	}
	return nil
}%s`, upperTableName, convertUpper(columnName, true), db, convertUpper(columnName, false), getType(dataType, columnName), dbConn, upperTableName, columnName, convertUpper(columnName, false), upperTableName, convertUpper(columnName, true), "\n")
	}

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------

func deleteTable(tableSchema, tableName string, rows *ksql.Rows, isSplit, isDivide bool) string {
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
	for rows.Next() {
		columnName, dataType := "", ""
		err := rows.Scan(&columnName)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}

		sqlStr := "select data_type from information_schema.columns where table_schema=? and table_name=? and column_name = ?;"
		model, _ := kinit.GetMysqlConnect("")
		row, err := model.Raw(sqlStr, tableSchema, tableName, columnName).Rows()
		if err != nil {
			kinit.LogError.Println("get  fail :", err)
			return err.Error()
		}
		if !row.Next() {
			kinit.LogError.Println("no next body row")
			return ""
		}
		err = row.Scan(&dataType)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}

		sql += fmt.Sprintf(`
func (%s) DeleteBy%s(%s, %s %s) error {
	%s
	var objs %s
		if err := tx.Where("%s=? ", %s).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}%s`, upperTableName, convertUpper(columnName, true), db, convertUpper(columnName, false), getType(dataType, columnName), dbConn, upperTableName, columnName, convertUpper(columnName, false), "\n")
	}

	return strings.TrimRight(sql, "\n")
}

//---------------------------------------------------------------------------

func deleteMustTable(tableSchema, tableName string, rows *ksql.Rows, isSplit, isDivide bool) string {
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
	for rows.Next() {
		columnName, dataType := "", ""
		err := rows.Scan(&columnName)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}

		sqlStr := "select data_type from information_schema.columns where table_schema=? and table_name=? and column_name = ?;"
		model, _ := kinit.GetMysqlConnect("")
		row, err := model.Raw(sqlStr, tableSchema, tableName, columnName).Rows()
		if err != nil {
			kinit.LogError.Println("get  fail :", err)
			return err.Error()
		}
		if !row.Next() {
			kinit.LogError.Println("no next body row")
			return ""
		}
		err = row.Scan(&dataType)
		if err != nil {
			kinit.LogError.Println("scan fail :", err)
			return ""
		}

		sql += fmt.Sprintf(`
func (%s) DeleteMustBy%s(%s, %s %s) error {
	%s
	var objs %s
	result := tx.Where("%s=? ", %s).Delete(objs)
	if result.RowsAffected == 0 {
		errMsg := errors.New("%s DeleteBy%s failed, rows is zero")
		if result.Error != nil {
			errMsg = result.Error
		}
		kinit.LogError.Println(errMsg)
		return errMsg
	}
	return nil
}%s`, upperTableName, convertUpper(columnName, true), db, convertUpper(columnName, false), getType(dataType, columnName), dbConn, upperTableName, columnName, convertUpper(columnName, false), upperTableName, convertUpper(columnName, true), "\n")
	}

	return strings.TrimRight(sql, "\n")
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
