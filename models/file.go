package models

import (
	file2 "knowledgebase/pkg/file"
	"log"
)

type File struct {
	ID      string `json:"id"`
	TopicID string `json:"topic_id"`
	FbName  string `json:"fb_name"`
	FbPath  string `json:"fb_path"`
	FbSize  int64  `json:"fb_size"`
	FbUrl   string `json:"fb_url"`
	Account string `json:"account"`
	Author  string `json:"author"`
	Uptime  string `json:"uptime"`
}

type FileList struct {
	FbName string `json:"fb_name"`
	FbUrl  string `json:"fb_url"`
}

func CreateFile(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

// 删除帖子附件
func DeleteFile(topic_id string) error {

	file, err := GetFile(topic_id)
	if err != nil {
		log.Println("删除文件中，获取文件列表失败")
		return err
	}

	if len(file) > 0 {
		for _, f := range file {
			log.Println(f.FbPath)
			_, err := file2.RemoveFile(f.FbPath)
			if err != nil {
				log.Println("服务器删除文件失败")
				return err
			}
		}
		if err := db.Where("topic_id = ?", topic_id).Delete(&File{}).Error; err == nil {
			return err
		}
		return nil
	}

	return nil

}

func UpdataFile(file map[string]string) error {
	if err := db.Model(&File{}).
		Where("fb_url=?", file["fb_url"]).Updates(file).Error; err != nil {
		return err
	}
	return nil
}

func GetFile(topic_id string) ([]*File, error) {
	var files []*File
	var err error

	if err = db.Model(File{}).Where("topic_id = ?", topic_id).Find(&files).Error; err == nil {
		return files, nil
	}
	return files, err
}

func GetFiles(pageSize, pageNo int) ([]*File, error) {
	var files []*File
	var err error

	if err = db.Model(File{}).Limit(pageSize).Offset(pageSize * (pageNo - 1)).
		Find(&files).Error; err == nil {
		return files, nil
	}
	return files, err
}

func FilesCnt() (cnt int) {
	if err := db.Table("file").
		Count(&cnt).Error; err != nil {
		cnt = 0
	}
	return cnt
}

// 删除单个附件
func DelFile(topic_id, file_id string) error {

	var file File
	if err := db.Debug().Where("topic_id = ?", topic_id).Where("id = ?", file_id).First(&file).Error; err != nil {
		log.Println("查找附件失败：", err)
		return err
	}
	log.Println(file.FbPath)
	_, err := file2.RemoveFile(file.FbPath)
	if err != nil {
		log.Println("服务器删除文件失败")
		return err
	}

	tx := db.Debug().Where("topic_id = ?", topic_id).Where("id = ?", file_id).Delete(&File{})
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

// 修改帖子批量删除副本
func EditTopicFile(file []string, topic_id string) error {

	var files []File

	tx := db.Debug().Model(&File{}).Where("topic_id = ?", topic_id)

	if len(file) > 0 {
		for _, s := range file {
			tx = tx.Where("fb_name != ?", s)
		}
	}

	// 查找并删除
	sel := tx
	if err := sel.Find(&files).Error; err != nil {
		return err
	}

	if len(files) > 0 {
		// 删除本地文件
		for _, f := range files {
			log.Println(f.FbPath)
			_, err := file2.RemoveFile(f.FbPath)
			if err != nil {
				log.Println("服务器删除文件失败")
				return err
			}
		}
		if err := tx.Delete(&File{}).Error; err != nil {
			return err
		}
	}
	return nil
}
