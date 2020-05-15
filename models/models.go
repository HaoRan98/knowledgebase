package models

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"NULL/knowledgebase/pkg/setting"
	"log"
)

var (
	db   *gorm.DB
	es   *elasticsearch.Client
	Info *InfoCount
)

func Setup() {
	var err error
	// Initialize database
	db, err = gorm.Open(setting.DatabaseSetting.Type,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			setting.DatabaseSetting.User,
			setting.DatabaseSetting.Password,
			setting.DatabaseSetting.Host,
			setting.DatabaseSetting.Name))

	//Initialize Es client
	es, err = NewEsClient()

	if err != nil {
		log.Fatalf("models setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	CheckTable()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CheckTable() {
	if !db.HasTable("agree") {
		db.CreateTable(Agree{})
	} else {
		db.AutoMigrate(Agree{})
	}
	if !db.HasTable("topic") {
		db.CreateTable(Topic{})
	} else {
		db.AutoMigrate(Topic{})
	}
	if !db.HasTable("reply") {
		db.CreateTable(Reply{})
	} else {
		db.AutoMigrate(Reply{})
	}
	if !db.HasTable("comment") {
		db.CreateTable(Comment{})
	} else {
		db.AutoMigrate(Comment{})
	}
	if !db.HasTable("notice") {
		db.CreateTable(Notice{})
	} else {
		db.AutoMigrate(Notice{})
	}
	if !db.HasTable("favorite") {
		db.CreateTable(Favorite{})
	} else {
		db.AutoMigrate(Favorite{})
	}
	if !db.HasTable("info_count") {
		db.CreateTable(InfoCount{})
	} else {
		db.AutoMigrate(InfoCount{})
	}
	if !db.HasTable("kind") {
		db.CreateTable(Kind{})
	} else {
		db.AutoMigrate(Kind{})
	}
}

func InitDb() {
	CheckTable()
	var cnt int
	err := db.Select("id").Model(&InfoCount{}).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("init InfoCount error: %v", err)
		return
	}
	if cnt == 0 {
		err = db.Table("info_count").Save(&InfoCount{ID: "infoCount", Browse: 0}).Error
		if err != nil {
			log.Printf("init save InfoCount error: %v", err)
			return
		}
	}
}
