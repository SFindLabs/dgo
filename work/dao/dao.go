package dao

import (
	kinit "dgo/work/base/initialize"
	kutils "dgo/work/utils"
	"encoding/hex"
	"fmt"
	jgorm "github.com/jinzhu/gorm"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ScanData(tx *jgorm.DB, name string, dbId int64, sql string, srcValues ...interface{}) (records []map[string]interface{}, err error) {
	if tx == nil {
		tx, err = kinit.GetMysqlConnect(name, dbId)
		if err != nil {
			return
		}
	}
	rows, err := tx.Raw(sql, srcValues...).Rows()
	if err != nil {
		kinit.LogError.Println(err)
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	columns, err := rows.Columns()
	if err != nil {
		kinit.LogError.Println(err)
		return
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		kinit.LogError.Println(err)
		return
	}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			scanArgs[i] = &values[i]
		}
		err = rows.Scan(scanArgs...)
		if err != nil {
			kinit.LogError.Println(err)
			return
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			v := values[i]
			dType := colTypes[i].DatabaseTypeName()

			if v == nil {
				switch dType {
				case "TINYINT", "SMALLINT", "MEDIUMINT", "INTEGER", "INT", "BIGINT", "FLOAT", "DOUBLE", "DECIMAL", "NUMERIC", "REAL":
					entry[col] = 0
				default:
					entry[col] = ""
				}
			} else {
				vType := reflect.TypeOf(v)
				vTypeStr := vType.String()

				switch dType {
				case "TINYINT", "SMALLINT", "MEDIUMINT", "INTEGER", "INT", "BIGINT":
					if vTypeStr == "[]uint8" {
						entry[col], _ = strconv.ParseInt(string(v.([]uint8)), 10, 64)
					} else {
						entry[col] = v
					}
				case "FLOAT", "DOUBLE", "DECIMAL", "NUMERIC", "REAL":
					if vTypeStr == "[]uint8" {
						entry[col], _ = strconv.ParseFloat(string(v.([]uint8)), 64)
					} else {
						entry[col] = v
					}
				case "BIT":
					if vTypeStr == "[]uint8" {
						entry[col] = kutils.Bytes2BitsString(v.([]uint8))
					} else {
						entry[col] = v
					}
				case "BINARY", "VARBINARY":
					if vTypeStr == "[]uint8" {
						entry[col] = hex.EncodeToString(v.([]uint8))
					} else {
						entry[col] = v
					}
				case "DATE", "TIME", "YEAR", "DATETIME", "TIMESTAMP":
					if vTypeStr == "[]uint8" {
						entry[col] = string(v.([]uint8))
					} else if vTypeStr == "time.Time" {
						entry[col] = v.(time.Time).Format("2006-01-02 15:04:05")
					} else {
						entry[col] = v
					}
				default:
					if vTypeStr == "[]uint8" {
						entry[col] = string(v.([]uint8))
					} else {
						entry[col] = v
					}
				}
			}

		}
		records = append(records, entry)
	}

	return
}

func InsertOne(tx *jgorm.DB, name string, dbId int64, tableName string, insertMap map[string]interface{}) (int64, error) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect(name, dbId)
	}
	var id []int64

	insertTime := time.Now().Format("2006-01-02 15:04:05")
	insertMap["created_at"] = insertTime
	insertMap["updated_at"] = insertTime

	countParam := len(insertMap)
	field := make([]string, 0)
	value := make([]interface{}, 0)
	for k, _ := range insertMap {
		field = append(field, fmt.Sprintf("`%s`", k))
		value = append(value, insertMap[k])
	}

	fieldStr := strings.Join(field, ",")
	sqlStr := strings.TrimRight(strings.Repeat("?,", countParam), ",")
	sql := fmt.Sprintf("insert into `%s`(%s) values(%s)", tableName, fieldStr, sqlStr)

	if err := tx.Exec(sql, value...).Error; err != nil {
		kinit.LogError.Println(err)
		return 0, err
	}

	tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id)
	return id[0], nil
}

func UpdateOne(tx *jgorm.DB, name string, dbId int64, tableName string, whereMap map[string]interface{}, updateMap map[string]interface{}) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect(name, dbId)
	}

	tx = tx.Table(tableName)
	updateMap["updated_at"] = time.Now().Format("2006-01-02 15:04:05")

	if len(whereMap) > 0 {
		for k, v := range whereMap {
			tx = tx.Where(k+"=?", v)
		}
	}
	if err := tx.Updates(updateMap).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}

	return nil
}

func DeleteOne(tx *jgorm.DB, name string, dbId int64, tableName string, whereMap map[string]interface{}) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect(name, dbId)
	}

	tx = tx.Table(tableName)

	if len(whereMap) > 0 {
		for k, v := range whereMap {
			tx = tx.Where(k+"=?", v)
		}
	}

	if err := tx.Delete(nil).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}

	return nil
}
