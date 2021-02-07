package cms

import (
	kinit "dgo/work/base/initialize"
	"fmt"
	jgorm "github.com/jinzhu/gorm"
	"time"
)

var CmsAdminRolesObj CmsAdminRoles

type RoleIdBlock struct {
	RoleId int64 `gorm:"column:role_id" json:"role_id"`
}

type CmsAdminRoles struct {
	ID        int64  `gorm:"primary_key" json:"-"`
	Name      string `gorm:"column:name" json:"name"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt string `gorm:"column:updated_at" json:"updated_at"`
}

func (CmsAdminRoles) TableName() string {
	return "cms_admin_roles"
}

func (CmsAdminRoles) Insert(tx *jgorm.DB, name string) (CmsAdminRoles, error) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	insertTime := time.Now().Format("2006-01-02 15:04:05")
	obj := CmsAdminRoles{
		Name:      name,
		CreatedAt: insertTime,
		UpdatedAt: insertTime,
	}
	if err := tx.Create(&obj).Error; err != nil {
		kinit.LogError.Println(err)
		return obj, err
	}
	return obj, nil
}

func (CmsAdminRoles) GetAll(tx *jgorm.DB) []CmsAdminRoles {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs []CmsAdminRoles
	tx.Find(&objs)
	return objs
}

func (CmsAdminRoles) GetIgnoreAll(tx *jgorm.DB, id int64) []CmsAdminRoles {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs []CmsAdminRoles
	tx.Where("id<>?", id).Find(&objs)
	return objs
}

func (CmsAdminRoles) GetById(tx *jgorm.DB, id int64) CmsAdminRoles {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminRoles
	tx.Where("id=? ", id).First(&objs)
	return objs
}

func (CmsAdminRoles) GetByName(tx *jgorm.DB, name string) CmsAdminRoles {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminRoles
	tx.Where("name=? ", name).First(&objs)
	return objs
}

func (CmsAdminRoles) UpdateById(tx *jgorm.DB, id int64, name string) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	updateMap := map[string]interface{}{
		"name":       name,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := tx.Model(CmsAdminRoles{}).Where("id=?", id).Updates(updateMap).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminRoles) DeleteById(tx *jgorm.DB, id int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminRoles
	if err := tx.Where("id=? ", id).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminRoles) GetIdRoleHasPermissions(tx *jgorm.DB, path string) []RoleIdBlock {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var ids []RoleIdBlock
	sql := "select role_id from `cms_admin_permissions` as ap right join " +
		"`cms_admin_role_has_permissions` as p on p.permission_id=ap.id where"
	sql += fmt.Sprintf(" ap.path = '%s'; ", path)

	tx.Raw(sql).Scan(&ids)
	return ids
}
