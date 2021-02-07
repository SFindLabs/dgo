package gorm

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

/**
 * option
		maxLifeTime  连接可以重用的最长时间
		maxIdleConns 空闲连接池中的最大连接数
		maxOpenConns 数据库的最大打开连接数
*/
func NewGorm(user, password, ip, port, mysqlDb string, option ...string) (*gorm.DB, error) {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		ip,
		port,
		mysqlDb,
	)
	gormDB, err := gorm.Open("mysql", addr)
	if err != nil {
		log.Println("gorm open fail:", err)
		return nil, err
	}

	initMaxLifeTime := 60

	countOption := len(option)
	if countOption == 1 || countOption == 2 || countOption == 3 {
		lifeTime, err := strconv.Atoi(option[0])
		if err != nil {
			log.Println("max life time convert fail:", err)
			return nil, err
		}
		if lifeTime > 0 {
			initMaxLifeTime = lifeTime
		}
	}

	// SetMaxIdleConns设置空闲连接池中的最大连接数。
	if countOption == 2 || countOption == 3 {
		idleConn, err := strconv.Atoi(option[1])
		if err != nil {
			log.Println("max idle connection convert fail:", err)
			return nil, err
		}
		if idleConn > 0 {
			gormDB.DB().SetMaxIdleConns(idleConn)
		}
	}

	// SetMaxOpenConns设置到数据库的最大打开连接数。
	if countOption == 3 {
		openConn, err := strconv.Atoi(option[2])
		if err != nil {
			log.Println("max open connection convert fail:", err)
			return nil, err
		}
		if openConn > 0 {
			gormDB.DB().SetMaxOpenConns(openConn)
		}
	}

	//SetConnMaxLifetime设置连接可以重用的最长时间
	gormDB.DB().SetConnMaxLifetime(time.Duration(initMaxLifeTime) * time.Second)

	return gormDB, nil
}
