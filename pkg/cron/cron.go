package cron

import (
	"github.com/robfig/cron"
	"knowledgebase/models"
	"knowledgebase/pkg/export"
	"knowledgebase/pkg/file"
	"knowledgebase/pkg/jkxm"
	"log"
	"os"
	"strings"
)

func Setup() {
	go func() {
		log.Println("crontab starting...")
		c := cron.New()
		// 每天1点清理超过1天的导出记录
		if err := c.AddFunc("0 0 1 * * *", CleanUpExportFiles); err != nil {
			log.Printf("WriteIntoFile crontab failed：%v", err)
		}
		// 每天2点同步开发区数据中台指标
		if err := c.AddFunc("0 0 2 * * *", SyncJkxm); err != nil {
			log.Printf("SyncJkxm crontab failed：%v", err)
		}
		c.Run()
	}()
}

//清理超过一天的导出文件
func CleanUpExportFiles() {
	dirpath := export.GetExcelFullPath()
	files, err := file.FindFilesOlderThanDate(dirpath, 1)
	errNotExist := "open : The system cannot find the file specified."
	if err != nil && err.Error() != errNotExist {
		log.Println("CleanUp ExportFiles err:", err)
		return
	}
	for _, fileInfo := range files {
		if strings.Contains(fileInfo.Name(), "jkxm_") {
			continue
		}
		err = os.Remove(dirpath + fileInfo.Name())
		if err != nil {
			log.Println("CleanUp ExportFiles err:", err)
		}
	}
}

//定时同步开发区数据中台指标
func SyncJkxm() {
	ldDate := models.NewLdDate("")
	jkxm.SyncJkxmQs(ldDate)
	jkxm.SyncJkxmCktsba()
}
