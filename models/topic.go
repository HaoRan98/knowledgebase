package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

//发帖
type Topic struct {
	ID          string `json:"id" gorm:"primary_key"`
	GroupId     string `json:"groupId"`
	GroupName   string `json:"groupName"`
	Del         bool   `json:"del" gorm:"COMMENT:'删除标志'"`
	Title       string `json:"title" gorm:"COMMENT:'标题'"`
	Content     string `json:"content" gorm:"COMMENT:'内容';size:65535"`
	FileName    string `json:"file_name" gorm:"COMMENT:'文件原始名'"`
	FileUrl     string `json:"file_url" gorm:"COMMENT:'文件真实路径'"`
	Fbs         string `json:"fbs"`
	Kind        string `json:"kind" gorm:"COMMENT:'类别'"`
	Author      string `json:"author"`
	Account     string `json:"account"`
	DeptID      string `json:"dept_id"`
	Deptname    string `json:"deptname"`
	ParentID    string `json:"parent_id"`
	JGDM        string `json:"jgdm"`
	JGMC        string `json:"jgmc"`
	Agree       int    `json:"agree" gorm:"COMMENT:'赞同数';default:'0'"`
	Createtime  string `json:"createtime"`
	Uptime      string `json:"uptime"`
	Browse      int    `json:"browse" gorm:"COMMENT:'浏览量';default:'0'"`
	Hot         bool   `json:"hot" gorm:"COMMENT:'热门';default:'0'"`
	Top         bool   `json:"top" gorm:"COMMENT:'置顶标志';default:'0'"`
	IsSelf      string `json:"is_self"`
	LastPublish string `json:"last_publish"`
	Replys      []Reply
}

type Files struct {
	FbName string `json:"fb_name"`
	FbURL  string `json:"fb_url"`
}

type Topic_Fuben struct {
	Topic
	FB []*File `json:"fb"`
}

