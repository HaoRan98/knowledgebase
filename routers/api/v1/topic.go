package v1

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/app"
	"NULL/knowledgebase/pkg/e"
	"NULL/knowledgebase/pkg/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TopicForm struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	FileName string `json:"file_name"`
	FileUrl  string `json:"file_url"`
	Kind     string `json:"kind"`
	Author   string `json:"author"`
	Account  string `json:"account"`
	Deptname string `json:"deptname"`
}

type TpResp struct {
	*models.Topic
	Replys []*RpResp
	Agreed bool `json:"agreed"`
	Faved  bool `json:"faved"`
}

func ImpTopic(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form TopicForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	topic := &models.Topic{
		ID:       "TP-" + util.RandomString(29),
		Title:    form.Title,
		Content:  form.Content,
		Author:   "第三方接入",
		Account:  "dsfjr",
		Deptname: "智税实验室",
		Uptime:   time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.CreateTopic(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	//insert document to Es
	var index = "topic"
	jsonStr, err := util.ToJson(topic)
	util.ShowError("EsDoc ToJson err", err)

	resp, errMsg := models.EsDocument(index, topic.ID, jsonStr)
	if errMsg != "" {
		log.Println("EsDocument err:", errMsg)
	} else {
		log.Println("EsDocument:", resp)
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func PostTopic(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form TopicForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	if form.Kind != "" {
		if !models.IsKindExist(form.Kind) {
			err := models.AddKind(&models.Kind{Mc: form.Kind})
			if err != nil {
				log.Println("add kind err:", err)
			}
		}
	}

	topic := &models.Topic{
		ID:       "TP-" + util.RandomString(29),
		Title:    form.Title,
		Content:  form.Content,
		FileName: form.FileName,
		FileUrl:  form.FileUrl,
		Kind:     form.Kind,
		Author:   form.Author,
		Account:  form.Account,
		Deptname: form.Deptname,
		Uptime:   time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.CreateTopic(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	//insert document to Es
	var index = "topic"
	jsonStr, err := util.ToJson(topic)
	util.ShowError("EsDoc ToJson err", err)

	resp, errMsg := models.EsDocument(index, topic.ID, jsonStr)
	if errMsg != "" {
		log.Println("EsDocument err:", errMsg)
	} else {
		log.Println("EsDocument:", resp)
	}

	BroadCastTopic()
	BroadCastCount()

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func EditTopic(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form TopicForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	if form.Kind != "" {
		if !models.IsKindExist(form.Kind) {
			err := models.AddKind(&models.Kind{Mc: form.Kind})
			if err != nil {
				log.Println("add kind err:", err)
			}
		}
	}

	topic := &models.Topic{
		ID:       form.ID,
		Title:    form.Title,
		Content:  form.Content,
		FileName: form.FileName,
		FileUrl:  form.FileUrl,
		Kind:     form.Kind,
		Uptime:   time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.EditTopic(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	//update document to Es
	topicEs, err := models.GetTopic(form.ID)
	if err != nil {
		log.Println("Get Topic By Key err:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	var index = "topic"
	jsonStr, err := util.ToJson(topicEs)
	util.ShowError("EsDoc ToJson err", err)

	resp, errMsg := models.EsDocument(index, topicEs.ID, jsonStr)
	if errMsg != "" {
		log.Println("EsDocument err:", errMsg)
	} else {
		log.Println("EsDocument:", resp)
	}

	BroadCastTopic()
	BroadCastCount()

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func GetTopic(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	topic, err := models.GetTopic(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(topic.ID) > 0 {
		//浏览量+1
		models.AddBrowse(topic.ID)
		BroadCastCount()

		loginId := util.GetLoginID("", c)
		replyResps := make([]*RpResp, 0)
		for _, rp := range topic.Replys {
			rpFlag := models.IsAgreed(rp.ID, loginId)
			reply := rp
			replyResps = append(replyResps, &RpResp{&reply, nil, rpFlag})
		}
		tpFlag := models.IsAgreed(topic.ID, loginId)
		favFlag := models.IsFavorite(topic.ID, loginId)
		appG.Response(http.StatusOK, e.SUCCESS, TpResp{
			Topic: topic, Replys: replyResps, Agreed: tpFlag, Faved: favFlag})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func GetTopics(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		err      error
		account  string
		kind     string
		pageSize int
		pageNo   int
	)
	switch c.Query("flag") {
	case "list":
		account = ""
	case "pub":
		account = util.GetLoginID("", c)
	default:
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, nil)
		return
	}
	if c.Query("kind") == "" {
		kind = ""
	} else {
		kind = c.Query("kind")
	}
	if c.Query("pageNo") == "" {
		pageNo = 1
	} else {
		pageNo, _ = strconv.Atoi(c.Query("pageNo"))
	}
	if c.Query("pageSize") == "" {
		pageSize = 100
	} else {
		pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	}
	topics, err := models.GetTopics(account, kind, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(topics) > 0 {
		tpResps := make([]TpResp, 0)
		for _, tp := range topics {
			loginId := util.GetLoginID("", c)
			replyResps := make([]*RpResp, 0)
			for _, rp := range tp.Replys {
				rpFlag := models.IsAgreed(rp.ID, loginId)
				reply := rp
				replyResps = append(replyResps, &RpResp{&reply, nil, rpFlag})
			}
			tpFlag := models.IsAgreed(tp.ID, loginId)
			favFlag := models.IsFavorite(tp.ID, loginId)
			tpResps = append(tpResps, TpResp{
				Topic: tp, Replys: replyResps, Agreed: tpFlag, Faved: favFlag})
		}
		appG.Response(http.StatusOK, e.SUCCESS,
			map[string]interface{}{
				"list": tpResps,
				"cnt":  models.GetTopicsCnt(account, kind),
			})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func TopTopic(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	topic := &models.Topic{
		ID:  id,
		Top: true,
	}
	if err := models.EditTopic(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func HotTopic(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	topic := &models.Topic{
		ID:  id,
		Hot: true,
	}
	if err := models.EditTopic(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func DelTopic(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	topic := &models.Topic{
		ID:     id,
		Del:    true,
		Uptime: time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.EditTopic(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
