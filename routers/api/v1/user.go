package v1

import (
	"NULL/knowledgebase/pkg/app"
	"NULL/knowledgebase/pkg/e"
	"NULL/knowledgebase/pkg/setting"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
	"net/http"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//用户登录
func Login(c *gin.Context) {
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
	// 引入"crypto/tls":解决golang https请求提示x509: certificate signed by unknown authority
	ts := &tls.Config{InsecureSkipVerify: true}
	pMap := map[string]string{
		"username": form.Username,
		"password": form.Password,
	}
	_, body, errs := gorequest.New().TLSClientConfig(ts).
		Post(setting.AppSetting.LoginUrl + "/jeecg-boot/sys/login").
		Type(gorequest.TypeJSON).SendMap(pMap).End()
	if len(errs) > 0 {
		data := fmt.Sprintf("login err:%v", errs[0])
		appG.Response(http.StatusOK, e.ERROR, data)
		return
	} else {
		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			data := fmt.Sprintf("unmarshall body error:%v", err)
			appG.Response(http.StatusOK, e.ERROR, data)
			return
		}
		if resp["result"] != nil {
			result = resp["result"].(map[string]interface{})
			userInfo = result["userInfo"].(map[string]interface{})
			departs := result["departs"].([]interface{})
			depart = departs[0].(map[string]interface{})
		}
	}

	//internet test
	/*token, _ := util.GenerateToken(form.Username, form.Password)
	if form.Username == "test" {
		resp["success"] = "True"
		resp["message"] = "登陆成功"
		result["token"] = token
		userInfo["id"] = "13706002531"
		userInfo["username"] = "张三"
		userInfo["userAccount"] = "test"
		depart["id"] = "13706130900"
		depart["departName"] = "XX市XX区信息中心"
		depart["parentId"] = "13706130000"
	}
	if form.Username == "test1" {
		resp["success"] = "True"
		resp["message"] = "登陆成功"
		result["token"] = token
		userInfo["id"] = "test1"
		userInfo["username"] = "王五"
		userInfo["userAccount"] = "test1"
		depart["id"] = "13706001800"
		depart["departName"] = "XX市信息中心"
		depart["parentId"] = "13706000000"
	}*/

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
	}

	appG.Response(http.StatusOK, e.SUCCESS, data)
}
