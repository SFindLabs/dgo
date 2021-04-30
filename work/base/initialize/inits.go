package initialize

import (
	"dgo/framework/tools/conf"
	kgorm "dgo/framework/tools/db/gorm"
	kredis "dgo/framework/tools/db/redis"
	klog "dgo/framework/tools/log"
	kruntime "dgo/framework/tools/runtime"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"log"
	"runtime"
	"strconv"
	"strings"
)

const CONFIG_PATH = "conf/config.json"

var Conf *conf.JsonConf

var LogError *klog.LlogFile
var LogWarn *klog.LlogFile
var LogInfo *klog.LlogFile
var LogDebug *klog.LlogFile

var GorMMap = make(map[string]map[string]map[int64]map[int64]*gorm.DB)
var RedisPoolMap = make(map[string]map[string]map[int64]map[int64]*kredis.RedisPool)
var SnowflakeNode *snowflake.Node

func init() {
	var err error
	runtime.GOMAXPROCS(runtime.NumCPU()) //多核设置

	LogError = klog.Error
	LogWarn = klog.Warn
	LogInfo = klog.Info
	LogDebug = klog.Debug

	SnowflakeNode, err = snowflake.NewNode(10)
	if err != nil {
		log.Println(err)
	}
}

func InitConf(confFile string) {
	if confFile == "" {
		confFile = CONFIG_PATH
	}
	tmpConf, err := conf.NewJsonConf(confFile)
	if err != nil {
		log.Panic("conf parse fail:", err)
	}
	Conf = tmpConf
}

//isSplit 是否每日分割日志
func InitLog(logfile string, isSplit ...bool) {
	root, _ := kruntime.GetCurrentPath()
	dir := fmt.Sprintf("%s/log", root)
	is, _ := kruntime.PathExists(dir)
	if !is {
		kruntime.CreateDir(dir)
	}
	logFile := fmt.Sprintf("%s/%s", dir, logfile)
	LogError, _ = klog.NewLlogFile(fmt.Sprintf("%s_err.log", logFile), "[Error]", klog.LstdFlags|klog.Lshortfile, klog.LOG_LEVEL_ERROR, klog.LOG_LEVEL_ERROR, 50, isSplit...)
	// LogWarn, _ = klog.NewLlogFile("log/"+logfile+"_warn.log", "[Warn]", klog.LstdFlags|klog.Lshortfile, klog.LOG_LEVEL_ERROR, klog.LOG_LEVEL_ERROR, 50, isSplit...)
	// LogInfo, _ = klog.NewLlogFile("log/"+logfile+"_info.log", "[Info]", klog.LstdFlags|klog.Lshortfile, klog.LOG_LEVEL_ERROR, klog.LOG_LEVEL_ERROR, 50, isSplit...)
	// LogDebug, _ = klog.NewLlogFile("log/"+logfile+"_debug.log", "[Debug]", klog.LstdFlags|klog.Lshortfile, klog.LOG_LEVEL_ERROR, klog.LOG_LEVEL_ERROR, 50, isSplit...)
	LogWarn = LogError
	LogInfo = LogError
	LogDebug = LogError
}

