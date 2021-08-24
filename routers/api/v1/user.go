package v1

import (
	"crypto/tls"
	"encoding/json"
	"github.com/mozillazg/request"
	"github.com/parnurzeal/gorequest"
	"knowledgebase/models"
	"knowledgebase/pkg/app"
	"knowledgebase/pkg/e"
	"knowledgebase/pkg/logging"
	"knowledgebase/pkg/setting"
	"knowledgebase/pkg/util"
	"log"
	"net/url"
	"strconv"
	"strings"

	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserForm struct {
	UserAccount string `json:"user_account"`
	Username    string `json:"username"`
	DepartID    string `json:"depart_id"`
	DepartName  string `json:"depart_name"`
	UserRole    string `json:"user_role"`
}
type Users struct {
	Deptid    string `json:"deptid"`
	Checked   bool   `json:"checked"`
	BeginTime string `json:"beginTime"`
	EndTime   string `json:"endTime"`
	Status    int    `json:"status"`
	PageNo    int    `json:"pageNo"`
	pageSize  int    `json:"page_size"`
	Page      int    `json:"page"`
}

var Usertoken = make(map[string]string)

// 用户登录
func Login(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form LoginForm
	)
	resp := map[string]interface{}{}
	//result := map[string]interface{}{}
	//userInfo := map[string]interface{}{}
	//depart := map[string]interface{}{}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	resp, err := login(form.Username, form.Password)
	if err != nil {
		logging.Error(err)
		data := map[string]string{
			"message": "登录失败，请检查账号密码是否输入正确，若无误，请联系管理员",
		}
		appG.Response(http.StatusOK, e.ERROR, data)
		return
	}

	if resp["code"].(float64) != 200 {
		appG.Response(http.StatusInternalServerError, e.ERROR, map[string]string{
			"message": "用户不存在/密码错误",
		})
		return
	}

	//token , err := util.GenerateToken(form.Username,strconv.FormatFloat(resp["dept_id"].(float64), 'f', -1, 64))
	token, err := util.GenerateToken(form.Username, resp["dept_id"].(string))
	if err != nil {
		logging.Error("生成token失败")
		log.Println("生成token失败")
		appG.Response(http.StatusInternalServerError, e.ERROR, "create token failed")
	}

	cookie := resp["token"].(string)
	//UserToken[resp["userName"]] = cookie
	resp_getinfo, err := getInfo(cookie)
	if err != nil {
		logging.Error(err)
		data := map[string]string{
			"message": "返回角色信息失败，请联系管理员",
		}
		appG.Response(http.StatusInternalServerError, e.ERROR, data)
		return
	}

	if resp["code"].(float64) != 200 {
		data := map[string]string{
			"message": "返回角色信息失败，请联系管理员",
		}
		appG.Response(http.StatusInternalServerError, e.ERROR, data)
		return
	}

	type role struct {
		RoleID   interface{} `json:"role_id"`
		RoleName interface{} `json:"role_name"`
	}

	var roleslice = make([]*role, 0)

	for _, i := range resp_getinfo["user"].(map[string]interface{})["roles"].([]interface{}) {

		var r = &role{
			RoleID:   i.(map[string]interface{})["roleId"],
			RoleName: i.(map[string]interface{})["roleName"],
		}

		roleslice = append(roleslice, r)

	}

	data := map[string]interface{}{
		"success":     "success",
		"message":     resp["msg"],
		"token":       token,
		"userid":      "",
		"username":    resp["nickName"],
		"userAccount": resp["userName"],
		"departID":    resp["dept_id"],
		"departName":  resp["dept_name"],
		"parentId":    resp["parent_id"],
		"userRole":    roleslice,
		"jgdm":        resp["SWJG_DM"],
		"jgmc":        resp["SWJG"],
	}

	Usertoken[form.Username] = cookie

	appG.Response(http.StatusOK, e.SUCCESS, data)

}

// 获取用户信息,存入session
func UserInfo(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		session = sessions.Default(c)
		form    UserForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	session.Set("userAccount", form.UserAccount)
	session.Set("username", form.Username)
	session.Set("departID", form.DepartID)
	session.Set("departName", form.DepartName)
	session.Set("userRole", form.UserRole)
	if err := session.Save(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 获取纳税人信息
func NsrInfo(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		nsrsbh = c.Query("nsrsbh")
		nsrmc  = "无此纳税人税务登记信息"
	)
	//todo:internet test
	/*if nsrsbh == "123" {
		nsrmc = "测试名称1"
	}
	if nsrsbh == "456" {
		nsrmc = "测试名称2"
	}*/
	pd, err := models.GetConfigSql("取纳税人名称")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	sql := fmt.Sprintf(pd.XmSql, nsrsbh)
	records, err := models.QueryData(sql)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	if len(records) > 0 {
		nsrmc = records[0]["NSRMC"]
	}
	appG.Response(http.StatusOK, e.SUCCESS, nsrmc)
}

// 代理转发智税平台用户登录
func Rlogin(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form LoginForm
	)
	resp := map[string]interface{}{}
	result := map[string]interface{}{}
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
			appG.Response(http.StatusOK, e.SUCCESS, result)
			return
		}
		appG.Response(http.StatusOK, e.ERROR, nil)
	}
}

