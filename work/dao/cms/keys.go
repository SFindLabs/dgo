package cms

import (
	kinit "dgo/work/base/initialize"
	"errors"
	jgorm "github.com/jinzhu/gorm"
	"time"
)

var CmsKeysObj CmsKeys

type CmsKeys struct {
	ID        int64  `gorm:"primary_key" json:"-"`
	Name      string `gorm:"column:name" json:"name"`
	Keyx1     string `gorm:"column:keyx1" json:"keyx1"`
	Keyx2     string `gorm:"column:keyx2" json:"keyx2"`
	Valuex    string `gorm:"column:valuex" json:"valuex"`
	Status    int64  `gorm:"column:status" json:"status"`
	SortNum   int64  `gorm:"column:sort_num" json:"sort_num"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt string `gorm:"column:updated_at" json:"updated_at"`
}

func (CmsKeys) TableName() string {
	return "cms_keys"
}

func (CmsKeys) Insert(tx *jgorm.DB, name string, keyx1 string, keyx2 string, valuex string, status int64) (CmsKeys, error) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	obj := CmsKeys{
		Name:      name,
		Keyx1:     keyx1,
		Keyx2:     keyx2,
		Valuex:    valuex,
		Status:    status,
		CreatedAt: timeStr,
		UpdatedAt: timeStr,
	}
	if err := tx.Create(&obj).Error; err != nil {
		kinit.LogError.Println(err)
		return obj, err
	}
	return obj, nil
}

func (CmsKeys) GetAll(tx *jgorm.DB) (objs []CmsKeys) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	tx.Find(&objs)
	return objs
}

func (CmsKeys) CountByKey(tx *jgorm.DB, key string) int {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	count := 0
	tx = tx.Model(CmsKeys{})
	if key != "" {
		tx = tx.Where("keyx1 like ? or keyx2 like ?", "%"+key+"%", "%"+key+"%")
	}
	tx.Count(&count)
	return count
}

func (CmsKeys) GetAllByKey(tx *jgorm.DB, key string, page int64, pageSize int64) []CmsKeys {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs []CmsKeys
	if key != "" {
		tx = tx.Where("keyx1 like ? or keyx2 like ?", "%"+key+"%", "%"+key+"%")
	}
	tx.Order("sort_num desc,id desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&objs)
	return objs
}

func (CmsKeys) GetById(tx *jgorm.DB, id int64) (objs CmsKeys) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	tx.Where("id=? ", id).First(&objs)
	return objs
}

func (CmsKeys) GetByKeyx1(tx *jgorm.DB, keyx1 string) (objs CmsKeys) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	tx.Where("keyx1=? ", keyx1).First(&objs)
	return objs
}

func (CmsKeys) GetByKeyx2(tx *jgorm.DB, keyx2 string) (objs CmsKeys) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	tx.Where("keyx2=? ", keyx2).First(&objs)
	return objs
}

func (CmsKeys) GetByAllKey(tx *jgorm.DB, key1, key2 string) (objs CmsKeys) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	tx.Where("keyx1=?", key1).Where("keyx2=?", key2).First(&objs)
	return objs
}

func (CmsKeys) UpdateById(tx *jgorm.DB, id int64, updateMap map[string]interface{}) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	if err := tx.Model(CmsKeys{}).Where("id=?", id).Updates(updateMap).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsKeys) UpdateByKeyx1(tx *jgorm.DB, keyx1 string, updateMap map[string]interface{}) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	if err := tx.Model(CmsKeys{}).Where("keyx1=?", keyx1).Updates(updateMap).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsKeys) UpdateByKeyx2(tx *jgorm.DB, keyx2 string, updateMap map[string]interface{}) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	if err := tx.Model(CmsKeys{}).Where("keyx2=?", keyx2).Updates(updateMap).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsKeys) UpdateMustById(tx *jgorm.DB, id int64, updateMap map[string]interface{}) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	result := tx.Model(CmsKeys{}).Where("id=?", id).Updates(updateMap)
	if result.RowsAffected == 0 {
		errMsg := errors.New("CmsKeys UpdateById failed, rows is zero")
		if result.Error != nil {
			errMsg = result.Error
		}
		kinit.LogError.Println(errMsg)
		return errMsg
	}
	return nil
}

func (CmsKeys) UpdateMustByKeyx1(tx *jgorm.DB, keyx1 string, updateMap map[string]interface{}) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	result := tx.Model(CmsKeys{}).Where("keyx1=?", keyx1).Updates(updateMap)
	if result.RowsAffected == 0 {
		errMsg := errors.New("CmsKeys UpdateByKeyx1 failed, rows is zero")
		if result.Error != nil {
			errMsg = result.Error
		}
		kinit.LogError.Println(errMsg)
		return errMsg
	}
	return nil
}

func (CmsKeys) UpdateMustByKeyx2(tx *jgorm.DB, keyx2 string, updateMap map[string]interface{}) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	result := tx.Model(CmsKeys{}).Where("keyx2=?", keyx2).Updates(updateMap)
	if result.RowsAffected == 0 {
		errMsg := errors.New("CmsKeys UpdateByKeyx2 failed, rows is zero")
		if result.Error != nil {
			errMsg = result.Error
		}
		kinit.LogError.Println(errMsg)
		return errMsg
	}
	return nil
}

func (CmsKeys) DeleteById(tx *jgorm.DB, id int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsKeys
	if err := tx.Where("id=? ", id).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsKeys) DeleteByKeyx1(tx *jgorm.DB, keyx1 string) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsKeys
	if err := tx.Where("keyx1=? ", keyx1).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsKeys) DeleteByKeyx2(tx *jgorm.DB, keyx2 string) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsKeys
	if err := tx.Where("keyx2=? ", keyx2).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsKeys) DeleteMustById(tx *jgorm.DB, id int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsKeys
	result := tx.Where("id=? ", id).Delete(objs)
	if result.RowsAffected == 0 {
		errMsg := errors.New("CmsKeys DeleteById failed, rows is zero")
		if result.Error != nil {
			errMsg = result.Error
		}
		kinit.LogError.Println(errMsg)
		return errMsg
	}
	return nil
}

func (CmsKeys) DeleteMustByKeyx1(tx *jgorm.DB, keyx1 string) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsKeys
	result := tx.Where("keyx1=? ", keyx1).Delete(objs)
	if result.RowsAffected == 0 {
		errMsg := errors.New("CmsKeys DeleteByKeyx1 failed, rows is zero")
		if result.Error != nil {
			errMsg = result.Error
		}
		kinit.LogError.Println(errMsg)
		return errMsg
	}
	return nil
}

func (CmsKeys) DeleteMustByKeyx2(tx *jgorm.DB, keyx2 string) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsKeys
	result := tx.Where("keyx2=? ", keyx2).Delete(objs)
	if result.RowsAffected == 0 {
		errMsg := errors.New("CmsKeys DeleteByKeyx2 failed, rows is zero")
		if result.Error != nil {
			errMsg = result.Error
		}
		kinit.LogError.Println(errMsg)
		return errMsg
	}
	return nil
}

func (CmsKeys) DeleteByIds(tx *jgorm.DB, ids []int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsKeys
	if err := tx.Where("id in (?) ", ids).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}