func InitMysql() {
	mysqlUser, _ := Conf.GetString("mysql.user")
	mysqlPassword, _ := Conf.GetString("mysql.password")
	mysqlIp, _ := Conf.GetString("mysql.host")
	mysqlPort, _ := Conf.GetInt("mysql.port")
	mysqlPortStr := fmt.Sprintf("%d", mysqlPort)
	mysqlDb, _ := Conf.GetString("mysql.db")
	isTest, _ := Conf.GetInt("server.is_test")

	//可选参数
	maxLifeTime, _ := Conf.GetInt("mysql.max_life_time")
	maxLifeTimeStr := fmt.Sprintf("%d", maxLifeTime)
	maxIdleConn, _ := Conf.GetInt("mysql.max_idle_conn")
	maxIdleConnStr := fmt.Sprintf("%d", maxIdleConn)
	maxOpenConn, _ := Conf.GetInt("mysql.max_open_conn")
	maxOpenConnStr := fmt.Sprintf("%d", maxOpenConn)

	defaultConnect, err := kgorm.NewGorm(mysqlUser, mysqlPassword, mysqlIp, mysqlPortStr, mysqlDb, maxLifeTimeStr, maxIdleConnStr, maxOpenConnStr)
	if err != nil {
		log.Panic(err)
	}
	if isTest == 1 {
		defaultConnect.LogMode(true)
	}
	GorMMap["default"] = map[string]map[int64]map[int64]*gorm.DB{"write": {0: {0: defaultConnect}}}

	result, err := Conf.GetArrayMap("mysql", "host", "port", "user", "password", "db", "max_life_time", "max_idle_conn", "max_open_conn")
	if err != nil {
		log.Panic(err)
	}
	for _, child := range result {
		for key, cc := range child {
			user, pwd, host, port, db, life, idle, open := "", "", "", "", "", "0", "0", "0"
			if _, ok := cc["user"]; !ok {
				user = mysqlUser
			} else {
				user = cc["user"]
			}
			if _, ok := cc["password"]; !ok {
				pwd = mysqlPassword
			} else {
				pwd = cc["password"]
			}
			if _, ok := cc["host"]; !ok {
				host = mysqlIp
			} else {
				host = cc["host"]
			}
			if _, ok := cc["port"]; !ok {
				port = mysqlPortStr
			} else {
				port = cc["port"]
			}
			if _, ok := cc["db"]; !ok {
				db = mysqlDb
			} else {
				db = cc["db"]
			}

			//可选参数
			if _, ok := cc["max_life_time"]; !ok {
				life = maxLifeTimeStr
			} else {
				life = cc["max_life_time"]
			}
			if _, ok := cc["max_idle_conn"]; !ok {
				idle = maxIdleConnStr
			} else {
				idle = cc["max_idle_conn"]
			}
			if _, ok := cc["max_open_conn"]; !ok {
				open = maxOpenConnStr
			} else {
				open = cc["max_open_conn"]
			}

			tmpConnect, err := kgorm.NewGorm(user, pwd, host, port, db, life, idle, open)
			if err != nil {
				log.Panic(err)
			}
			if isTest == 1 {
				tmpConnect.LogMode(true)
			}

			gMapKey := strings.Split(key, "_")
			gMapKey2, _ := strconv.ParseInt(gMapKey[2], 10, 64)
			gMapKey3, _ := strconv.ParseInt(gMapKey[3], 10, 64)

			if val, ok := GorMMap[gMapKey[0]]; ok {
				if childVal, ok := val[gMapKey[1]]; ok {
					if childIndex, ok := childVal[gMapKey2]; ok {
						childIndex[gMapKey3] = tmpConnect
					} else {
						childVal[gMapKey2] = map[int64]*gorm.DB{gMapKey3: tmpConnect}
					}
				} else {
					val[gMapKey[1]] = map[int64]map[int64]*gorm.DB{gMapKey2: {gMapKey3: tmpConnect}}
				}
			} else {
				GorMMap[gMapKey[0]] = map[string]map[int64]map[int64]*gorm.DB{gMapKey[1]: {gMapKey2: {gMapKey3: tmpConnect}}}
			}
		}
	}
}

func InitRedis() {
	redisHost, _ := Conf.GetString("redis.host")
	redisPort, _ := Conf.GetInt("redis.port")
	redisPortStr := fmt.Sprintf("%d", redisPort)
	redisAuth, _ := Conf.GetString("redis.auth")

	//可选参数
	redisDb, _ := Conf.GetInt("redis.db")
	redisDbStr := fmt.Sprintf("%d", redisDb)
	maxIdle, _ := Conf.GetInt("redis.max_idle")
	maxIdleStr := fmt.Sprintf("%d", maxIdle)
	maxActive, _ := Conf.GetInt("redis.max_active")
	maxActiveStr := fmt.Sprintf("%d", maxActive)
	idleTimeout, _ := Conf.GetInt("redis.idle_timeout")
	idleTimeoutStr := fmt.Sprintf("%d", idleTimeout)

	defaultRedisPool, err := kredis.NewRedisPool(redisHost, redisPortStr, redisAuth, redisDbStr, maxIdleStr, maxActiveStr, idleTimeoutStr)
	if err != nil {
		log.Panic(err)
	}
	RedisPoolMap["default"] = map[string]map[int64]map[int64]*kredis.RedisPool{"write": {0: {0: defaultRedisPool}}}

	result, err := Conf.GetArrayMap("redis", "host", "port", "auth", "db", "max_idle", "max_active", "idle_timeout")
	if err != nil {
		log.Panic(err)
	}
	for _, child := range result {
		for key, cc := range child {
			host, port, auth, db, idle, active, timeout := "", "", "", "0", "0", "0", "0"
			if _, ok := cc["host"]; !ok {
				host = redisHost
			} else {
				host = cc["host"]
			}
			if _, ok := cc["port"]; !ok {
				port = redisPortStr
			} else {
				port = cc["port"]
			}
			if _, ok := cc["auth"]; !ok {
				auth = redisAuth
			} else {
				auth = cc["auth"]
			}

			//可选参数
			if _, ok := cc["db"]; !ok {
				db = redisDbStr
			} else {
				db = cc["db"]
			}
			if _, ok := cc["max_idle"]; !ok {
				idle = maxIdleStr
			} else {
				idle = cc["max_idle"]
			}
			if _, ok := cc["max_active"]; !ok {
				active = maxActiveStr
			} else {
				active = cc["max_active"]
			}
			if _, ok := cc["idle_timeout"]; !ok {
				timeout = idleTimeoutStr
			} else {
				timeout = cc["idle_timeout"]
			}
			tmpRedisPool, err := kredis.NewRedisPool(host, port, auth, db, idle, active, timeout)
			if err != nil {
				log.Panic(err)
			}

			rMapKey := strings.Split(key, "_")
			rMapKey2, _ := strconv.ParseInt(rMapKey[2], 10, 64)
			rMapKey3, _ := strconv.ParseInt(rMapKey[3], 10, 64)

			if val, ok := RedisPoolMap[rMapKey[0]]; ok {
				if childVal, ok := val[rMapKey[1]]; ok {
					if childIndex, ok := childVal[rMapKey2]; ok {
						childIndex[rMapKey3] = tmpRedisPool
					} else {
						childVal[rMapKey2] = map[int64]*kredis.RedisPool{rMapKey3: tmpRedisPool}
					}
				} else {
					val[rMapKey[1]] = map[int64]map[int64]*kredis.RedisPool{rMapKey2: {rMapKey3: tmpRedisPool}}
				}
			} else {
				RedisPoolMap[rMapKey[0]] = map[string]map[int64]map[int64]*kredis.RedisPool{rMapKey[1]: {rMapKey2: {rMapKey3: tmpRedisPool}}}
			}
		}
	}
}

