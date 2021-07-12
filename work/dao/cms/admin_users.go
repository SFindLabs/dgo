package cms

import (
	kinit "dgo/work/base/initialize"
	"fmt"
	jgorm "github.com/jinzhu/gorm"
	"time"
)

var CmsAdminUsersObj CmsAdminUsers

type CmsAdminUsers struct {
	ID        int64  `gorm:"primary_key" json:"-"`
	Name      string `gorm:"column:name" json:"name"`
	Password  string `gorm:"column:password" json:"password"`
	Avatar    string `gorm:"column:avatar" json:"avatar"`
	Status    int64  `gorm:"column:status" json:"status"`
	LoginIp   string `gorm:"column:login_ip" json:"login_ip"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
	LoginAt   string `gorm:"column:login_at" json:"login_at"`
}

func (CmsAdminUsers) TableName() string {
	return "cms_admin_users"
}

func (CmsAdminUsers) Insert(tx *jgorm.DB, name string, password string, avatar string, loginIp string) (CmsAdminUsers, error) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	insertTime := time.Now().Format("2006-01-02 15:04:05")
	obj := CmsAdminUsers{
		Name:      name,
		Password:  password,
		Avatar:    avatar,
		Status:    1,
		LoginIp:   loginIp,
		CreatedAt: insertTime,
		LoginAt:   insertTime,
	}
	if err := tx.Create(&obj).Error; err != nil {
		kinit.LogError.Println(err)
		return obj, err
	}
	return obj, nil
}

func (CmsAdminUsers) GetById(tx *jgorm.DB, id int64) CmsAdminUsers {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminUsers
	tx.Where("id=? ", id).First(&objs)
	return objs
}

func (CmsAdminUsers) GetStatusById(tx *jgorm.DB, id int64) []int64 {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var status []int64
	tx.Model(CmsAdminUsers{}).Where("id=? ", id).Pluck("status", &status)
	return status
}

func (CmsAdminUsers) GetAllOrByName(tx *jgorm.DB, name string) []CmsAdminUsers {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs []CmsAdminUsers
	if name != "" {
		tx.Where("name=?", name).Find(&objs)
	} else {
		tx.Find(&objs)
	}
	return objs
}

func (CmsAdminUsers) GetByName(tx *jgorm.DB, name string) CmsAdminUsers {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminUsers
	tx.Where("name=? ", name).First(&objs)
	return objs
}

func (CmsAdminUsers) UpdatePicById(tx *jgorm.DB, id int64, avatar string) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	if err := tx.Model(CmsAdminUsers{}).Where("id=?", id).Updates(map[string]interface{}{"avatar": avatar}).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminUsers) UpdatePasswordById(tx *jgorm.DB, id int64, password string) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	if err := tx.Model(CmsAdminUsers{}).Where("id=?", id).Updates(map[string]interface{}{"password": password}).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminUsers) UpdateStatusById(tx *jgorm.DB, id, status int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	updateMap := map[string]interface{}{
		"status": status,
	}
	if err := tx.Model(CmsAdminUsers{}).Where("id=?", id).Updates(updateMap).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminUsers) UpdateLoginInfoById(tx *jgorm.DB, id int64, ip string) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	updateMap := map[string]interface{}{
		"login_ip": ip,
		"login_at": time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := tx.Model(CmsAdminUsers{}).Where("id=?", id).Updates(updateMap).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminUsers) DeleteById(tx *jgorm.DB, id int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsAdminUsers
	if err := tx.Where("id=? ", id).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsAdminUsers) GetIdUserHasRole(tx *jgorm.DB, userId int64) RoleIdBlock {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var id RoleIdBlock
	sql := "select role_id from `cms_admin_users` as u left join " +
		"`cms_admin_user_has_roles` as r on u.id=r.admin_id where "
	sql += fmt.Sprintf(" u.id = %d; ", userId)

	tx.Raw(sql).Scan(&id)
	return id
}

func (CmsAdminUsers) CountByName(tx *jgorm.DB, name string) (count int) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	tx = tx.Model(CmsAdminUsers{})
	if name != "" {
		tx = tx.Where("name like ?", "%"+name+"%")
	}
	tx.Count(&count)
	return
}

func (CmsAdminUsers) GetByAllName(tx *jgorm.DB, count int, name string, page int64, pageSize int64) []CmsAdminUsers {
	var objs []CmsAdminUsers
	if count == 0 {
		return objs
	}
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	tx = tx.Model(CmsAdminUsers{})
	if name != "" {
		tx = tx.Where("name like ?", "%"+name+"%")
	}
	tx.Limit(pageSize).Offset((page - 1) * pageSize).Find(&objs)
	return objs
}
