package cms

import (
	kinit "dgo/work/base/initialize"
	jgorm "github.com/jinzhu/gorm"
	"time"
)

var CmsAdminPermissionsObj CmsAdminPermissions

type PermissionMenu struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Pid     int64  `json:"pid"`
	Level   int64  `json:"level"`
	Open    bool   `json:"open"`
	Checked bool   `json:"checked"`
}

type CmsAdminPermissions struct {
	ID        int64  `gorm:"primary_key" json:"-"`
	Name      string `gorm:"column:name" json:"name"`
	Pid       int64  `gorm:"column:pid" json:"pid"`
	Path      string `gorm:"column:path" json:"path"`
	IsShow    int64  `gorm:"column:is_show" json:"is_show"`
	IsRecord  int64  `gorm:"column:is_record" json:"is_record"`
	IsModify  int64  `gorm:"column:is_modify" json:"is_modify"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt string `gorm:"column:updated_at" json:"updated_at"`
}

func (CmsAdminPermissions) TableName() string {
	return "cms_admin_permissions"
}

func (CmsAdminPermissions) Insert(tx *jgorm.DB, name string, pid int64, path string, isShow, isModify, isRecord int64) (CmsAdminPermissions, error) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	insertTime := time.Now().Format("2006-01-02 15:04:05")
	obj := CmsAdminPermissions{
		Name:      name,
		Pid:       pid,
		Path:      path,
		IsShow:    isShow,
		IsModify:  isModify,
		IsRecord:  isRecord,
		CreatedAt: insertTime,
		UpdatedAt: insertTime,
	}
	if err := tx.Create(&obj).Error; err != nil {
		kinit.LogError.Println(err)
		return obj, err
	}
	return obj, nil
}

func (CmsAdminPermissions) GetAll(tx *jgorm.DB) []CmsAdminPermissions {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs []CmsAdminPermissions
	tx.Find(&objs)
	return objs
}

func (CmsAdminPermissions) GetModifyAll(tx *jgorm.DB, isModify int64) []CmsAdminPermissions {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs []CmsAdminPermissions
	tx.Where("is_modify=?", isModify).Find(&objs)
	return objs
}

func (CmsAdminPermissions) GetById(tx *jgorm.DB, id int64) CmsAdminPermissions {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminPermissions
	tx.Where("id=? ", id).First(&objs)
	return objs
}

func (CmsAdminPermissions) GetByPidName(tx *jgorm.DB, name string, pid int64) CmsAdminPermissions {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminPermissions
	tx.Where("name=? ", name).Where("pid=? ", pid).First(&objs)
	return objs
}

func (CmsAdminPermissions) GetByLikePidName(tx *jgorm.DB, name string, pid int64) (ids []int64) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	tx = tx.Model(CmsAdminPermissions{})
	if name != "" {
		tx.Where("name like (?)", "%"+name+"%").Pluck("id", &ids)
	} else {
		tx.Where("pid = ?", pid).Pluck("id", &ids)
	}
	return
}

func (CmsAdminPermissions) GetByPath(tx *jgorm.DB, path string) CmsAdminPermissions {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminPermissions
	tx.Where("path=? ", path).First(&objs)
	return objs
}

func (CmsAdminPermissions) UpdateById(tx *jgorm.DB, id int64, name string, pid int64, path string, isShow, isModify, isRecord int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	updateMap := map[string]interface{}{
		"name":       name,
		"pid":        pid,
		"path":       path,
		"is_show":    isShow,
		"is_modify":  isModify,
		"is_record":  isRecord,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := tx.Model(CmsAdminPermissions{}).Where("id=?", id).Updates(updateMap).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminPermissions) DeleteByIds(tx *jgorm.DB, ids []int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminPermissions
	if err := tx.Where("id in (?) ", ids).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminPermissions) GetRoleIdHasPermissions(tx *jgorm.DB, roleId int64) []CmsAdminPermissions {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs []CmsAdminPermissions
	sql := "select * from `cms_admin_permissions` where is_show = ? and id in (" +
		"select permission_id from `cms_admin_role_has_permissions` where role_id = ?); "

	tx.Raw(sql, 1, roleId).Scan(&objs)
	return objs
}
