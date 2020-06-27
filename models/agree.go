package models

import "github.com/jinzhu/gorm"

type Agree struct {
	ID      string `json:"id" gorm:"primary_key"`
	Agreeid string `json:"agreeid" gorm:"COMMENT:'对应点赞的id'"`
	Account string `json:"account"`
	Uptime  string `json:"uptime"`
}

func Agreed(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func RemoveAgreed(agreeid, account string) error {
	if err := db.Where("agreeid=? and account=?", agreeid, account).
		Delete(Agree{}).Error; err != nil {
		return err
	}
	return nil
}

func IsAgreed(agreeid, account string) bool {
	var agree Agree
	err := db.Where("agreeid=? and account=?", agreeid, account).First(&agree).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}
	if err == gorm.ErrRecordNotFound {
		return false
	}
	if len(agree.ID) > 0 {
		return true
	}
	return false
}
