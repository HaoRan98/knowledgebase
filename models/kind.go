package models

import (
	"errors"
	"log"
)

type Kind struct {
	ID uint   `gorm:"primary_key"`
	Mc string `json:"mc"`
}

func AddKind(kind *Kind) error {
	if err := db.Create(kind).Error; err != nil {
		return err
	}
	return nil
}

func EditKind(kind map[string]interface{}) error {

	if err := db.Model(&Kind{}).
		Where("id=?", kind["id"]).Updates(kind).Error; err != nil {
		return err
	}
	return nil
}

func GetKind(id uint) (error, bool) {

	var kind Kind

	if err := db.Model(&Kind{}).
		Where("id=?", id).First(&kind).Error; err != nil {
		return err, false
	}

	//log.Println(kind)

	if kind.Mc == "" {
		log.Println("该分类为空")
		return errors.New("该分类为空"), false
	}

	var topics []*Topic
	if err := db.Debug().Where("kind = ?", kind.Mc).Find(&topics).Error; err != nil {
		return err, false
	}

	if len(topics) > 0 {
		return nil, true
	}

	return nil, false

}

func DelKind(id uint) error {
	if err := db.Where("id=?", id).Delete(Kind{}).Error; err != nil {
		return err
	}
	return nil
}
func GetKinds() ([]*Kind, error) {
	var kinds []*Kind
	if err := db.Find(&kinds).Error; err != nil {
		return nil, err
	}
	return kinds, nil
}
func IsKindExist(mc string) bool {
	var kind Kind
	if err := db.Where("mc=?", mc).First(&kind).Error; err != nil {
		return false
	}
	if kind.ID > 0 {
		return true
	}
	return false
}

func Getkind(id uint) (Kind, error) {
	var kind Kind
	if err := db.Where("id = ?", id).First(&kind).Error; err != nil {
		return kind, err
	}
	return kind, nil
}
