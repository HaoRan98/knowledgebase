package models

import "github.com/jinzhu/gorm"

//发帖
type Favorite struct {
	ID      string `json:"id" gorm:"primary_key"`
	Account string `json:"account"`
	Uptime  string `json:"uptime"`
	TopicID string `json:"topic_id"`
	Topic   Topic
}

func AddFavorite(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}
func CancelFavorite(topicId, account string) error {
	if err := db.Table("favorite").
		Where("topic_id=? and account=?", topicId, account).Delete(Favorite{}).Error; err != nil {
		return err
	}
	return nil
}
func GetFavorites(account string, pageNo, pageSize int) ([]*Favorite, error) {
	var favorites []*Favorite
	if err := db.Preload("Topic").
		Where("account=?", account).Limit(pageSize).Offset(pageSize * (pageNo - 1)).
		Find(&favorites).Error; err != nil {
		return nil, err
	}
	if len(favorites) > 0 {
		return favorites, nil
	}
	return nil, nil
}
func GetFavoritesCnt(account string) (cnt int) {
	if err := db.Table("favorite").
		Where("account=?", account).Count(&cnt).Error; err != nil {
		cnt = 0
	}
	return cnt
}
func GetCollector(topicId string) ([]*Favorite, error) {
	var favorites []*Favorite
	if err := db.Table("favorite").
		Where("topic_id=?", topicId).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}
func DelFavorite(id string) error {
	if err := db.Table("favorite").
		Where("id=?", id).Delete(Favorite{}).Error; err != nil {
		return err
	}
	return nil
}
func IsFavorite(topicId, account string) bool {
	var fav Favorite
	err := db.Where("topic_id=? and account=?", topicId, account).First(&fav).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}
	if err == gorm.ErrRecordNotFound {
		return false
	}
	if len(fav.ID) > 0 {
		return true
	}
	return false
}
