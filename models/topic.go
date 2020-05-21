package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

//发帖
type Topic struct {
	ID       string `json:"id" gorm:"primary_key"`
	Del      bool   `json:"del" gorm:"COMMENT:'删除标志'"`
	Title    string `json:"title" gorm:"COMMENT:'标题'"`
	Content  string `json:"content" gorm:"COMMENT:'内容';size:65535"`
	FileName string `json:"file_name" gorm:"COMMENT:'文件原始名'"`
	FileUrl  string `json:"file_url" gorm:"COMMENT:'文件真实路径'"`
	Kind     string `json:"kind" gorm:"COMMENT:'类别'"`
	Author   string `json:"author"`
	Account  string `json:"account"`
	Deptname string `json:"deptname"`
	Uptime   string `json:"uptime"`
	Browse   int    `json:"browse" gorm:"COMMENT:'浏览量';default:'0'"`
	Hot      bool   `json:"hot" gorm:"COMMENT:'热门';default:'0'"`
	Top      bool   `json:"top" gorm:"COMMENT:'置顶标志';default:'0'"`
	Replys   []Reply
}

func CreateTopic(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}
func EditTopic(topic *Topic) error {
	if err := db.Model(&Topic{}).
		Where("id=?", topic.ID).Updates(topic).Error; err != nil {
		return err
	}
	return nil
}
func AddBrowse(id string) {
	s := fmt.Sprintf("update topic set browse=browse+1 where id='%s'", id)
	db.Exec(s)
}
func GetTopic(id string) (*Topic, error) {
	var topic Topic
	if err := db.Preload("Replys", func(db *gorm.DB) *gorm.DB {
		return db.Order("reply.accept,reply.floor")
	}).
		Where("id=?", id).First(&topic).Error; err != nil {
		return nil, err
	}
	if len(topic.ID) > 0 {
		return &topic, nil
	}
	return nil, nil
}
func GetTopics(account, kind string, pageNo, pageSize int) ([]*Topic, error) {
	var topics []*Topic
	if err := db.Preload("Replys", func(db *gorm.DB) *gorm.DB {
		return db.Order("reply.accept,reply.floor")
	}).
		Where("account like ?", "%"+account+"%").
		Where("kind like ?", "%"+kind+"%").
		Order("top desc,hot desc,uptime desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).
		Find(&topics).Error; err != nil {
		return nil, err
	}
	if len(topics) > 0 {
		return topics, nil
	}
	return nil, nil
}
func GetTopicsCnt(account string) (cnt int) {
	if err := db.Table("topic").
		Where("account like ?", "%"+account+"%").
		Count(&cnt).Error; err != nil {
		cnt = 0
	}
	return cnt
}
func GetTopicBrowseCnt() int {
	type wordCnt struct {
		C int
	}
	var cnt wordCnt
	if err := db.Raw(`SELECT sum(browse) c from topic`).Scan(&cnt).Error; err != nil {
		log.Println("Get Word Cnt err:", err)
		return 0
	}
	return cnt.C
}
