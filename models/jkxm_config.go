package models

import (
	"time"
)

type Config struct {
	XmId       string    `json:"xm_id" gorm:"primary_key;column:XM_ID"`
	XmMc       string    `json:"xm_mc" gorm:"column:XM_MC;COMMENT:'指标名称'"`
	XmSql      string    `json:"xm_sql" gorm:"column:XM_SQL;COMMENT:'执行SQL'"`
	InsertDate time.Time `json:"insert_date"  gorm:"column:INSERT_DATE;COMMENT:'插入时间'"`
}

func GetConfigSql(xmMc string) (*Config, error) {
	var pd Config
	err := db.Table("klib.CONFIG").Where("XM_MC=?", xmMc).First(&pd).Error
	if err != nil {
		return nil, err
	}
	return &pd, nil
}