// ====================================获取连接池连接===============================================

/**
 * name 模块数据库名
 * dbId 第几个库(用于分库)
 */
func GetMysqlConnect(name string, dbIds ...int64) (*gorm.DB, error) {
	var dbId int64 = 0
	if name == "" {
		name = "default"
	}
	switch len(dbIds) {
	case 1:
		dbId = dbIds[0]
	}
	if connect, ok := (((GorMMap[name])["write"])[dbId])[0]; ok {
		return connect, nil
	}
	err := fmt.Sprintf("the %s %d mysql write connect not exist", name, dbId)
	LogError.Println(err)
	return nil, errors.New(err)
}

/**
 * name 模块数据库名
 * dbId 第几个库(用于分库)
 * readIndex  第几个从库(用于读写分离)
 */
func GetMysqlReadConnect(name string, dbIds ...int64) (*gorm.DB, error) {
	if name == "" {
		name = "default"
	}
	var dbId, readIndex int64 = 0, 0
	switch len(dbIds) {
	case 1:
		dbId = dbIds[0]
	case 2:
		dbId = dbIds[0]
		readIndex = dbIds[1]
	}

	if connect, ok := (((GorMMap[name])["read"])[dbId])[readIndex]; ok {
		return connect, nil
	}
	LogError.Println(fmt.Sprintf("the %s %d mysql read connect not exist", name, dbId))
	if connect, ok := (((GorMMap[name])["write"])[dbId])[0]; ok {
		return connect, nil
	}
	err := fmt.Sprintf("the %s %d mysql connect not exist", name, dbId)
	LogError.Println(err)
	return nil, errors.New(err)
}

/**
 * name 模块数据库名
 * dbId 第几个库(用于分库)
 */
func GetRedisConnect(name string, dbIds ...int64) (*kredis.RedisPool, error) {
	var dbId int64 = 0
	if name == "" {
		name = "default"
	}
	switch len(dbIds) {
	case 1:
		dbId = dbIds[0]
	}
	if connect, ok := (((RedisPoolMap[name])["write"])[dbId])[0]; ok {
		return connect, nil
	}
	err := fmt.Sprintf("the %s %d redis pool connect not exist", name, dbId)
	LogError.Println(err)
	return nil, errors.New(err)
}

/**
 * name 模块数据库名
 * dbId 第几个库(用于分库)
 * readIndex  第几个从库(用于读写分离)
 */
func GetRedisReadConnect(name string, dbIds ...int64) (*kredis.RedisPool, error) {
	if name == "" {
		name = "default"
	}
	var dbId, readIndex int64 = 0, 0
	switch len(dbIds) {
	case 1:
		dbId = dbIds[0]
	case 2:
		dbId = dbIds[0]
		readIndex = dbIds[1]
	}
	if connect, ok := (((RedisPoolMap[name])["read"])[dbId])[readIndex]; ok {
		return connect, nil
	}
	LogError.Println(fmt.Sprintf("the %s %d redis read pool connect not exist", name, dbId))
	if connect, ok := (((RedisPoolMap[name])["write"])[dbId])[0]; ok {
		return connect, nil
	}
	err := fmt.Sprintf("the %s %d redis pool connect not exist", name, dbId)
	LogError.Println(err)
	return nil, errors.New(err)
}

//-----------------------------------------取余-----------------------------------------------
/**
 * name 模块数据库名
 * dbId 第几个库(用于分库)
 */
