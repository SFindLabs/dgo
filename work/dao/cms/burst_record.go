package cms

import (
	kinit "dgo/work/base/initialize"
	jgorm "github.com/jinzhu/gorm"
	"time"
)

var CmsBurstRecordObj CmsBurstRecord

type CmsBurstRecord struct {
	ID             int64  `gorm:"primary_key" json:"-"`
	Uid            int64  `gorm:"column:uid" json:"uid"`
	TempFolderName string `gorm:"column:temp_folder_name" json:"temp_folder_name"`
	FileName       string `gorm:"column:file_name" json:"file_name"`
	FileTotalSize  string `gorm:"column:file_total_size" json:"file_total_size"`
	BurstCount     int64  `gorm:"column:burst_count" json:"burst_count"`
	BurstTotal     int64  `gorm:"column:burst_total" json:"burst_total"`
	CreatedAt      string `gorm:"column:created_at" json:"created_at"`
}

func (CmsBurstRecord) TableName() string {
	return "cms_burst_record"
}

func (CmsBurstRecord) Insert(tx *jgorm.DB, uid int64, tempFolderName string, fileName string, fileTotalSize string, burstTotal int64) (CmsBurstRecord, error) {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	obj := CmsBurstRecord{
		Uid:            uid,
		TempFolderName: tempFolderName,
		FileName:       fileName,
		FileTotalSize:  fileTotalSize,
		BurstTotal:     burstTotal,
		CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := tx.Create(&obj).Error; err != nil {
		kinit.LogError.Println(err)
		return obj, err
	}
	return obj, nil
}

func (CmsBurstRecord) GetByUidAndFilename(tx *jgorm.DB, uid int64, fileName string) CmsBurstRecord {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsBurstRecord
	tx.Where("uid=? ", uid).Where("file_name = ?", fileName).First(&objs)
	return objs
}

func (CmsBurstRecord) UpdateById(tx *jgorm.DB, id int64, burstCount int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}

	if err := tx.Model(CmsBurstRecord{}).Where("id=?", id).Updates(map[string]interface{}{"burst_count": burstCount}).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}

func (CmsBurstRecord) DeleteById(tx *jgorm.DB, id int64) error {
	if tx == nil {
		tx, _ = kinit.GetMysqlConnect("")
	}
	var objs CmsBurstRecord
	if err := tx.Where("id = ? ", id).Delete(objs).Error; err != nil {
		kinit.LogError.Println(err)
		return err
	}
	return nil
}
