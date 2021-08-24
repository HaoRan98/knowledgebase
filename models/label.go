package models

import (
	"fmt"
	"log"
	"time"
)

//标签
type Label struct {
	ID        string     `json:"id" gorm:"primary_key"`
	DeletedAt *time.Time `sql:"index"`
	TopicID   string     `json:"topic_id"`
	Content   string     `json:"content" gorm:"COMMENT:'标签内容';size:65535"`
	Author    string     `json:"author"`
	Account   string     `json:"account"`
	Deptname  string     `json:"deptname"`
	Uptime    string     `json:"uptime"`
	Agree     int        `json:"agree" gorm:"COMMENT:'赞同数';default:'0'"`
	Accept    bool       `json:"accept" gorm:"COMMENT:'采纳标志';default:'0'"`
}

// 创建标签
func CreateLabel(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func Is_label_exist(content string, topic_id string) ([]*Label, bool) {
	var label []*Label
	db.Table("label").Where("content=?", content).
		Where("topic_id=?", topic_id).
		Find(&label)

	if len(label) == 0 {
		//log.Println("label是空的")
		return nil, false
	} else {
		//log.Println("label不是空的")
		return label, true
	}
}

// 获取标签
func GetLabel(id string) (*Label, error) {
	var label Label
	if err := db.Where("id=?", id).First(&label).Error; err != nil {
		return nil, err
	}
	return &label, nil
}

// 修改标签
func EditLabel(label *Label) error {
	if err := db.Model(&Label{}).Where("id=?", label.ID).Updates(label).Error; err != nil {
		return err
	}
	return nil
}

// 获取标签列表
func GetLabels(topicId string, pageNo, pageSize int) ([]*Label, error) {
	var labels []*Label
	if err := db.Table("label").
		Where("topic_id like ?", "%"+topicId+"%").
		Order("accept desc,agree desc,uptime desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&labels).Error; err != nil {
		return nil, err
	}
	if len(labels) > 0 {
		return labels, nil
	}
	return nil, nil
}

// 删除标签
func DelLabel(id string) error {
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
	if err := tx.Table("label").
		Where("id=?", id).Delete(Label{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//if err := tx.Table("label").
	//	Where("label_id=?", id).Delete(Comment{}).Error; err != nil {
	//	tx.Rollback()
	//	return err
	//}
	return tx.Commit().Error
}

func LabelAgree(id string) error {
	s := fmt.Sprintf("update label set agree=agree+1 where id='%s'", id)
	if err := db.Exec(s).Error; err != nil {
		return err
	}
	return nil
}
func RemoveLabelAgree(id string) error {
	s := fmt.Sprintf("update label set agree=agree-1 where id='%s'", id)
	if err := db.Exec(s).Error; err != nil {
		return err
	}
	return nil
}

func AgreeBool(id, account string) bool {

	var agree []*Agree
	db.Table("agree").Where("agreeid=?", id).Where("account=?", account).Find(&agree)

	if len(agree) == 0 {
		log.Println("空的", agree)
		return false
	} else {
		log.Println("不空", agree)
		return true
	}

}

func GetLabelsCnt(topicId string) (cnt int) {
	if err := db.Table("label").
		Where("topic_id like ?", "%"+topicId+"%").Count(&cnt).Error; err != nil {
		cnt = 0
	}
	return cnt
}