func CreateTopic(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTopic(id string) error {
	if err := db.Where("id = ?", id).Delete(&Topic{}).Error; err == nil {
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

func EditTopic1(topic map[string]interface{}) error {
	if err := db.Model(&Topic{}).
		Where("id=?", topic["id"]).Updates(topic).Error; err != nil {
		return err
	}
	return nil
}

func AddBrowse(id string) {
	s := fmt.Sprintf("update topic set browse=browse+1 where id='%s'", id)
	db.Exec(s)
}

func GetTopic(id string) (*Topic_Fuben, error) {
	var topic Topic
	if err := db.Preload("Replys", func(db *gorm.DB) *gorm.DB {
		return db.Order("reply.accept,reply.floor")
	}).
		Where("id=?", id).First(&topic).Error; err != nil {
		return nil, err
	}
	if len(topic.ID) > 0 {

		fb, err := GetFile(id)
		if err != nil {
			log.Println("附件查询有误")
			return nil, err
		}

		var top_fu = &Topic_Fuben{
			Topic: topic,
			FB:    []*File{},
		}
		top_fu.FB = fb

		return top_fu, nil
	}
	return nil, nil
}

func GetTopics(account, kind, groupId string, pageNo, pageSize int) ([]*Topic, error) {
	var topics []*Topic

	tx := db.Preload("Replys", func(db *gorm.DB) *gorm.DB {
		return db.Order("reply.accept,reply.floor")
	}).
		Where("account like ?", "%"+account+"%").
		Where("kind like ?", "%"+kind+"%").
		Order("top desc,hot desc,uptime desc,agree desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1))

	if groupId != "" {
		tx.Where("del = 0").Where("groupId = ?", groupId)
	}

	if account == "" {
		if err := tx.Debug().Where("is_self = 0").Find(&topics).Error; err != nil {
			return nil, err
		}
	} else {
		if err := tx.Debug().Find(&topics).Error; err != nil {
			return nil, err
		}
	}

	if len(topics) > 0 {
		return topics, nil
	}
	return nil, nil
}

func GetTopicsCnt(account, kind string) (cnt int) {
	if err := db.Table("topic").
		Where("del = 0").
		Where("account like ?", "%"+account+"%").
		Where("kind like ?", "%"+kind+"%").
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

func TopicAgree(id string) error {
	s := fmt.Sprintf("update topic set agree=agree+1 where id='%s'", id)
	if err := db.Exec(s).Error; err != nil {
		return err
	}
	return nil
}
func RemoveTopicAgree(id string) error {
	s := fmt.Sprintf("update topic set agree=agree-1 where id='%s'", id)
	if err := db.Exec(s).Error; err != nil {
		return err
	}
	return nil
}

func GetTitle(id string) Topic {
	var topic Topic
	if err := db.Model(Topic{}).Where("id = ?", id).First(&topic).Error; err != nil {
		log.Println(err)
		return topic
	}
	return topic
}

type Rank struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Topic_num int    `json:"topic_num"`
	Browses   int    `json:"browses"`
	Agrees    int    `json:"agrees"`
	Replys    int    `json:"replys"`
}

func TopicRankRY(pageSize, pageNo int, start, end string) ([]*Rank, error, int) {

	var rank1 []*Rank
	var rank []*Rank
	var cnt int
	var tx1 = db
	var tx2 = db

	var sql = `SELECT
	topic.account as id,
	ANY_VALUE(topic.author) AS name,
	count( topic.account ) AS topic_num,
	sum( topic.browse ) AS browses,
	sum( topic.agree ) AS agrees
FROM
	topic 
where STR_TO_DATE(uptime,'%Y-%m-%d') >= STR_TO_DATE('` + start + `','%Y-%m-%d')
And STR_TO_DATE(uptime,'%Y-%m-%d') <= STR_TO_DATE('` + end + `','%Y-%m-%d')
And del = 0
GROUP BY
	topic.account 
ORDER BY
	topic_num DESC`

	if start != "" && end != "" {
		if err := db.Raw(sql).Limit(pageSize).Offset(pageSize * (pageNo - 1)).Scan(&rank1).Error; err == nil {

			var rank2 []*Rank
			if err := db.Raw(`
			SELECT
				reply.account as id,
				ANY_VALUE(reply.author) as name,
				count(reply.account) as replys
			From
				reply
			where STR_TO_DATE(uptime,'%Y-%m-%d') >= STR_TO_DATE('` + start + `','%Y-%m-%d')
			And STR_TO_DATE(uptime,'%Y-%m-%d') <= STR_TO_DATE('` + end + `','%Y-%m-%d')
			And deleted_at is null
			GROUP BY
				reply.account
		`).Scan(&rank2).Error; err == nil {

				set_rank1 := make(map[string]struct{})
				for _, rank := range rank1 {
					set_rank1[rank.ID] = struct{}{}
				}

				for _, rank := range rank2 {
					if _, ok := set_rank1[rank.ID]; ok {
						for _, r := range rank1 {
							if r.ID == rank.ID {
								r.Replys = rank.Replys
								break
							}
						}
					} else {
						rank1 = append(rank1, rank)
					}
				}
			}

			if err := tx1.Raw(sql).Scan(&rank).Error; err == nil {
				cnt = len(rank)
			}

			return rank1, nil, cnt
		} else {
			return nil, err, 0
		}
	} else {

		var sql1 = `SELECT
	topic.account as id,
	ANY_VALUE(topic.author) AS name,
	count( topic.account ) AS topic_num,
	sum( topic.browse ) AS browses,
	sum( topic.agree ) AS agrees
FROM
	topic 
where del = 0
GROUP BY
	topic.account 
ORDER BY
	topic_num DESC`

		if err := db.Raw(sql1).Limit(pageSize).Offset(pageSize * (pageNo - 1)).Scan(&rank1).Error; err == nil {

			var rank2 []*Rank
			if err := db.Raw(`
			SELECT
				reply.account as id,
				ANY_VALUE(reply.author) as name,
				count(reply.account) as replys
			From
				reply
			where deleted_at is null
			GROUP BY
				reply.account
		`).Scan(&rank2).Error; err == nil {

				set_rank1 := make(map[string]struct{})
				for _, rank := range rank1 {
					set_rank1[rank.ID] = struct{}{}
				}

				for _, rank := range rank2 {
					if _, ok := set_rank1[rank.ID]; ok {
						for _, r := range rank1 {
							if r.ID == rank.ID {
								r.Replys = rank.Replys
								break
							}
						}
					} else {
						rank1 = append(rank1, rank)
					}
				}
			}

			if err := tx2.Raw(sql1).Scan(&rank).Error; err == nil {
				cnt = len(rank)
			}

			return rank1, nil, cnt
		} else {
			return nil, err, 0
		}
	}
}

func TopicRankJG(pageSize, pageNo int, start, end string) ([]*Rank, error, int) {

	// 查询帖子

	var rank []*Rank
	var rank1 []*Rank
	var cnt int
	var tx = db
	var sql string
	if start != "" && end != "" {
		sql = `SELECT
	topic.jgdm as id,
	ANY_VALUE(topic.jgmc) AS name,
	count( topic.jgdm ) AS topic_num,
	sum( topic.browse ) AS browses,
	sum( topic.agree ) AS agrees
FROM
	topic 
where STR_TO_DATE(uptime,'%Y-%m-%d') >= STR_TO_DATE('` + start + `','%Y-%m-%d')
And STR_TO_DATE(uptime,'%Y-%m-%d') <= STR_TO_DATE('` + end + `','%Y-%m-%d')
And del = 0
GROUP BY
	topic.jgdm 
ORDER BY
	topic_num DESC`
		if err := db.Raw(sql).Limit(pageSize).Offset(pageSize * (pageNo - 1)).Scan(&rank1).Error; err == nil {
			// 查询回帖
			var rank2 []*Rank
			if err := db.Raw(`
			SELECT
				reply.jgdm as id,
				ANY_VALUE(reply.jgmc) AS name,
				count(reply.jgdm) as replys
			From
				reply
			where STR_TO_DATE(uptime,'%Y-%m-%d') >= STR_TO_DATE('` + start + `','%Y-%m-%d')
			And STR_TO_DATE(uptime,'%Y-%m-%d') <= STR_TO_DATE('` + end + `','%Y-%m-%d')
			And deleted_at is null
			GROUP BY
				reply.jgdm
		`).Scan(&rank2).Error; err == nil {

				set_rank1 := make(map[string]struct{})
				for _, rank := range rank1 {
					set_rank1[rank.ID] = struct{}{}
				}

				for _, rank := range rank2 {
					if _, ok := set_rank1[rank.ID]; ok {
						for _, r := range rank1 {
							if r.ID == rank.ID {
								r.Replys = rank.Replys
								break
							}
						}
					} else {
						rank1 = append(rank1, rank)
					}
				}

				if err := tx.Raw(sql).Scan(&rank).Error; err == nil {
					cnt = len(rank)
				}

				return rank1, nil, cnt
			} else {
				return nil, err, 0
			}

		} else {
			return nil, err, 0
		}
	} else {

		sql = `SELECT
	topic.jgdm as id,
	ANY_VALUE(topic.jgmc) AS name,
	count( topic.jgdm ) AS topic_num,
	sum( topic.browse ) AS browses,
	sum( topic.agree ) AS agrees
FROM
	topic 
where del = 0
GROUP BY
	topic.jgdm 
ORDER BY
	topic_num DESC`

		if err := db.Raw(sql).Limit(pageSize).Offset(pageSize * (pageNo - 1)).Scan(&rank1).Error; err == nil {
			// 查询回帖
			var rank2 []*Rank
			if err := db.Raw(`
			SELECT
				reply.jgdm as id,
				ANY_VALUE(reply.jgmc) AS name,
				count(reply.jgdm) as replys
			From
				reply
			where deleted_at is null
			GROUP BY
				reply.jgdm
		`).Scan(&rank2).Error; err == nil {

				set_rank1 := make(map[string]struct{})
				for _, rank := range rank1 {
					set_rank1[rank.ID] = struct{}{}
				}

				for _, rank := range rank2 {
					if _, ok := set_rank1[rank.ID]; ok {
						for _, r := range rank1 {
							if r.ID == rank.ID {
								r.Replys = rank.Replys
								break
							}
						}
					} else {
						rank1 = append(rank1, rank)
					}
				}

				if err := tx.Raw(sql).Scan(&rank).Error; err == nil {
					cnt = len(rank)
				}

				return rank1, nil, cnt
			} else {
				return nil, err, 0
			}

		} else {
			return nil, err, 0
		}
	}
}

// 查找删除的帖子

type Topic_del struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Kind     string `json:"kind"`
	Author   string `json:"author"`
	Deptname string `json:"deptname"`
	Uptime   string `json:"uptime"`
}

// 下级的
func QueryDelTopic_jg(parent_id string, pageSize, pageNo int) ([]*Topic_del, error, int) {

	var topics []*Topic_del

	if err := db.Table("topic").Where("del = 1").
		Where("dept_id like ?", "%"+parent_id+"%").
		Where("is_self = 0").
		Order("top desc,hot desc,uptime desc,agree desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&topics).Error; err != nil {
		return topics, err, 0
	}

	var cnt int
	if err := db.Table("topic").Where("del = 1").
		Where("dept_id like ?", "%"+parent_id+"%").
		Where("is_self = 0").
		Count(&cnt).Error; err != nil {
		cnt = 0
	}

	return topics, nil, cnt

}

// 下级全部帖子
func QueryTopic_jg(parent_id string, pageSize, pageNo int) ([]*Topic, error) {

	var topics []*Topic
	if parent_id == "1370600" {
		parent_id = "13706"
	}

	tx := db.Preload("Replys", func(db *gorm.DB) *gorm.DB {
		return db.Order("reply.accept,reply.floor")
	}).Table("topic").
		Where("jgdm like ?", "%"+parent_id+"%").
		Where("is_self = ?", 0).
		Order("top desc,hot desc,uptime desc,agree desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&topics)

	if err := tx.Debug().Error; err != nil {
		log.Println(err)
		return topics, err
	}

	return topics, nil
}

func CntTopic_jg(parent_id string) (cnt int) {

	tx := db.Preload("Replys", func(db *gorm.DB) *gorm.DB {
		return db.Order("reply.accept,reply.floor")
	}).Table("topic").
		Where("jgdm like ?", "%"+parent_id+"%").
		Where("is_self = ?", 0).Count(&cnt)

	if err := tx.Error; err != nil {
		cnt = 0
	}

	return cnt

}

// 自己的
func QueryDelTopic_user(account string, pageSize, pageNo int) ([]*Topic_del, error, int) {

	var topics []*Topic_del

	if err := db.Table("topic").Where("del = 1").
		Where("account = ? ", account).
		Order("top desc,hot desc,uptime desc,agree desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&topics).Error; err != nil {
		return topics, err, 0
	}

	var cnt int
	if err := db.Table("topic").Where("del = 1").
		Where("account = ? ", account).
		Count(&cnt).Error; err != nil {
		cnt = 0
	}

	return topics, nil, cnt
}

// 自己的全部帖子
func QueryTopic_user(account string, pageSize, pageNo int) ([]*Topic_del, error) {

	var topics []*Topic_del

	if err := db.Table("topic").
		Where("account = ? ", account).
		Order("top desc,hot desc,uptime desc,agree desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&topics).Error; err != nil {
		return nil, err
	}

	return topics, nil
}

// 查是否能修改
func IsEditTopic(isgly_deptid, dept_id string, isgly bool) (string, bool, error) {

	if isgly_deptid[:7] == dept_id[:7] && isgly {
		return "", true, nil
	} else if isgly_deptid[:7] == "1370600" && isgly {
		return "", true, nil
	} else {
		return "", false, errors.New("该用户不是管理员或不是该局管理员")
	}
}

func UpDateKind(kind1, kind2 string) error {

	if err := db.Debug().Table("topic").Where("kind = ?", kind1).
		Update("kind", kind2).Error; err != nil {
		return err
	}
	return nil

}
