package v1

import (
	"github.com/gin-gonic/gin"
	"knowledgebase/pkg/app"
	"knowledgebase/pkg/e"
	"knowledgebase/pkg/util"
	"log"
	"net/http"
	"strconv"
)

func Login1(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form LoginForm
	)
	resp := map[string]interface{}{}
	result := map[string]interface{}{}
	userInfo := map[string]interface{}{}
	depart := map[string]interface{}{}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	if resp["result"] != nil {
		result = resp["result"].(map[string]interface{})
		userInfo = result["userInfo"].(map[string]interface{})
		departs := result["departs"].([]interface{})
		depart = departs[0].(map[string]interface{})
	}
	//}

	//todo:internet test

	log.Println(strconv.FormatFloat(1.3706e+10, 'f', -1, 64))

	token1, _ := util.GenerateToken(form.Username, strconv.FormatFloat(1.3706e+10, 'f', -1, 64))
	token2, _ := util.GenerateToken(form.Username, strconv.FormatFloat(13706130800, 'f', -1, 64))
	if form.Username == "test" {
		resp["success"] = "True"
		resp["message"] = "操作成功"
		result["token"] = token1
		userInfo["id"] = "test"
		userInfo["username"] = "张三"
		userInfo["userAccount"] = "test"
		depart["id"] = "13706130900"
		depart["departName"] = "XX市XX区信息中心"
		depart["parentId"] = 13706130000
		depart["jgmc"] = "烟台市莱山区税务局"
		depart["jgdm"] = 13706130000
	}
	if form.Username == "test1" {
		resp["success"] = "True"
		resp["message"] = "操作成功"
		result["token"] = token2
		userInfo["id"] = "test1"
		userInfo["username"] = "王五"
		userInfo["userAccount"] = "test1"
		depart["id"] = 13706001800
		depart["departName"] = "XX市信息中心"
		depart["parentId"] = 13706000000
		depart["jgmc"] = "烟台市税务局"
		depart["jgdm"] = 13706000000
	}

	type role struct {
		RoleID   string `json:"role_id"`
		RoleName string `json:"role_name"`
	}

	roles := make([]*role, 2)
	roles = []*role{
		&role{
			RoleID:   "10010",
			RoleName: "市局普通用户",
		},
		&role{
			RoleID:   "10011",
			RoleName: "全局管理员",
		},
	}

	data := map[string]interface{}{
		"success":     resp["success"],
		"message":     resp["message"],
		"token":       result["token"],
		"userid":      userInfo["id"],
		"username":    userInfo["username"],
		"userAccount": userInfo["userAccount"],
		"departID":    depart["id"],
		"departName":  depart["departName"],
		"parentId":    depart["parentId"],
		"jgmc":        depart["jgmc"],
		"jgdm":        depart["jgdm"],
		//todo:internet test
		"userRole": roles,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func GetUserList1(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
	)

	users := make([]map[string]string, 2)

	appG.Response(http.StatusOK, e.SUCCESS, users)
}
