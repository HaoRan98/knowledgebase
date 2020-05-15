package models

type Notice struct {
	ID      string `json:"id" gorm:"primary_key"`
	Account string `json:"account"`
	Uptime  string `json:"uptime"`
	Msg     string `json:"msg"`
	TopicID string `json:"topic_id"`
}

func AddNotice(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}
func GetNotices(account string, pageNo, pageSize int) ([]*Notice, error) {
	var notices []*Notice
	if err := db.
		Where("account=?", account).Limit(pageSize).Offset(pageSize * (pageNo - 1)).
		Find(&notices).Error; err != nil {
		return nil, err
	}
	if len(notices) > 0 {
		return notices, nil
	}
	return nil, nil
}
func GetNoticesCnt(account string) (cnt int) {
	if err := db.Table("notice").
		Where("account=?", account).Count(&cnt).Error; err != nil {
		cnt = 0
	}
	return cnt
}
func DelNotice(id string) error {
	if err := db.Table("notice").
		Where("id=?", id).Delete(Notice{}).Error; err != nil {
		return err
	}
	return nil
}
