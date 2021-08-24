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

type Group struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Creator string `json:"creator"`
	Account string `json:"account"`
	Type    string `json:"type"`
}

type Userinfo struct {
	Account string `json:"account"`
}

// 创建群
// @Tags 创建团队
// @Summary POST方法 创建团队
// @Accept application/json
// @Produce  application/json
// @Param data body Group true "ID(不用传) , 团队名称 , 创建人 , 创建人账号 , 类型（0临时团队，1长期团队）"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":nil}
// @Router /api/v1/group/create [post]
func CreateGroup(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Group
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	ok, err := models.GroupIsExit(form.Name)
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if ok {
		log.Println("该团队名称已被使用")
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	var user *models.SysUser

	// 获取当前用户信息
	Clients.Range(func(k, v interface{}) bool {
		if v.(*models.SysUser).UserAccount == form.Account {
			user = v.(*models.SysUser)
			return false
		}
		return true
	})

	jgdm, jgmc := util.SWJG(strconv.FormatFloat(user.DepID, 'f', -1, 64)[:7])

	group := &models.Group{
		ID:         "gp-" + util.RandomString(29),
		Name:       form.Name,
		Creator:    form.Creator,
		Account:    form.Account,
		JGMC:       jgmc,
		JGDM:       jgdm,
		Type:       form.Type,
		UserCnt:    0,
		Del:        0,
		Uptime:     time.Now().Format("2006-01-02 15:04:05"),
		LastChated: time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := models.CreateGroup(group); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	//jgdm , jgmc := util.SWJG(strconv.FormatFloat(user.DepID,'f',-1,64)[:7])

	member := models.Member{
		ID:        "mb-" + util.RandomString(29),
		GroupID:   group.ID,
		GroupName: group.Name,
		UserName:  user.Username,
		Account:   user.UserAccount,
		DeptID:    strconv.FormatFloat(user.DepID, 'f', -1, 64),
		JGMC:      jgmc,
		JGDM:      jgdm,
		Status:    0,
	}

	if err := models.CreateGroupMem(member); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	// 成员+1
	models.AddUserCnt(group.ID)

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// 修改群
// @Tags 修改团队
// @Summary POST方法 修改团队
// @Accept application/json
// @Produce  application/json
// @Param data body Group true "ID , 团队名称 , 创建人 , 创建人账号 , 类型（0临时团队，1长期团队）"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":nil}
// @Router /api/v1/group/edit [post]
func EditGroup(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Group
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	if form.Name != "" {
		ok, err := models.GroupIsExit(form.Name)
		if err != nil {
			log.Println(err)
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}

		if ok {
			log.Println("该团队名称已被使用")
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
	}

	group := &models.Group{
		ID:         form.ID,
		Name:       form.Name,
		Creator:    form.Creator,
		Account:    form.Account,
		Type:       form.Type,
		LastChated: time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := models.EditGroup(group); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// 获取当前用户所在团队列表
// @Tags 获取当前用户所在团队列表
// @Summary POST方法 获取当前用户所在团队列表
// @Accept application/json
// @Produce  application/json
// @Param data body Userinfo true "当前账号"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":{[{"groupId":"qwertgh","groupName":"第一群"}]}}
// @Router /api/v1/group/groups [post]
func GetGroups(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Userinfo
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	groups, err := models.MyJoinedGroup(form.Account)
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, groups)

}

//我管理的群
// @Tags 获取我创建的团队列表
// @Summary POST方法 获取我创建的团队列表
// @Accept application/json
// @Produce  application/json
// @Param data body models.Group true "只用传account"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":"返回值同传参列表，对象数组"}
// @Router /api/v1/group/mygroups [post]
func MyGroups(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Userinfo
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	groups, err := models.GetGroups(form.Account)
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, groups)
}

// 删除群
// @Tags 删除团队
// @Summary POST方法 删除团队
// @Accept application/json
// @Produce  application/json
// @Param data body Group true "只传id"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":nil}
// @Router /api/v1/group/del [post]
func DelGroup(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Group
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	if err := models.DelGroup(form.ID); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// 查询团队
// @Tags 查询团队
// @Summary POST方法 查询团队
// @Accept application/json
// @Produce  application/json
// @Param data body Group true "只传name"
// @Success 200 {string} string "{"code":200,"msg":"ok","data":"返回值同传参列表，对象数组"}
// @Router /api/v1/group/select [post]
func SeleteGroup(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
		form Group
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	group, err := models.SeleteGroup(form.Name)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.Select_ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, group)

}
