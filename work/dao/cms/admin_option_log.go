package cms

import (
	kinit "dgo/work/base/initialize"
	jgorm "github.com/jinzhu/gorm"
	"time"
)

var CmsAdminOptionLogObj CmsAdminOptionLog

type CmsAdminOptionLog struct {
	ID        int64  `gorm:"primary_key" json:"-"`
	UserName  string `gorm:"column:user_name" json:"user_name"`
	UserId    int64  `gorm:"column:user_id" json:"user_id"`
	Path      string `gorm:"column:path" json:"path"`
	Method    string `gorm:"column:method" json:"method"`
	Option    string `gorm:"column:option" json:"option"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}

func (CmsAdminOptionLog) TableName() string {
	return "cms_admin_option_log"
}

func (CmsAdminOptionLog) Insert(tx *jgorm.DB, userName string, userId int64, path string, method string, option string) (CmsAdminOptionLog, error) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	obj := CmsAdminOptionLog{
		UserName:  userName,
		UserId:    userId,
		Path:      path,
		Method:    method,
		Option:    option,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := tx.Create(&obj).Error; err != nil {
		kinit.LogError.Println(err)
		return obj, err
	}
	return obj, nil
}

func (CmsAdminOptionLog) CountByUserName(tx *jgorm.DB, userName string) (int, []int64) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	count := 0
	var adminIds []int64
	if userName != "" {
		tx.Model(CmsAdminUsers{}).Where("name like ?", "%"+userName+"%").Pluck("id", &adminIds)
		count = len(adminIds)
	} else {
		tx.Model(CmsAdminOptionLog{}).Count(&count)
	}
	return count, adminIds
}

func (CmsAdminOptionLog) GetByUserName(tx *jgorm.DB, count int, userId []int64, page int64, pageSize int64) []CmsAdminOptionLog {
	var objs []CmsAdminOptionLog
	if count == 0 {
		return objs
	}
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	var id []int64
	if len(userId) != 0 {
		tx = tx.Where("user_id in (?)", userId)
	}
	tx.Model(CmsAdminOptionLog{}).Limit(1).Offset((page-1)*pageSize).Order("id desc").Pluck("id", &id)
	if len(id) > 0 {
		tx.Where("id <= ?", id[0]).Order("id desc").Limit(pageSize).Find(&objs)
	}
	return objs
}

func (CmsAdminOptionLog) DeleteByIds(tx *jgorm.DB, ids []int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminOptionLog
	if err := tx.Where("id in (?) ", ids).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}