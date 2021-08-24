package v1

import (
	"encoding/json"
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

type TopicForm struct {
	ID       string              `json:"id"`
	Title    string              `json:"title"`
	Content  string              `json:"content"`
	FileName string              `json:"file_name"`
	FileUrl  string              `json:"file_url"`
	FB       []map[string]string `json:"fb"`
	Kind     string              `json:"kind"`
	Author   string              `json:"author"`
	Account  string              `json:"account"`
	Deptname string              `json:"deptname"`
	DeptID   string              `json:"dept_id"`
	JGDM     string              `json:"jgdm"`
	JGMC     string              `json:"jgmc"`
	ParentID string              `json:"parent_id"`
	IsSelf   string              `json:"is_self"`
	Group    []map[string]string `json:"group"`
}

type TpResp struct {
	*models.Topic_Fuben
	Replys []*RpResp
	Agreed bool   `json:"agreed"`
	Faved  bool   `json:"faved"`
	Isgly  bool   `json:"isgly"`
	Gly    string `json:"gly"`
}

type TpResp1 struct {
	*models.Topic
	Replys []*RpResp
	Agreed bool   `json:"agreed"`
	Faved  bool   `json:"faved"`
	Isgly  bool   `json:"isgly"`
	Gly    string `json:"gly"`
}

type Ranking_form struct {
	PageSize int      `json:"page_size"`
	PageNo   int      `json:"page_no"`
	TimeArr  []string `json:"time_arr"`
}

type Topic_ID struct {
	ID string `json:"id"`
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

	id := "TP-" + util.RandomString(29)

	var fb []byte
	var err error
	if form.FB == nil {
		fb = []byte("[]")
	} else {

		for _, m := range form.FB {
			file := map[string]string{
				"topic_id": id,
				"fb_name":  m["fb_name"],
				"fb_url":   m["fb_url"],
				"account":  form.Account,
				"author":   form.Author,
			}

			if err := models.UpdataFile(file); err != nil {
				log.Println("附件更新失败")
				appG.Response(http.StatusOK, e.ERROR, "附件更新失败")
				return
			}

		}

		fb, err = json.Marshal(form.FB)
		if err != nil {
			log.Println("副本序列化失败")
			appG.Response(http.StatusOK, e.ERROR, "附件序列化失败")
			return
		}
	}

	topic := &models.Topic{
		ID:          id,
		Title:       form.Title,
		Content:     form.Content,
		FileName:    form.FileName,
		FileUrl:     form.FileUrl,
		Fbs:         string(fb),
		Kind:        form.Kind,
		Author:      form.Author,
		Account:     form.Account,
		Deptname:    form.Deptname,
		Createtime:  time.Now().Format("2006-01-02 15:04:05"),
		Uptime:      time.Now().Format("2006-01-02 15:04:05"),
		DeptID:      form.DeptID,
		ParentID:    form.ParentID,
		IsSelf:      form.IsSelf,
		JGDM:        form.JGDM,
		JGMC:        form.JGMC,
		LastPublish: form.Author,
	}

	// 如果有团队
	if len(form.Group) > 0 {
		for _, m := range form.Group {
			topic.GroupName = m["groupName"]
			topic.GroupId = m["groupId"]
			if err := models.CreateTopic(topic); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR, err)
				return
			}

			BroadCastTopic()
			BroadCastCount()

			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		}
	}

	if err := models.CreateTopic(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	if form.IsSelf == "0" {
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

	var fb []byte
	var err error
	var fbslice = make([]string, 0)
	if form.FB != nil {

		topic, err := models.GetTopic(form.ID)
		if err != nil {
			log.Println("修改帖子，获取帖子失败")
			appG.Response(http.StatusInternalServerError, e.ERROR, "获取帖子失败")
		}

		for _, m := range form.FB {
			fbslice = append(fbslice, m["fb_name"])
			file := map[string]string{
				"topic_id": form.ID,
				"fb_name":  m["fb_name"],
				"fb_url":   m["fb_url"],
				"account":  topic.Account,
				"author":   form.Author,
			}

			if err := models.UpdataFile(file); err != nil {
				log.Println("附件更新失败")
				appG.Response(http.StatusOK, e.ERROR, "附件更新失败")
				return
			}

		}

		fb, err = json.Marshal(form.FB)
		if err != nil {
			log.Println("副本序列化失败")
			appG.Response(http.StatusOK, e.ERROR, "附件序列化失败")
			return
		}

		// 删除没发过来的附件
		err = models.EditTopicFile(fbslice, form.ID)
		if err != nil {
			log.Println(err)
			appG.Response(http.StatusOK, e.ERROR, "附件删除失败")
			return
		}

	} else {
		// 删附件库
		if err := models.DeleteFile(form.ID); err != nil {
			log.Println("附件删除失败")
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
	}

	account := util.GetLoginID("", c)

	var user *models.SysUser

	// 获取当前用户信息
	Clients.Range(func(k, v interface{}) bool {
		if v.(*models.SysUser).UserAccount == account {
			user = v.(*models.SysUser)
			return false
		}
		return true
	})

	topic := &models.Topic{
		ID:          form.ID,
		Title:       form.Title,
		Content:     form.Content,
		FileName:    form.FileName,
		FileUrl:     form.FileUrl,
		Kind:        form.Kind,
		Uptime:      time.Now().Format("2006-01-02 15:04:05"),
		Fbs:         string(fb),
		DeptID:      form.DeptID,
		ParentID:    form.ParentID,
		IsSelf:      form.IsSelf,
		LastPublish: user.Username,
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

	type topic_form struct {
		Id    string `json:"id"`
		IsGly bool   `json:"isgly"`
	}

	var (
		appG = app.Gin{C: c}
		form topic_form
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	topic, err := models.GetTopic(form.Id)
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

		//log.Println(util.GetDeptID("",c))

		gly, isgly, err := models.IsEditTopic(util.GetDeptID("", c), topic.DeptID, form.IsGly)
		//log.Println(topic)
		topic.Browse = topic.Browse + 1
		if err != nil {
			log.Println("当前角色是否能修改帖子身份错误 ", err)
			appG.Response(http.StatusOK, e.SUCCESS, TpResp{
				Topic_Fuben: topic, Replys: replyResps, Agreed: tpFlag, Faved: favFlag, Isgly: false, Gly: ""})
			return
		}

		appG.Response(http.StatusOK, e.SUCCESS, TpResp{
			Topic_Fuben: topic, Replys: replyResps, Agreed: tpFlag, Faved: favFlag, Isgly: isgly, Gly: gly})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func GetTopics(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		groupId  string
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
		log.Println(account)
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
	if c.Query("groupId") == "" {
		groupId = ""
	} else {
		groupId = c.Query("groupId")
	}
	topics, err := models.GetTopics(account, kind, groupId, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(topics) > 0 {
		tpResps := make([]TpResp1, 0)
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
			tpResps = append(tpResps, TpResp1{
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

type IsTop struct {
	ID    string `json:"id"`
	Istop bool   `json:"istop"`
}

func TopTopic(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
		form IsTop
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	topic := map[string]interface{}{
		"id":  form.ID,
		"top": form.Istop,
	}
	if err := models.EditTopic1(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type IsHot struct {
	ID    string `json:"id"`
	Ishot bool   `json:"ishot"`
}

func HotTopic(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form IsHot
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	topic := map[string]interface{}{
		"id":  form.ID,
		"top": form.Ishot,
	}
	if err := models.EditTopic1(topic); err != nil {
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

	account := util.GetLoginID("", c)

	var user *models.SysUser

	// 获取当前用户信息
	Clients.Range(func(k, v interface{}) bool {
		if v.(*models.SysUser).UserAccount == account {
			user = v.(*models.SysUser)
			return false
		}
		return true
	})

	topic := &models.Topic{
		ID:          id,
		Del:         true,
		Uptime:      time.Now().Format("2006-01-02 15:04:05"),
		LastPublish: user.Username,
		// 添加删除人，不是本人则为管理员，则不允许恢复
	}
	//log.Println(topic)
	if err := models.EditTopic(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	var index = "topic"
	errMsg := models.ESDeleteSingle(index, topic.ID)
	if errMsg != nil {
		log.Println("EsDeleteSingle err:", errMsg)
	} else {
		log.Println("EsDeleteSingle: ok")
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type Rec_topic struct {
	ID string `json:"id"`
}

func RecTopic(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
		form Rec_topic
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	// 获取帖子，如果lastpublish和author相同，则可以恢复
	// 如果不相同，则为管理员删除，判断部门是否相同
	topic, err := models.GetTopic(form.ID)
	if err != nil {
		log.Println("获取帖子详情失败", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	// 获取当前用户depid
	//account  := util.GetLoginID("",c)
	//if topic.LastPublish != topic.Author && topic.Account == account {
	//	appG.Response(http.StatusInternalServerError, e.ERROR, "当前用户没有权限恢复帖子")
	//	return
	//}

	account := util.GetLoginID("", c)

	var user *models.SysUser

	// 获取当前用户信息
	Clients.Range(func(k, v interface{}) bool {
		if v.(*models.SysUser).UserAccount == account {
			user = v.(*models.SysUser)
			return false
		}
		return true
	})

	a := map[string]interface{}{
		"id":           form.ID,
		"del":          false,
		"uptime":       time.Now().Format("2006-01-02 15:04:05"),
		"last_publish": user.Username,
	}

	//log.Println(a)
	if err := models.EditTopic1(a); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	//update document to Es
	var index = "topic"
	jsonStr, err := util.ToJson(topic)
	util.ShowError("EsDoc ToJson err", err)

	resp, errMsg := models.EsDocument(index, topic.ID, jsonStr)
	if errMsg != "" {
		log.Println("EsDocument err:", errMsg)
	} else {
		log.Println("EsDocument:", resp)
	}

	appG.Response(http.StatusOK, e.SUCCESS, "操作成功")
}

func DelTopic_zhen(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
		form Topic_ID
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	// 删帖
	if err := models.DeleteTopic(form.ID); err != nil {
		log.Println("帖子删除失败")
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	// 删回帖和评论
	if err := models.DeleteReply(form.ID); err != nil {
		log.Println("回帖和评论删除失败")
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	// 删附件库
	if err := models.DeleteFile(form.ID); err != nil {
		log.Println("附件删除失败")
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	// 删ES数据库
	var index = "topic"
	errMsg := models.ESDeleteSingle(index, form.ID)
	if errMsg != nil {
		log.Println("EsDeleteSingle err:", errMsg)
	} else {
		log.Println("EsDeleteSingle: ok")
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func QueryDel(c *gin.Context) {

	type querydel struct {
		Parent_id string `json:"parent_id"`
		Username  string `json:"username"`
		PageSize  int    `json:"page_size"`
		PageNo    int    `json:"page_no"`
	}

	var (
		appG = app.Gin{C: c}
		form querydel
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	if form.Parent_id != "" {
		if form.Parent_id == "0" {
			topics, err, cnt := models.QueryDelTopic_jg("13706", form.PageSize, form.PageNo)
			if err != nil {
				log.Println("下级部门删除帖子查询有误，", err)
				appG.Response(http.StatusOK, e.ERROR, "下级部门查询有误，"+error.Error(err))
				return
			}
			data := map[string]interface{}{
				"list":  topics,
				"total": cnt,
			}

			appG.Response(http.StatusOK, e.SUCCESS, data)
			return
		}

		if form.Parent_id[:7] == "1370600" {
			topics, err, cnt := models.QueryDelTopic_jg("13706", form.PageSize, form.PageNo)
			if err != nil {
				log.Println("下级部门删除帖子查询有误，", err)
				appG.Response(http.StatusOK, e.ERROR, "下级部门查询有误，"+error.Error(err))
				return
			}
			data := map[string]interface{}{
				"list":  topics,
				"total": cnt,
			}

			appG.Response(http.StatusOK, e.SUCCESS, data)
			return
		}

		topics, err, cnt := models.QueryDelTopic_jg(form.Parent_id[:7], form.PageSize, form.PageNo)
		if err != nil {
			log.Println("下级部门删除帖子查询有误，", err)
			appG.Response(http.StatusOK, e.ERROR, "下级部门查询有误，"+error.Error(err))
			return
		}
		data := map[string]interface{}{
			"list":  topics,
			"total": cnt,
		}

		appG.Response(http.StatusOK, e.SUCCESS, data)
		return

	} else {

		topics, err, cnt := models.QueryDelTopic_user(form.Username, form.PageSize, form.PageNo)
		if err != nil {
			log.Println("个人删除帖子查询有误，", err)
			appG.Response(http.StatusOK, e.ERROR, "个人帖子查询有误，"+error.Error(err))
			return
		}
		data := map[string]interface{}{
			"list":  topics,
			"total": cnt,
		}

		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}

}

func QueryTopic(c *gin.Context) {
	type querydel struct {
		Deptid   string `json:"deptid"`
		PageSize int    `json:"pagesize"`
		PageNo   int    `json:"pageno"`
	}

	var (
		appG = app.Gin{C: c}
		form querydel
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	topics, err := models.QueryTopic_jg(form.Deptid[:7], form.PageSize, form.PageNo)
	if err != nil {
		log.Println("下级部门帖子查询有误，", err)
		appG.Response(http.StatusOK, e.ERROR, "下级部门查询有误，"+error.Error(err))
		return
	}
	if len(topics) > 0 {
		tpResps := make([]TpResp1, 0)
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
			tpResps = append(tpResps, TpResp1{
				Topic: tp, Replys: replyResps, Agreed: tpFlag, Faved: favFlag})
		}
		appG.Response(http.StatusOK, e.SUCCESS,
			map[string]interface{}{
				"list":  tpResps,
				"total": models.CntTopic_jg(form.Deptid[:7]),
			})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func TopicAgree(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Query("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}

	account := util.GetLoginID("", c)
	if models.AgreeBool(id, account) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "您已经点过赞了，请不要重复点赞")
		return
	} else {

		if err := models.TopicAgree(id); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, "点赞错误"+error.Error(err))
			return
		}
		agreed := &models.Agree{
			ID:      "AG-" + util.RandomString(30),
			Agreeid: id,
			Account: util.GetLoginID("", c),
			Uptime:  time.Now().Format("2006-01-02 15:04:05"),
		}
		if err := models.Agreed(agreed); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, "存储点赞表错误"+error.Error(err))
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, "点赞成功")
		return
	}
}

func RemoveTopicAgree(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Query("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	if err := models.RemoveTopicAgree(id); err != nil {
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

// 帖子排名
func Topic_user_Rank(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Ranking_form
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	topic, err, cnt := models.TopicRankRY(form.PageSize, form.PageNo, form.TimeArr[0], form.TimeArr[1])
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, "排名查询有误")
		log.Println("人员排名 failed：", err)
		return
	}

	data := map[string]interface{}{
		"list":  topic,
		"total": cnt,
	}

	appG.Response(http.StatusOK, e.SUCCESS, data)

}

func Topic_dept_Rank(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Ranking_form
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	topic, err, cnt := models.TopicRankJG(form.PageSize, form.PageNo, form.TimeArr[0], form.TimeArr[1])
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, "排名查询有误")
		log.Println("机关排名 failed：", err)
		return
	}

	data := map[string]interface{}{
		"list":  topic,
		"total": cnt,
	}

	appG.Response(http.StatusOK, e.SUCCESS, data)

}
