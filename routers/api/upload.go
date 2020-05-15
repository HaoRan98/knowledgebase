package api

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"NULL/knowledgebase/pkg/app"
	"NULL/knowledgebase/pkg/e"
	"NULL/knowledgebase/pkg/logging"
	"NULL/knowledgebase/pkg/upload"
)

func UploadFile(c *gin.Context) {
	appG := app.Gin{C: c}
	file, fHeader, err := c.Request.FormFile("file")
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if fHeader == nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	fileName := upload.GetFileName(fHeader.Filename)
	fullPath := upload.GetFileFullPath()
	src := fullPath + fileName

	if !upload.CheckFileExt(fileName) {
		appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_FILE_FORMAT, nil)
		return
	}
	if !upload.CheckFileSize(file) {
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

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"file_name": fHeader.Filename,
		"file_url":  upload.GetFileFullUrl(fileName),
	})
}
