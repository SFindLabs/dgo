package cms

import (
	kinit "dgo/work/base/initialize"
	"bytes"
	"fmt"
	jgorm "github.com/jinzhu/gorm"
	"time"
)

var CmsAdminRoleHasPermissionsObj CmsAdminRoleHasPermissions

type CmsAdminRoleHasPermissions struct {
	ID           int64  `gorm:"primary_key" json:"-"`
	RoleId       int64  `gorm:"column:role_id" json:"role_id"`
	PermissionId int64  `gorm:"column:permission_id" json:"permission_id"`
	CreatedAt    string `gorm:"column:created_at" json:"created_at"`
}

func (CmsAdminRoleHasPermissions) TableName() string {
	return "cms_admin_role_has_permissions"
}

func (CmsAdminRoleHasPermissions) InsertByRoleId(tx *jgorm.DB, roleId int64, permissionId []int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	insertTime := time.Now().Format("2006-01-02 15:04:05")
	sql := "insert into `cms_admin_role_has_permissions` (`role_id`,`permission_id`,`created_at`) values "

	var buffer bytes.Buffer
	buffer.WriteString(sql)
	for k, v := range permissionId {
		if len(permissionId)-1 == k {
			buffer.WriteString(fmt.Sprintf("(%d, %d, '%s');", roleId, v, insertTime))
		} else {
			buffer.WriteString(fmt.Sprintf("(%d, %d, '%s'),", roleId, v, insertTime))
		}
	}
	return tx.Exec(buffer.String()).Error
}

func (CmsAdminRoleHasPermissions) GetMapByRoleId(tx *jgorm.DB, roleId int64) map[int64]CmsAdminRoleHasPermissions {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs []CmsAdminRoleHasPermissions
	var dict map[int64]CmsAdminRoleHasPermissions
	tx.Where("role_id=? ", roleId).Find(&objs)
	dict = make(map[int64]CmsAdminRoleHasPermissions)
	count := len(objs)
	for i := 0; i < count; i++ {
		dict[objs[i].PermissionId] = objs[i]
	}
	return dict
}

func (CmsAdminRoleHasPermissions) DeleteByRoleIdPermissionIds(tx *jgorm.DB, roleId int64, ids []int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminRoleHasPermissions
	if err := tx.Where("role_id=? ", roleId).Where("permission_id in (?) ", ids).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminRoleHasPermissions) DeleteByRoleId(tx *jgorm.DB, roleId int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminRoleHasPermissions
	if err := tx.Where("role_id=? ", roleId).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminRoleHasPermissions) DeleteByPermissionIds(tx *jgorm.DB, ids []int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminRoleHasPermissions
	if err := tx.Where("permission_id in (?) ", ids).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}
