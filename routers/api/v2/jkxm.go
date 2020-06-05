package v2

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/app"
	"NULL/knowledgebase/pkg/e"
	"NULL/knowledgebase/pkg/export"
	"NULL/knowledgebase/pkg/jkxm"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

type ShZjForm struct {
	XmDm string   `json:"xm_dm"`
	Id   []string `json:"id"`
}

// 导入监控项目
func ImpJkxm(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		session = sessions.Default(c)
		xmDm    = c.Query("xm_dm")
	)
	userAccount := session.Get("userAccount").(string)
	username := session.Get("username").(string)
	departID := session.Get("departID").(string)
	departName := session.Get("departName").(string)
	userInfo := map[string]string{
		"userAccount": userAccount,
		"username":    username,
		"departID":    departID,
		"departName":  departName,
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	impMsg := jkxm.ImpJkxm(file, xmDm, userInfo)
	appG.Response(http.StatusOK, e.SUCCESS, impMsg)
}

// 审核监控项目
func ShJkxm(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		session = sessions.Default(c)
		form    ShZjForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	userAccount := session.Get("userAccount").(string)
	username := session.Get("username").(string)
	departID := session.Get("departID").(string)
	departName := session.Get("departName").(string)
	userRole := session.Get("userRole").(string)
	if !strings.Contains(userRole, form.XmDm) {
		appG.Response(http.StatusInternalServerError, e.ERROR,
			"该用户无此监控指标操作权限!")
		return
	}
	shMap := map[string]string{}
	if strings.Contains(c.Request.URL.Path, "lrsh") {
		shMap = map[string]string{
			"shr_account":  userAccount,
			"shr_name":     username,
			"shr_depid":    departID,
			"shr_deptname": departName,
			"shrq":         time.Now().Format("2006-01-02 15:04:05"),
			"shbz":         "Y",
		}
	}
	if strings.Contains(c.Request.URL.Path, "zjsh") {
		shMap = map[string]string{
			"zjshr_account":  userAccount,
			"zjshr_name":     username,
			"zjshr_depid":    departID,
			"zjshr_deptname": departName,
			"zjshrq":         time.Now().Format("2006-01-02 15:04:05"),
			"zjshbz":         "Y",
		}
	}
	var data []string
	for _, id := range form.Id {
		if err := models.ShJkxm(form.XmDm, id, shMap); err != nil {
			log.Println(err)
			data = append(data, id)
			continue
		}
	}
	if len(data) == 0 {
		data = append(data, "审核成功!")
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 终结监控项目
func ZjJkxm(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		session = sessions.Default(c)
		form    ShZjForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	userRole := session.Get("userRole").(string)
	if !strings.Contains(userRole, form.XmDm) {
		appG.Response(http.StatusInternalServerError, e.ERROR,
			"该用户无此监控指标操作权限!")
		return
	}
	userAccount := session.Get("userAccount").(string)
	username := session.Get("username").(string)
	departID := session.Get("departID").(string)
	departName := session.Get("departName").(string)
	zjMap := map[string]string{
		"zjr_account":  userAccount,
		"zjr_name":     username,
		"zjr_depid":    departID,
		"zjr_deptname": departName,
		"zjrq":         time.Now().Format("2006-01-02 15:04:05"),
		"zjbz":         "Y",
	}
	var data []string
	for _, id := range form.Id {
		if err := models.ZjJkxm(form.XmDm, id, zjMap); err != nil {
			log.Println(err)
			data = append(data, id)
			continue
		}
	}
	if len(data) == 0 {
		data = append(data, "终结成功!")
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 根据审核标志获取对应项目列表
func GetJkxmByShbz(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		nsrsbh = c.Query("nsrsbh")
		xmDm   = c.Query("xm_dm")
		shbz   = c.Query("shbz")
		squery string
	)
	if strings.Contains(c.Request.URL.Path, "lrsh") {
		squery = fmt.Sprintf(
			`select * from %s where shbz='%s'`, xmDm, shbz)
	}
	if strings.Contains(c.Request.URL.Path, "zjsh") {
		if c.Query("flag") == "" {
			squery = fmt.Sprintf(
				`select * from %s where zjshbz='%s' and zjbz='Y'`, xmDm, shbz)
		} else {
			squery = fmt.Sprintf(
				`select * from %s where zjshbz='N' and shbz='Y'`, xmDm)
		}
	}
	if nsrsbh != "" {
		squery += fmt.Sprintf(" and nsrsbh='%s'", nsrsbh)
	}
	data, err := models.QueryData(squery)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 下载根据审核标志获取对应项目列表
func DownloadJkxmByShbz(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		nsrsbh = c.Query("nsrsbh")
		xmDm   = c.Query("xm_dm")
		shbz   = c.Query("shbz")
	)
	squery := fmt.Sprintf(
		`select * from %s where zjshbz='%s' and shbz='Y'`, xmDm, shbz)
	if nsrsbh != "" {
		squery += fmt.Sprintf(" and nsrsbh='%s'", nsrsbh)
	}
	records, err := models.QueryData(squery)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	var url = "该监控项目没有异常数据"
	if len(records) > 0 {
		fileName := models.GetJkxmMcByDm(xmDm)
		url, err = export.WriteIntoExcel(fileName, records)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, url)
}

// 根据终结标志获取对应项目列表
func GetJkxmByZjbz(c *gin.Context) {
	appG := app.Gin{C: c}
	squery := fmt.Sprintf(`select * from %s where zjbz='%s' and shbz='Y'`,
		c.Query("xm_dm"), c.Query("zjbz"))
	nsrsbh := c.Query("nsrsbh")
	if nsrsbh != "" {
		squery += fmt.Sprintf(" and nsrsbh='%s'", nsrsbh)
	}
	data, err := models.QueryData(squery)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type Resp struct {
	*models.JkxmMcdm
	Cnt int `json:"cnt"`
}

// 获取所有监控项目异常(未终结审核)数量
func GetJkxms(c *gin.Context) {
	appG := app.Gin{C: c}
	nsrsbh := c.Query("nsrsbh")
	mcdms, err := models.GetJkxmMcdms()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	var nsrmc string
	var resps = make([]*Resp, 0)
	if len(mcdms) > 0 {
		for _, mcdm := range mcdms {
			var cond string
			if nsrsbh != "" {
				cond = fmt.Sprintf("nsrsbh='%s'", nsrsbh)
				if nsrmc == "" {
					query := fmt.Sprintf(
						`select distinct nsrmc from %s where nsrsbh='%s'`,
						mcdm.Dm, nsrsbh)
					nsrmcs, _ := models.QueryData(query)
					if len(nsrmcs) > 0 {
						nsrmc = nsrmcs[0]["nsrmc"]
					}
				}
			} else {
				cond = "nsrsbh like '%'"
				nsrmc = "合计"
			}
			cnt := models.CountJkxms(mcdm.Dm, cond)
			resps = append(resps, &Resp{mcdm, cnt})
		}
	}
	data := map[string]interface{}{
		"nsrshb": nsrsbh,
		"nsrmc":  nsrmc,
		"list":   resps,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 获取监控项目名称代码
func GetJkxmMcdms(c *gin.Context) {
	appG := app.Gin{C: c}
	data, err := models.GetJkxmMcdms()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
