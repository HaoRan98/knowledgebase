package v1

import (
	"github.com/gin-gonic/gin"
	"knowledgebase/models"
	"knowledgebase/pkg/app"
	"knowledgebase/pkg/e"
	"knowledgebase/pkg/util"
	"log"
	"net/http"
)

type Member struct {
	GroupID   string `json:"groupId"`
	GroupName string `json:"groupName"`
	UserName  string `json:"userName"`
	Account   string `json:"account"`
	DeptID    string `json:"deptid"`
	deptname  string `json:"deptname"`
	JGMC      string `json:"jgmc"`
	JGDM      string `json:"jgdm"`
}

// @Tags 添加成员
// @Summary POST方法 添加成员
// @Accept application/json
// @Produce  application/json
// @Param data body Member true "ID(不用传) , 团队名称 , 人 , 创账号 , 部门id , 部门名称 , 机关名称 , 机关代码"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":nil}
// @Router /api/v1/member/add [post]
func AddMember(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Member
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	ok, err := models.MemberIsExit(form.GroupID, form.Account)
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if ok {
		log.Println("该成员已加入")
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	jgdm, jgmc := util.SWJG(form.DeptID[:7])

	member := &models.Member{
		ID:        "MB-" + util.RandomString(29),
		GroupID:   form.GroupID,
		GroupName: form.GroupName,
		UserName:  form.UserName,
		Account:   form.Account,
		DeptID:    form.DeptID,
		JGMC:      jgmc,
		JGDM:      jgdm,
		Status:    0,
	}

	err = models.CreateGroupMem(member)
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// @Tags 获取团队成员列表
// @Summary POST方法 获取团队成员列表
// @Accept application/json
// @Produce  application/json
// @Param data body Member true "只传ID"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":"返回值为传参的对象数组"}
// @Router /api/v1/member/list [post]
func GetMembers(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Member
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	members, err := models.GetMembers(form.GroupID)
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, members)

}

// @Tags 退出团队或删人
// @Summary POST方法 退出团队或删人
// @Accept application/json
// @Produce  application/json
// @Param data body Member true "只传ID和account"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":nil}
// @Router /api/v1/member/dropout [post]
func DropOut(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Member
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	if err := models.DropOut(form.GroupID, form.Account); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}