//获取数据中台路由表
func GetRoutes(c *gin.Context) {
	var appG = app.Gin{C: c}
	var IeHeader = `Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Win64; x64; Trident/4.0; .NET CLR 2.0.limit727; SLCC2; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET4.0C; .NET4.0E)`

	token := c.Query("token")
	// 引入"crypto/tls":解决golang https请求提示x509: certificate signed by unknown authority
	ts := &tls.Config{InsecureSkipVerify: true}
	url := setting.AppSetting.LoginUrl + "/jeecg-boot/sys/permission/getUserPermissionByToken" +
		"?token=" + token
	_, body, errs := gorequest.New().TLSClientConfig(ts).Get(url).
		Set("Content-Type", "application/json").
		Set("X-Access-Token", token).Set("Accept-Language", "zh-CN").
		Set("Accept", "*/*").Set("User-Agent", IeHeader).End()
	if len(errs) > 0 {
		data := fmt.Sprintf("Get routes err:%v", errs[0])
		appG.Response(http.StatusOK, e.ERROR, data)
		return
	} else {
		if !strings.Contains(body, "查询成功") {
			appG.Response(http.StatusOK, e.ERROR, body)
			return
		}
		resp := make(map[string]interface{})
		result := map[string]interface{}{}
		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			data := fmt.Sprintf("unmarshall body error:%v", err)
			appG.Response(http.StatusOK, e.ERROR, data)
			return
		}
		if resp["result"] != nil {
			result = resp["result"].(map[string]interface{})
			menu := make(map[string]interface{})
			for _, m := range result["menu"].([]interface{}) {
				mMap := m.(map[string]interface{})
				if mMap["path"].(string) == "/shuiwuyewu" {
					menu = map[string]interface{}{
						"redirect":  mMap["redirect"],
						"path":      mMap["path"],
						"component": mMap["component"],
						"route":     mMap["route"],
					}
					for _, c := range mMap["children"].([]interface{}) {
						cMap := c.(map[string]interface{})
						if cMap["path"].(string) == "/zxfxsm" {
							menu["children"] = []map[string]interface{}{
								{
									"path":      cMap["path"],
									"component": cMap["component"],
									"route":     cMap["route"],
									"children":  cMap["children"],
								},
							}
						}
					}
				}
			}
			appG.Response(http.StatusOK, e.SUCCESS,
				map[string]interface{}{
					"allAuth": result["allAuth"],
					"auth":    result["auth"],
					"menu": []map[string]interface{}{
						menu,
					},
				})
			return
		}
		appG.Response(http.StatusOK, e.ERROR, nil)
	}
}

