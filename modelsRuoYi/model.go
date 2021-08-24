package modelsRuoYi

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"knowledgebase/pkg/setting"
	"log"
)

var db *gorm.DB

func Setup() {
	var err error
	// Initialize database
	db, err = gorm.Open(setting.DatabaseSetting.Type,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			setting.DatabaseSetting.User,
			setting.DatabaseSetting.Password,
			setting.DatabaseSetting.Host,
			"ruoyi"))

	//Initialize Es client

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

}
