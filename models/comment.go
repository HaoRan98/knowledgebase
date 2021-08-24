package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

//评论回帖
type Comment struct {
	ID         string     `json:"id" gorm:"primary_key"`
	DeletedAt  *time.Time `sql:"index"`
	TopicID    string     `json:"topic_id"`
	ReplyID    string     `json:"reply_id"`
	Floor      int        `json:"floor" gorm:"COMMENT:'楼层'"`
	Content    string     `json:"content" gorm:"COMMENT:'评论内容';size:65535"`
	Author     string     `json:"author"`
	Account    string     `json:"account"`
	ToAccount  string     `json:"to_account"`
	ToDeptname string     `json:"to_deptname"`
	ToAuthor   string     `json:"to_author"`
	Deptname   string     `json:"deptname"`
	Createtime string     `json:"createtime"`
	Uptime     string     `json:"uptime"`
	Agree      int        `json:"agree" gorm:"COMMENT:'赞同数';default:'0'"`
	JGDM       string     `json:"jgdm"`
	JGMC       string     `json:"jgmc"`
}

func CreateComment(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func GenCommentFloor(replyId string) (int, error) {
	var comment Comment
	err := db.Where("reply_id=?", replyId).Order("floor desc").
		Limit(1).First(&comment).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return 1, nil
	}
	return comment.Floor + 1, nil
}
func GetCommentByID(id string) (*Comment, error) {
	var comment Comment
	if err := db.Where("id=?", id).First(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}
func EditComment(comment *Comment) error {
	if err := db.Model(&Comment{}).Where("id=?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}
	return nil
}
func CommentAgree(id string) error {
	s := fmt.Sprintf("update comment set agree=agree+1 where id='%s'", id)
	if err := db.Exec(s).Error; err != nil {
		return err
	}
	return nil
}
func RemoveCommentAgree(id string) error {
	s := fmt.Sprintf("update comment set agree=agree-1 where id='%s'", id)
	if err := db.Exec(s).Error; err != nil {
		return err
	}
	return nil
}
func DelComment(id string) error {
	if err := db.Table("comment").Where("id=?", id).Delete(Comment{}).Error; err != nil {
		return err
	}
	return nil
}
func GetComments(replyId string, pageNo, pageSize int) ([]*Comment, error) {
	var comments []*Comment
	if err := db.
		Where("reply_id like ?", "%"+replyId+"%").
		Order("agree desc,uptime desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&comments).Error; err != nil {
		return nil, err
	}
	if len(comments) > 0 {
		return comments, nil
	}
	return nil, nil
}
func GetCommentsCnt(replyId string) (cnt int) {
	if err := db.Table("comment").
		Where("reply_id like ?", "%"+replyId+"%").Count(&cnt).Error; err != nil {
		cnt = 0
	}
	return cnt
}

func RepGetComments(account string, pageNo, pageSize int) ([]*Comment, error) {
	var comments []*Comment
	if err := db.Debug().
		Where("account = ?", account).
		Where("deleted_at is null").
		Order("agree desc,uptime desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&comments).Error; err != nil {
		return nil, err
	}
	if len(comments) > 0 {
		return comments, nil
	}
	return nil, nil
}

func GetCommentCnt(account string) (cnt int) {
	if err := db.Table("comment").
		Where("account=?", account).Count(&cnt).Error; err != nil {
		cnt = 0
	}
	return cnt
}
