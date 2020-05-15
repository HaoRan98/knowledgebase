package models

import (
	"fmt"
	"log"
)

type InfoCount struct {
	ID     string `json:"id" gorm:"primary_key"`
	Browse int    `json:"browse" gorm:"COMMENT:'搜索统计'"`
}

func SaveInfoCount(id string, browse int) error {
	if err := db.Table("info_count").
		Where("id=?", id).Update("browse", browse).Error; err != nil {
		return err
	}
	return nil
}
func GetInfoCount() (*InfoCount, error) {
	var cnt InfoCount
	if err := db.Table("info_count").First(&cnt).Error; err != nil {
		return nil, err
	}
	return &cnt, nil
}
func GetTopicWordCnt() int {
	type wordCnt struct {
		T int
		C int
	}
	var cnt wordCnt
	sql := fmt.Sprintf(`SELECT sum(length(title)) t,sum(length(content)) c from topic`)
	if err := db.Raw(sql).Scan(&cnt).Error; err != nil {
		log.Println("Get Word Cnt err:", err)
		return 0
	}
	return cnt.T + cnt.C
}
func GetWordCnt(tname string) int {
	type wordCnt struct {
		C int
	}
	var cnt wordCnt
	sql := fmt.Sprintf(`SELECT sum(length(content)) c from %s`, tname)
	if err := db.Raw(sql).Scan(&cnt).Error; err != nil {
		log.Println("Get Word Cnt err:", err)
		return 0
	}
	return cnt.C
}
