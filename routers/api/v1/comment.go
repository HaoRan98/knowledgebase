package v1

import (
	"github.com/gin-gonic/gin"
	"knowledgebase/models"
	"knowledgebase/pkg/app"
	"knowledgebase/pkg/e"
	"knowledgebase/pkg/util"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CommentForm struct {
	ID       string `json:"id"`
	ReplyID  string `json:"reply_id"`
	Floor    int    `json:"floor"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Account  string `json:"account"`
	Deptname string `json:"deptname"`
	JGDM     string `json:"jgdm"`
	JGMC     string `json:"jgmc"`
}

type CtResp struct {
	*models.Comment
	Agreed bool `json:"agreed"`
}

type CtResp1 struct {
	*models.Comment
	Kind  string `json:"kind"`
	Title string `json:"title"`
}

func PostComment(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form CommentForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	floor, err := models.GenCommentFloor(form.ReplyID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	reply, err := models.GetReply(form.ReplyID)
	if err != nil {
		log.Println("GetTopic in comment err:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	// 调用.5接口 获取人员信息

	t := time.Now().Format("2006-01-02 15:04:05")
	comment := &models.Comment{
		ID:         "CMT-" + util.RandomString(28),
		TopicID:    reply.TopicID,
		ReplyID:    form.ReplyID,
		Content:    form.Content,
		Author:     form.Author,
		Account:    form.Account,
		Deptname:   form.Deptname,
		Floor:      floor,
		Createtime: t,
		Uptime:     t,
		JGDM:       form.JGDM,
		JGMC:       form.JGMC,
	}
	if err := models.CreateComment(comment); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	// 更新帖子回复时间
	topicmap := map[string]interface{}{
		"uptime":       t,
		"last_publish": form.Author,
		"id":           reply.TopicID,
	}
	if err := models.EditTopic1(topicmap); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, "更新回帖时间失败")
		return
	}

	//insert document to Es
	var index = "comment"
	jsonStr, err := util.ToJson(comment)
	util.ShowError("EsDoc ToJson err", err)

	resp, errMsg := models.EsDocument(index, comment.ID, jsonStr)
	if errMsg != "" {
		log.Println("EsDocument err:", errMsg)
	} else {
		log.Println("EsDocument:", resp)
	}

	BroadCastCount()

	// 通知回帖人
	notice := &models.Notice{
		ID:      "NTE-" + util.RandomString(28),
		TopicID: reply.TopicID,
		Account: reply.Account,
		Msg:     "您的回帖收到一条新评论",
		Uptime:  t,
	}
	BroadCastReply(notice)

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func EditComment(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form CommentForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	comment := &models.Comment{
		ID:      form.ID,
		Content: form.Content,
		Uptime:  t,
	}
	if err := models.EditComment(comment); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	//update document to Es
	commentEs, err := models.GetCommentByID(form.ID)
	if err != nil {
		log.Println("Get Comment By Key err:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	var index = "comment"
	jsonStr, err := util.ToJson(commentEs)
	util.ShowError("EsDoc ToJson err", err)

	resp, errMsg := models.EsDocument(index, commentEs.ID, jsonStr)
	if errMsg != "" {
		log.Println("EsDocument err:", errMsg)
	} else {
		log.Println("EsDocument:", resp)
	}

	BroadCastCount()

	// 通知回帖人
	reply, err := models.GetReply(commentEs.ReplyID)
	if err != nil {
		log.Println("GetTopic in comment err:", err)
	} else {
		notice := &models.Notice{
			ID:      "NTE-" + util.RandomString(28),
			TopicID: reply.TopicID,
			Account: reply.Account,
			Msg:     "您的回帖收到一条新评论",
			Uptime:  t,
		}
		BroadCastReply(notice)
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func GetComments(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		err      error
		replyId  string
		pageSize int
		pageNo   int
	)
	if c.Query("id") == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	} else {
		replyId = c.Query("id")
	}
	if c.Query("pageNo") == "" {
		pageNo = 1
	} else {
		pageNo, _ = strconv.Atoi(c.Query("pageNo"))
	}
	if c.Query("pageSize") == "" {
		pageSize = 10000
	} else {
		pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	}
	comments, err := models.GetComments(replyId, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(comments) > 0 {
		ctResps := make([]*CtResp, 0)
		for _, ct := range comments {
			loginId := util.GetLoginID("", c)
			flag := models.IsAgreed(ct.ID, loginId)
			ctResps = append(ctResps, &CtResp{ct, flag})
		}
		appG.Response(http.StatusOK, e.SUCCESS, ctResps)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func CommentAgree(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	if err := models.CommentAgree(id); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	agreed := &models.Agree{
		ID:      "AG-" + util.RandomString(30),
		Agreeid: id,
		Account: util.GetLoginID("", c),
		Uptime:  time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.Agreed(agreed); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func RemoveCommentAgree(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	if err := models.RemoveCommentAgree(id); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	account := util.GetLoginID("", c)
	if err := models.RemoveAgreed(id, account); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func DelComment(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not nil")
		return
	}
	if err := models.DelComment(id); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	var index = "comment"
	errMsg := models.ESDeleteSingle(index, id)
	if errMsg != nil {
		log.Println("EsDeleteSingle err:", errMsg)
	} else {
		log.Println("EsDeleteSingle: ok")
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
