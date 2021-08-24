package models

import (
	"errors"
	"fmt"
	"log"
)

type Group struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Creator    string `json:"creator"`
	Account    string `json:"account"`
	JGDM       string `json:"jgdm"`
	JGMC       string `json:"jgmc"`
	Type       string `json:"type"`
	UserCnt    int    `json:"userCnt"`
	Del        int    `json:"del"`
	Uptime     string `json:"uptime"`
	LastChated string `json:"last_chated"`
	Members    []Member
}

// 创建团队
func CreateGroup(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

// 团队是否存在
func GroupIsExit(name string) (bool, error) {
	var group []Group
	if err := db.Where("name = ?", name).Find(&group).Error; err != nil {
		return false, err
	}

	if len(group) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// 获取团队列表
func GetGroups(account string) ([]Group, error) {
	var group []Group
	err := db.Debug().Where("account = ?", account).Find(&group).Error
	if err != nil {
		log.Println(err)
		return group, err
	}

	return group, nil
}

// 删团队
func DelGroup(groupid string) error {
	if err := db.Debug().Table("group").
		Where("id = ?", groupid).
		Update("del", 1).Error; err != nil {
		return err
	}

	if err := db.Debug().Table("member").Where("groupId = ?", groupid).
		Update("status", 1).Error; err != nil {
		return err
	}

	return nil
}

// 改团队信息
func EditGroup(data *Group) error {
	if err := db.Model(&Group{}).Where("id=?", data.ID).
		Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// 团队成员+1
func AddUserCnt(id string) {
	s := fmt.Sprintf("update group set userCnt=userCnt+1 where id='%s'", id)
	db.Exec(s)
}

// 查询团队
func SeleteGroup(name string) ([]Group, error) {
	var group []Group
	err := db.Debug().Where("name like ?", "%"+name+"%").Find(&group).Error
	if err != nil {
		log.Println(err)
		return group, nil
	}

	if len(group) == 0 {
		return group, errors.New("该团队不存在")
	}

	return group, nil
}
