package v1

import (
	"fmt"
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

type ReplyForm struct {
	ID       string `json:"id"`
	TopicID  string `json:"topic_id"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Account  string `json:"account"`
	Deptname string `json:"deptname"`
	JGDM     string `json:"jgdm"`
	JGMC     string `json:"jgmc"`
}

type RpResp struct {
	*models.Reply
	Comments []*CtResp
	Agreed   bool `json:"agreed"`
}

type RpResp1 struct {
	*models.Reply
	Comments []*CtResp1
	Kind     string `json:"kind"`
	Title    string `json:"title"`
}

func PostReply(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form ReplyForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	floor, err := models.GenReplyFloor(form.TopicID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	reply := &models.Reply{
		ID:         "RE-" + util.RandomString(29),
		TopicID:    form.TopicID,
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
	if err := models.CreateReply(reply); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	// 更新帖子回复时间
	topicmap := map[string]interface{}{
		"uptime":       t,
		"last_publish": form.Author,
		"id":           form.TopicID,
	}
	if err := models.EditTopic1(topicmap); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, "更新回帖时间失败")
		return
	}

	//insert document to Es
	var index = "reply"
	jsonStr, err := util.ToJson(reply)
	util.ShowError("EsDoc ToJson err", err)

	resp, errMsg := models.EsDocument(index, reply.ID, jsonStr)
	if errMsg != "" {
		log.Println("EsDocument err:", errMsg)
	} else {
		log.Println("EsDocument:", resp)
	}

	BroadCastCount()

	// 通知发帖人
	topic, err := models.GetTopic(form.TopicID)
	if err != nil {
		log.Println("GetTopic in reply err:", err)
	} else {
		msg := fmt.Sprintf("您发布的帖子:\"%s\",收到一条新回帖", topic.Title)
		notice := &models.Notice{
			ID:      "NTE-" + util.RandomString(28),
			TopicID: topic.ID,
			Account: topic.Account,
			Msg:     msg,
			Uptime:  t,
		}
		BroadCastReply(notice)
	}

	// 通知收藏人
	favorites, err := models.GetCollector(form.TopicID)
	if err != nil {
		log.Println("GetCollector in reply err:", err)
	} else {
		if len(favorites) > 0 {
			msg := fmt.Sprintf("您收藏的帖子:\"%s\" ,收到一条新回帖", topic.Title)
			for _, favorite := range favorites {
				notice := &models.Notice{
					ID:      "NTE-" + util.RandomString(28),
					TopicID: form.TopicID,
					Account: favorite.Account,
					Msg:     msg,
					Uptime:  t,
				}
				BroadCastReply(notice)
			}
		}
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func EditReply(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form ReplyForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	reply := &models.Reply{
		ID:      form.ID,
		Content: form.Content,
		Uptime:  t,
	}
	if err := models.EditReply(reply); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	//update document to Es
	replyEs, err := models.GetReply(form.ID)
	if err != nil {
		log.Println("Get Reply By Key err:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	var index = "reply"
	jsonStr, err := util.ToJson(replyEs)
	util.ShowError("EsDoc ToJson err", err)

	resp, errMsg := models.EsDocument(index, replyEs.ID, jsonStr)
	if errMsg != "" {
		log.Println("EsDocument err:", errMsg)
	} else {
		log.Println("EsDocument:", resp)
	}

	BroadCastCount()

	// 通知发帖人
	topic, err := models.GetTopic(replyEs.TopicID)
	if err != nil {
		log.Println("GetTopic in reply err:", err)
	} else {
		msg := fmt.Sprintf("您发布的帖子:\"%s\" ,收到一条新回帖", topic.Title)
		notice := &models.Notice{
			ID:      "NTE-" + util.RandomString(28),
			TopicID: topic.ID,
			Account: topic.Account,
			Msg:     msg,
			Uptime:  t,
		}
		BroadCastReply(notice)
	}

	// 通知收藏人
	favorites, err := models.GetCollector(form.TopicID)
	if err != nil {
		log.Println("GetCollector in reply err:", err)
	} else {
		if len(favorites) > 0 {
			msg := fmt.Sprintf("您收藏的帖子:\"%s\" ,收到一条新回帖", topic.Title)
			for _, favorite := range favorites {
				notice := &models.Notice{
					ID:      "NTE-" + util.RandomString(28),
					TopicID: form.TopicID,
					Account: favorite.Account,
					Msg:     msg,
					Uptime:  t,
				}
				BroadCastReply(notice)
			}
		}
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func GetReplies(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		err      error
		topicId  string
		pageSize int
		pageNo   int
	)
	if c.Query("id") == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	} else {
		topicId = c.Query("id")
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

	//浏览量+1
	//models.AddBrowse(topicId)
	//BroadCastCount()

	replies, err := models.GetReplies(topicId, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(replies) > 0 {
		rpResps := make([]*RpResp, 0)
		for _, rp := range replies {
			loginId := util.GetLoginID("", c)
			ctResps := make([]*CtResp, 0)
			for _, ct := range rp.Comments {
				ctFlag := models.IsAgreed(ct.ID, loginId)
				comment := ct
				ctResps = append(ctResps, &CtResp{&comment, ctFlag})
			}
			rpFlag := models.IsAgreed(rp.ID, loginId)
			rpResps = append(rpResps, &RpResp{rp, ctResps, rpFlag})
		}
		appG.Response(http.StatusOK, e.SUCCESS,
			map[string]interface{}{
				"list": rpResps,
				"cnt":  models.GetRepliesCnt(topicId),
			})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func ReplyAgree(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	if err := models.ReplyAgree(id); err != nil {
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

func RemoveReplyAgree(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	if err := models.RemoveReplyAgree(id); err != nil {
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

func AcceptReply(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	reply := &models.Reply{
		ID:     id,
		Accept: true,
	}
	if err := models.EditReply(reply); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 删除评论
func DelReply(c *gin.Context) {
	var appG = app.Gin{C: c}
	if err := models.DelReply(c.Param("id")); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	var index = "reply"
	errMsg := models.ESDeleteSingle(index, c.Param("id"))
	if errMsg != nil {
		log.Println("EsDeleteSingle err:", errMsg)
	} else {
		log.Println("EsDeleteSingle: ok")
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 我的回帖和评论
func GetRepCom(c *gin.Context) {

	type user struct {
		Account  string `json:"account"`
		PageSize int    `json:"pagesize"`
		PageNo   int    `json:"pageno"`
	}

	var (
		appG = app.Gin{C: c}
		form user
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	var data = make([]interface{}, 0)

	replies, err := models.ComGetReplies(form.Account, form.PageNo, form.PageSize)
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	if len(replies) > 0 {
		for _, rp := range replies {
			ctResps := make([]*CtResp1, 0)
			for _, ct := range rp.Comments {
				comment := ct
				ctResps = append(ctResps, &CtResp1{&comment, "0", ""})
			}
			title := models.GetTitle(rp.TopicID)
			//rpResps = append(rpResps, &RpResp{rp, ctResps, rpFlag})
			data = append(data, &RpResp1{rp, ctResps, "1", title.Title})
		}
	} else {
		log.Println("没查到回帖")
	}

	repliesCnt := models.GetAccountRepliesCnt(form.Account)

	comments, err := models.RepGetComments(form.Account, form.PageNo, form.PageSize)
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	if len(comments) > 0 {
		for _, ct := range comments {
			title := models.GetTitle(ct.TopicID)
			data = append(data, &CtResp1{ct, "0", title.Title})
		}
	} else {
		log.Println("没查到评论")
	}

	commentCnt := models.GetCommentCnt(form.Account)

	appG.Response(http.StatusOK, e.SUCCESS,
		map[string]interface{}{
			"list":  data,
			"total": commentCnt + repliesCnt,
		})

}
