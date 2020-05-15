package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

//回帖
type Reply struct {
	ID        string     `json:"id" gorm:"primary_key"`
	DeletedAt *time.Time `sql:"index"`
	TopicID   string     `json:"topic_id"`
	Floor     int        `json:"floor" gorm:"COMMENT:'楼层'"`
	Content   string     `json:"content" gorm:"COMMENT:'回帖内容';size:65535"`
	Author    string     `json:"author"`
	Account   string     `json:"account"`
	Deptname  string     `json:"deptname"`
	Uptime    string     `json:"uptime"`
	Agree     int        `json:"agree" gorm:"COMMENT:'赞同数';default:'0'"`
	Accept    bool       `json:"accept" gorm:"COMMENT:'采纳标志';default:'0'"`
	Comments  []Comment
}

func CreateReply(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}
func GenReplyFloor(topicId string) (int, error) {
	var reply Reply
	err := db.Where("topic_id=?", topicId).Order("floor desc").
		Limit(1).First(&reply).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return 1, nil
	}
	return reply.Floor + 1, nil
}
func GetReply(id string) (*Reply, error) {
	var reply Reply
	if err := db.Where("id=?", id).First(&reply).Error; err != nil {
		return nil, err
	}
	return &reply, nil
}
func EditReply(reply *Reply) error {
	if err := db.Model(&Reply{}).Where("id=?", reply.ID).Updates(reply).Error; err != nil {
		return err
	}
	return nil
}
func ReplyAgree(id string) error {
	s := fmt.Sprintf("update reply set agree=agree+1 where id='%s'", id)
	if err := db.Exec(s).Error; err != nil {
		return err
	}
	return nil
}
func DelReply(id string) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	if err := tx.Table("reply").
		Where("id=?", id).Delete(Reply{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Table("comment").
		Where("reply_id=?", id).Delete(Comment{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func GetReplies(topicId string, pageNo, pageSize int) ([]*Reply, error) {
	var replies []*Reply
	if err := db.
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("comment.floor")
		}).
		Where("topic_id like ?", "%"+topicId+"%").
		Order("accept desc,agree desc,uptime desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&replies).Error; err != nil {
		return nil, err
	}
	if len(replies) > 0 {
		return replies, nil
	}
	return nil, nil
}
func GetRepliesCnt(topicId string) (cnt int) {
	if err := db.Table("reply").
		Where("topic_id like ?", "%"+topicId+"%").Count(&cnt).Error; err != nil {
		cnt = 0
	}
	return cnt
}
