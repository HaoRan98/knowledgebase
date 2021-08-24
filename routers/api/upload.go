package api

import (
	"github.com/gin-gonic/gin"
	"knowledgebase/models"
	"knowledgebase/pkg/setting"
	"knowledgebase/pkg/util"
	v1 "knowledgebase/routers/api/v1"
	"log"
	"net/http"
	"time"

	"knowledgebase/pkg/app"
	"knowledgebase/pkg/e"
	"knowledgebase/pkg/logging"
	"knowledgebase/pkg/upload"
)

func UploadFile(c *gin.Context) {
	appG := app.Gin{C: c}
	fHeader, err := c.FormFile("file")
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if fHeader == nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// 限制总大小，如果单人上传附件总大小大于10G，则不让上传

	account := util.GetLoginID("", c)

	var user *models.SysUser

	// 获取当前用户信息
	v1.Clients.Range(func(k, v interface{}) bool {
		if v.(*models.SysUser).UserAccount == account {
			user = v.(*models.SysUser)
			return false
		}
		return true
	})

	fileName := upload.GetFileName(fHeader.Filename)
	fullPath := upload.GetFileFullPath()
	src := fullPath + fileName

	if !upload.CheckFileExt(fileName) {
		appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_FILE_FORMAT, nil)
		return
	}
	if !(fHeader.Size <= (int64(setting.AppSetting.FileMaxSize) << 20)) {
		appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_FILE_SIZE, nil)
		return
	}

	err = upload.CheckFile(fullPath)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_CHECK_FILE_FAIL, nil)
		return
	}

	if err := c.SaveUploadedFile(fHeader, src); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_FILE_FAIL, nil)
		return
	}

	files := &models.File{
		ID:      "FN" + util.RandomString(29),
		TopicID: "",
		FbName:  fHeader.Filename,
		FbPath:  src,
		FbSize:  fHeader.Size,
		FbUrl:   upload.GetFileFullUrl(fileName),
		Account: account,
		Author:  user.Username,
		Uptime:  time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := models.CreateFile(files); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, "附件上传失败 : "+error.Error(err))
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"file_name": fHeader.Filename,
		"file_url":  upload.GetFileFullUrl(fileName),
	})
}

type FileName struct {
	Topicid  string `json:"topicID"`
	Fileid   string `json:"fileID"`
	Pagesize int    `json:"pageSize"`
	Pageno   int    `json:"pageNo"`
}

func DeleteFile(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
		form FileName
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	err := models.DelFile(form.Topicid, form.Fileid)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, "附件删除失败 : "+error.Error(err))
		return
	}

}

func GetFiles(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form FileName
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	files, err := models.GetFiles(form.Pagesize, form.Pageno)
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, "获取附件列表失败 : "+error.Error(err))
		return
	}

	cnt := models.FilesCnt()

	data := map[string]interface{}{
		"files": files,
		"total": cnt,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)

}
