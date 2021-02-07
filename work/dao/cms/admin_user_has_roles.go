package cms

import (
	kinit "dgo/work/base/initialize"
	jgorm "github.com/jinzhu/gorm"
	"time"
)

var CmsAdminUserHasRolesObj CmsAdminUserHasRoles

type CmsAdminUserHasRoles struct {
	ID        int64  `gorm:"primary_key" json:"-"`
	AdminId   int64  `gorm:"column:admin_id" json:"admin_id"`
	RoleId    int64  `gorm:"column:role_id" json:"role_id"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}

func (CmsAdminUserHasRoles) TableName() string {
	return "cms_admin_user_has_roles"
}

func (CmsAdminUserHasRoles) Insert(tx *jgorm.DB, adminId, roleId int64) (CmsAdminUserHasRoles, error) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	obj := CmsAdminUserHasRoles{
		AdminId:   adminId,
		RoleId:    roleId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := tx.Create(&obj).Error; err != nil {
		kinit.LogError.Println(err)
		return obj, err
	}
	return obj, nil
}

func (CmsAdminUserHasRoles) GetByAdminId(tx *jgorm.DB, adminId int64) CmsAdminUserHasRoles {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminUserHasRoles
	tx.Where("admin_id=? ", adminId).First(&objs)
	return objs
}

func (CmsAdminUserHasRoles) UpdateByAdminId(tx *jgorm.DB, adminId int64, roleId int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	if err := tx.Model(CmsAdminUserHasRoles{}).Where("admin_id=?", adminId).Updates(map[string]interface{}{"role_id": roleId}).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminUserHasRoles) DeleteByRoleId(tx *jgorm.DB, roleId int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminUserHasRoles
	if err := tx.Where("role_id=? ", roleId).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminUserHasRoles) DeleteByAdminId(tx *jgorm.DB, adminId int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminUserHasRoles
	if err := tx.Where("admin_id=? ", adminId).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}
