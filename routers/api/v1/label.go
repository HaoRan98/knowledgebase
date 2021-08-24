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

type LabelForm struct {
	ID       string `json:"id" gorm:"primary_key"`                    // 标签ID
	TopicID  string `json:"topic_id"`                                 // 帖子ID
	Content  string `json:"content" gorm:"COMMENT:'标签内容';size:65535"` //内容
	Author   string `json:"author"`                                   // 作者
	Account  string `json:"account"`                                  // 账号
	Deptname string `json:"deptname"`                                 // 部门
}

type LbResp struct {
	*models.Label
	Agreed bool `json:"agreed"`
}

// @Tags Base
// @Summary POST方法 发布标签
// @Accept application/json
// @Produce  application/json
// @Param data body LabelForm true "标签ID , 帖子ID , 标签内容 , 作者 , 账号 , 部门 "
// @Success 200 {string} string "{"code":200,"msg":"ok","data":"标签发布成功"}"
// @Router /api/v1/label/post [post]
func PostLabel(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form LabelForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	isLabel, exist := models.Is_label_exist(form.Content, form.TopicID)
	//log.Println("是否存在:",exist,len(isLabel))
	if exist == true {
		if models.AgreeBool(isLabel[0].ID, form.Account) {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		} else {

			if err := models.LabelAgree(isLabel[0].ID); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR, "点赞错误"+error.Error(err))
				return
			}
			agreed := &models.Agree{
				ID:      "AG-" + util.RandomString(30),
				Agreeid: isLabel[0].ID,
				Account: util.GetLoginID("", c),
				Uptime:  time.Now().Format("2006-01-02 15:04:05"),
			}
			if err := models.Agreed(agreed); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR, "存储点赞表错误"+error.Error(err))
				return
			}
			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		}
	}

	t := time.Now().Format("2006-01-02 15:04:05")
	label := &models.Label{
		ID:       "CMT-" + util.RandomString(28),
		TopicID:  form.TopicID,
		Content:  form.Content,
		Author:   form.Author,
		Account:  form.Account,
		Deptname: form.Deptname,
		Uptime:   t,
	}
	if err := models.CreateLabel(label); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "标签发布成功")

	//insert document to Es
	var index = "label"
	jsonStr, err := util.ToJson(label)
	util.ShowError("EsDoc ToJson err", err)

	resp, errMsg := models.EsDocument(index, label.ID, jsonStr)
	if errMsg != "" {
		log.Println("EsDocument err:", errMsg)
	} else {
		log.Println("EsDocument:", resp)
	}

	BroadCastCount()

	// 通知发帖人
	topic, err := models.GetTopic(form.TopicID)
	if err != nil {
		log.Println("GetTopic in label err:", err)
	} else {
		msg := fmt.Sprintf("您发布的帖子:\"%s\",收到一条新标签", topic.Title)
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
		log.Println("GetCollector in label err:", err)
	} else {
		if len(favorites) > 0 {
			msg := fmt.Sprintf("您收藏的帖子:\"%s\" ,收到一条新标签", topic.Title)
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

}

// @Tags Base
// @Summary POST方法 修改标签
// @Accept application/json
// @Produce  application/json
// @Param data body LabelForm true "标签ID , 帖子ID , 标签内容 , 作者 , 账号 , 部门 "
// @Success 200 {string} string "{"code":200,"msg":"ok","data":"标签修改成功"}"
// @Router /api/v1/label/edit [post]
func EditLabel(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form LabelForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	label := &models.Label{
		ID:      form.ID,
		Content: form.Content,
		Uptime:  t,
	}
	if err := models.EditLabel(label); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

	//update document to Es
	labelEs, err := models.GetLabel(form.ID)
	if err != nil {
		log.Println("Get Label By Key err:", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	var index = "label"
	jsonStr, err := util.ToJson(labelEs)
	util.ShowError("EsDoc ToJson err", err)

	resp, errMsg := models.EsDocument(index, labelEs.ID, jsonStr)
	if errMsg != "" {
		log.Println("EsDocument err:", errMsg)
	} else {
		log.Println("EsDocument:", resp)
	}

	BroadCastCount()

	// 通知发帖人
	topic, err := models.GetTopic(labelEs.TopicID)
	if err != nil {
		log.Println("GetTopic in label err:", err)
	} else {
		msg := fmt.Sprintf("您发布的帖子:\"%s\" ,收到一条新标签", topic.Title)
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
		log.Println("GetCollector in label err:", err)
	} else {
		if len(favorites) > 0 {
			msg := fmt.Sprintf("您收藏的帖子:\"%s\" ,收到一条新标签", topic.Title)
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

}

// @Tags Base
// @Summary GET方法 获取标签列表
// @Accept application/json
// @Produce  application/json
// @Param id query string true "帖子id"
// @Param pageSage query string true "条数"
// @Param pageNo query string true "当前页"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":"标签列表获取成功"}"
// @Router /api/v1/labels [get]
func GetLabels(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		err      error
		labelId  string
		pageSize int
		pageNo   int
	)
	if c.Query("id") == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	} else {
		labelId = c.Query("id")
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
	labels, err := models.GetLabels(labelId, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(labels) > 0 {
		lbresps := make([]*LbResp, 0)
		for _, ct := range labels {
			loginId := util.GetLoginID("", c)
			flag := models.IsAgreed(ct.ID, loginId)
			lbresps = append(lbresps, &LbResp{ct, flag})
		}
		appG.Response(http.StatusOK, e.SUCCESS, lbresps)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Tags Base
// @Summary GET方法 删除标签
// @Accept application/json
// @Produce  application/json
// @Param id query string true "标签id"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":"标签删除成功"}"
// @Router /api/v1/label/del [get]
func DelLabel(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Query("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not nil")
		return
	}
	if err := models.DelLabel(id); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	var index = "label"
	errMsg := models.ESDeleteSingle(index, id)
	if errMsg != nil {
		log.Println("EsDeleteSingle err:", errMsg)
	} else {
		log.Println("EsDeleteSingle: ok")
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Tags Base
// @Summary GET方法 标签点赞
// @Accept application/json
// @Produce  application/json
// @Param id query string true "标签id"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":"点赞成功"}"
// @Router /api/v1/label/agree [get]
func LabelAgree(c *gin.Context) {
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

		if err := models.LabelAgree(id); err != nil {
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

// @Tags Base
// @Summary GET方法 取消标签点赞
// @Accept application/json
// @Produce  application/json
// @Param id query string true "标签id"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":"标签删除成功"}"
// @Router /api/v1/label/agree_cancel [get]
func RemoveLabelAgree(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Query("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	account := util.GetLoginID("", c)

	if err := models.RemoveLabelAgree(id); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	if err := models.RemoveAgreed(id, account); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