func GetMysqlDivideConnect(name string, dbIds ...int64) (*gorm.DB, error) {
	var dbId int64 = 0
	if name == "" {
		name = "default"
	}
	switch len(dbIds) {
	case 1:
		dbId = dbIds[0]
	}
	totalDb := int64(len((GorMMap[name])["write"]))
	if totalDb == 0 {
		return nil, errors.New(fmt.Sprintf("the %s mysql write connect not exist", name))
	}
	dbIndex := dbId % totalDb
	if connect, ok := (((GorMMap[name])["write"])[dbIndex])[0]; ok {
		return connect, nil
	}
	err := fmt.Sprintf("the %s %d mysql write connect not exist", name, dbIndex)
	LogError.Println(err)
	return nil, errors.New(err)
}

/**
 * name 模块数据库名
 * dbId 第几个库(用于分库)
 * readIndex  第几个从库(用于读写分离)
 */
func GetMysqlReadDivideConnect(name string, dbIds ...int64) (*gorm.DB, error) {
	if name == "" {
		name = "default"
	}
	var dbId, readIndex, dbIndex int64 = 0, 0, 0
	switch len(dbIds) {
	case 1:
		dbId = dbIds[0]
	case 2:
		dbId = dbIds[0]
		readIndex = dbIds[1]
	}
	totalDb := int64(len((GorMMap[name])["read"]))
	if totalDb > 0 {
		dbIndex = dbId % totalDb
		if connect, ok := (((GorMMap[name])["read"])[dbIndex])[readIndex]; ok {
			return connect, nil
		}
		LogError.Println(fmt.Sprintf("the %s %d mysql read connect not exist", name, dbIndex))
	} else {
		LogError.Println(fmt.Sprintf("the %s mysql read connect not exist", name))
	}

	totalDb = int64(len((GorMMap[name])["write"]))
	if totalDb == 0 {
		return nil, errors.New(fmt.Sprintf("the %s mysql connect not exist", name))
	}
	dbIndex = dbId % totalDb
	if connect, ok := (((GorMMap[name])["write"])[dbIndex])[0]; ok {
		return connect, nil
	}
	err := fmt.Sprintf("the %s %d mysql connect not exist", name, dbIndex)
	LogError.Println(err)
	return nil, errors.New(err)
}

/**
 * name 模块数据库名
 * dbId 第几个库(用于分库)
 */
func GetRedisDivideConnect(name string, dbIds ...int64) (*kredis.RedisPool, error) {
	var dbId int64 = 0
	if name == "" {
		name = "default"
	}
	switch len(dbIds) {
	case 1:
		dbId = dbIds[0]
	}
	totalRedis := int64(len((RedisPoolMap[name])["write"]))
	if totalRedis == 0 {
		return nil, errors.New(fmt.Sprintf("the %s redis pool connect not exist", name))
	}
	dbIndex := dbId % totalRedis
	if connect, ok := (((RedisPoolMap[name])["write"])[dbIndex])[0]; ok {
		return connect, nil
	}
	err := fmt.Sprintf("the %s %d redis pool connect not exist", name, dbIndex)
	LogError.Println(err)
	return nil, errors.New(err)
}

/**
 * name 模块数据库名
 * dbId 第几个库(用于分库)
 * readIndex  第几个从库(用于读写分离)
 */
func GetRedisReadDivideConnect(name string, dbIds ...int64) (*kredis.RedisPool, error) {
	if name == "" {
		name = "default"
	}
	var dbId, readIndex, dbIndex int64 = 0, 0, 0
	switch len(dbIds) {
	case 1:
		dbId = dbIds[0]
	case 2:
		dbId = dbIds[0]
		readIndex = dbIds[1]
	}
	totalRedis := int64(len((RedisPoolMap[name])["read"]))
	if totalRedis > 0 {
		dbIndex = dbId % totalRedis
		if connect, ok := (((RedisPoolMap[name])["read"])[dbIndex])[readIndex]; ok {
			return connect, nil
		}
		LogError.Println(fmt.Sprintf("the %s %d redis read pool connect not exist", name, dbIndex))
	} else {
		LogError.Println(fmt.Sprintf("the %s redis read pool connect not exist", name))
	}

	totalRedis = int64(len((RedisPoolMap[name])["write"]))
	if totalRedis == 0 {
		return nil, errors.New(fmt.Sprintf("the %s redis pool connect not exist", name))
	}
	dbIndex = dbId % totalRedis
	if connect, ok := (((RedisPoolMap[name])["write"])[dbIndex])[0]; ok {
		return connect, nil
	}
	err := fmt.Sprintf("the %s %d redis pool connect not exist", name, dbIndex)
	LogError.Println(err)
	return nil, errors.New(err)
}