func Zhpt_login(c *gin.Context) {

	type token_param struct {
		Token string
	}

	var (
		appG = app.Gin{C: c}
		form token_param
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	res, err := getInfo(form.Token)
	if err != nil {
		logging.Error(err)
		data := map[string]string{
			"message": "智慧平台转登录返回角色信息失败，请联系管理员",
		}
		appG.Response(http.StatusOK, e.ERROR, data)
	}

	if res["code"].(float64) != 200 {

		appG.Response(http.StatusInternalServerError, e.ERROR, "自动失败,平台登录超时")
		return
	}

	user := res["user"].(map[string]interface{})
	dept := user["dept"].(map[string]interface{})

	type role struct {
		RoleID   interface{} `json:"role_id"`
		RoleName interface{} `json:"role_name"`
	}
	var roleslice = make([]*role, 0)

	for _, i := range res["user"].(map[string]interface{})["roles"].([]interface{}) {

		var r = &role{
			RoleID:   i.(map[string]interface{})["roleId"],
			RoleName: i.(map[string]interface{})["roleName"],
		}

		roleslice = append(roleslice, r)

	}

	token, err := util.GenerateToken(user["userName"].(string), strconv.FormatFloat(user["deptId"].(float64), 'f', -1, 64))
	//token , err := util.GenerateToken(user["userName"].(string),strconv.FormatFloat(user["deptId"].(float64), 'f', -1, 64))
	if err != nil {
		log.Println("生成token失败")
		appG.Response(http.StatusOK, e.ERROR, "智慧平台转登录token生成失败")
		return
	}

	data := map[string]interface{}{
		"success":     "success",
		"message":     res["msg"],
		"token":       token,
		"userid":      "",
		"username":    user["nickName"],
		"userAccount": user["userName"],
		"departID":    dept["deptId"],
		"departName":  dept["deptName"],
		"parentId":    dept["parentId"],
		"userRole":    roleslice,
		"jgdm":        user["swjgdm"],
		"jgmc":        user["swjgmc"],
	}

	appG.Response(http.StatusOK, e.SUCCESS, data)

}

// @Tags 用户列表
// @Summary POST方法 用户列表
// @Accept application/json
// @Produce  application/json
// @Param data body Users true "部门ID , 本级(默认为false，如果是查同级别，则为true) ,状态（0正常，1停用） ,开始时间 , 结束时间 "
// @Success 200 {string} string "{"code":200,"msg":"ok","data":{"list":[{"account":"test","username":"张三","deptname":"XX区信息中心","deptid":"13706130900","jgmc":"烟台市莱山区税务局","jgdm":"13706130000"}],"total":100}}"
// @Router /api/v1/userinfo [post]
func GetUserList(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
		form Users
	)

	account := util.GetLoginID("", c)
	cookies := Usertoken[account]
	resp, err := userlist(cookies, form)
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if resp["code"].(float64) != 200 {
		log.Println(resp)
		appG.Response(http.StatusInternalServerError, e.ERROR, map[string]interface{}{
			"message": resp["data"],
		})
		return
	}

	type userstruct struct {
		Account  string `json:"account"`
		Username string `json:"username"`
		Deptname string `json:"deptname"`
		DeptID   string `json:"deptId"`
		JGMC     string `json:"jgmc"`
		JGDM     string `json:"jgdm"`
	}

	users := resp["rows"].([]map[string]interface{})
	userslice := make([]*userstruct, 0)
	for _, user := range users {

		deptid := strconv.FormatFloat(user["deptId"].(float64), 'f', -1, 64)
		jgdm, jgmc := util.SWJG(deptid)

		usersingle := &userstruct{
			Account:  user["userName"].(string),
			Username: user["nickName"].(string),
			Deptname: user["deptName"].(string),
			DeptID:   deptid,
			JGDM:     jgdm,
			JGMC:     jgmc,
		}

		userslice = append(userslice, usersingle)

	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"list": userslice, "total": resp["total"].(int),
	})

}

func login(username, password string) (map[string]interface{}, error) {

	c := new(http.Client)
	req := request.NewRequest(c)
	req.Json = map[string]string{
		"username": username,
		"password": password,
	}

	req.Headers = map[string]string{
		"tjtoken": util.RandomString(30),
	}

	log.Println("username", username)
	log.Println("password", password)

	p := url.Values{"username": {username}, "password": {password}}

	res, err := req.Post(setting.AppSetting.LoginUrl + "/loginNoCode?" + p.Encode())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp := make(map[string]interface{})
	rrr, err := res.Content()
	err = json.Unmarshal(rrr, &resp)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(resp)

	return resp, nil

}

func getInfo(cookies string) (map[string]interface{}, error) {

	c := new(http.Client)
	req := request.NewRequest(c)
	req.Cookies = map[string]string{
		"sidebarStatus": "0",
		"Admin-Token":   cookies,
	}

	req.Headers = map[string]string{
		"Authorization": "Bearer " + cookies,
		"tjtoken":       util.RandomString(30),
	}

	res, err := req.Get(setting.AppSetting.LoginUrl + "/getInfo")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp := make(map[string]interface{})
	rrr, err := res.Content()
	err = json.Unmarshal(rrr, &resp)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(resp)

	return resp, nil

}

func userlist(cookies string, form Users) (map[string]interface{}, error) {

	c := new(http.Client)
	req := request.NewRequest(c)
	req.Cookies = map[string]string{
		"sidebarStatus": "0",
		"Admin-Token":   cookies,
	}

	req.Headers = map[string]string{
		"Authorization": "Bearer " + cookies,
		"tjtoken":       util.RandomString(30),
	}

	deptid := form.Deptid[:7] + "0000"

	p := url.Values{"pageNum": {strconv.Itoa(form.PageNo)}, "pageSize": {strconv.Itoa(form.pageSize)},
		"deptId": {deptid}, "checked": {strconv.FormatBool(form.Checked)},
		"beginTime": {form.BeginTime}, "endTime": {form.EndTime},
		"page": {strconv.Itoa(form.Page)}, "status": {strconv.Itoa(form.Status)}}

	res, err := req.Get(setting.AppSetting.LoginUrl + "/system/user/list?" + p.Encode())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp := make(map[string]interface{})
	rrr, err := res.Content()
	err = json.Unmarshal(rrr, &resp)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return resp, nil
}
