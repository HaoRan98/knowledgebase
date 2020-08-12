package v2

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/app"
	"NULL/knowledgebase/pkg/cron"
	"NULL/knowledgebase/pkg/e"
	"NULL/knowledgebase/pkg/export"
	"NULL/knowledgebase/pkg/util"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type JbzxForm struct {
	Nsrsbh string `json:"nsrsbh"`
	Nsrmc  string `json:"nsrmc"`
}

// 即办注销监控项目
func JbzxJkxm(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		session = sessions.Default(c)
		form    JbzxForm
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
	jbzx := models.JkxmJbzx{
		ID:          "JBZX-" + util.RandomString(15),
		Nsrsbh:      form.Nsrsbh,
		Nsrmc:       form.Nsrmc,
		JbrAccount:  userAccount,
		JbrName:     username,
		JbrDepid:    departID,
		JbrDeptname: departName,
		Qrrq:        time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.JbZxJkxm(&jbzx); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 统计工作量
func CountWorkload(c *gin.Context) {
	var (
		appG  = app.Gin{C: c}
		fqrqq = c.Query("fqrqq")
		fqrqz = c.Query("fqrqz")
		zjrqq = c.Query("zjrqq")
		zjrqz = c.Query("zjrqz")
	)
	if fqrqq == "" && fqrqz == "" {
		fqrqq, fqrqz = "2020-01-01", "2099-12-31"
	}
	var cond = fmt.Sprintf(
		"fqrq>='%s 00:00:00' and fqrq<='%s 23:59:59'", fqrqq, fqrqz)
	if zjrqq != "" && zjrqz != "" {
		cond += fmt.Sprintf(
			" and zjrq>='%s 00:00:00' and zjrq<='%s 23:59:59'", zjrqq, zjrqz)
	}
	mcdms, err := models.GetJkxmMcdms()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	var total, resolve, unsolved int
	if len(mcdms) > 0 {
		for _, mcdm := range mcdms {
			total += models.CountJkxmsToal(mcdm.Dm, cond)
			resolve += models.CountJkxmRsolved(mcdm.Dm, cond)
			unsolved += models.CountJkxmsUnsolved(mcdm.Dm, cond)
		}
	}
	data := map[string]interface{}{
		"total":    total,    //疑点指标总数量
		"resolve":  resolve,  //已经解除疑点数量
		"unsolved": unsolved, //尚未解除的疑点数量
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type CountResp struct {
	*models.JkxmMcdm
	TotalCnt     int `json:"total_cnt"`
	ResolveCnt   int `json:"resolve_cnt"`
	UnresolveCnt int `json:"unresolve_cnt"`
}

// 汇总所有监控项目异常数量
func CountJkxms(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		nsrsbh = c.Query("nsrsbh")
		fqrqq  = c.Query("fqrqq")
		fqrqz  = c.Query("fqrqz")
		zjrqq  = c.Query("zjrqq")
		zjrqz  = c.Query("zjrqz")
	)
	mcdms, err := models.GetJkxmMcdms()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	var nsrmc string
	var resps = make([]*CountResp, 0)
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
			if fqrqq == "" && fqrqz == "" {
				fqrqq, fqrqz = "2020-01-01", "2099-12-31"
			}
			cond += fmt.Sprintf(
				" and fqrq>='%s 00:00:00' and fqrq<='%s 23:59:59'", fqrqq, fqrqz)
			if zjrqq != "" && zjrqz != "" {
				cond += fmt.Sprintf(
					" and zjrq>='%s 00:00:00' and zjrq<='%s 23:59:59'", zjrqq, zjrqz)
			}
			resps = append(resps, &CountResp{
				JkxmMcdm:     mcdm,
				TotalCnt:     models.CountJkxmsToal(mcdm.Dm, cond),
				ResolveCnt:   models.CountJkxmRsolved(mcdm.Dm, cond),
				UnresolveCnt: models.CountJkxmsUnsolved(mcdm.Dm, cond),
			})
		}
	}
	data := map[string]interface{}{
		"nsrsbh": nsrsbh,
		"nsrmc":  nsrmc,
		"list":   resps,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type Resp struct {
	*models.JkxmMcdm
	Cnt int `json:"cnt"`
}

// 获取所有监控项目异常总数量
func GetJkxmsTotal(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		nsrsbh = c.Query("nsrsbh")
		fqrqq  = c.Query("fqrqq")
		fqrqz  = c.Query("fqrqz")
		zjrqq  = c.Query("zjrqq")
		zjrqz  = c.Query("zjrqz")
	)
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
			if fqrqq == "" && fqrqz == "" {
				fqrqq, fqrqz = "2020-01-01", "2099-12-31"
			}
			cond += fmt.Sprintf(
				" and fqrq>='%s 00:00:00' and fqrq<='%s 23:59:59'", fqrqq, fqrqz)
			if zjrqq != "" && zjrqz != "" {
				cond += fmt.Sprintf(
					" and zjrq>='%s 00:00:00' and zjrq<='%s 23:59:59'", zjrqq, zjrqz)
			}
			cnt := models.CountJkxmsToal(mcdm.Dm, cond)
			resps = append(resps, &Resp{mcdm, cnt})
		}
	}
	data := map[string]interface{}{
		"nsrsbh": nsrsbh,
		"nsrmc":  nsrmc,
		"list":   resps,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 获取所有监控项目异常(已经解除疑点)数量
func GetJkxmsResolve(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		nsrsbh = c.Query("nsrsbh")
		fqrqq  = c.Query("fqrqq")
		fqrqz  = c.Query("fqrqz")
		zjrqq  = c.Query("zjrqq")
		zjrqz  = c.Query("zjrqz")
	)
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
			if fqrqq == "" && fqrqz == "" {
				fqrqq, fqrqz = "2020-01-01", "2099-12-31"
			}
			cond += fmt.Sprintf(
				" and fqrq>='%s 00:00:00' and fqrq<='%s 23:59:59'", fqrqq, fqrqz)
			if zjrqq != "" && zjrqz != "" {
				cond += fmt.Sprintf(
					" and zjrq>='%s 00:00:00' and zjrq<='%s 23:59:59'", zjrqq, zjrqz)
			}
			cnt := models.CountJkxmRsolved(mcdm.Dm, cond)
			resps = append(resps, &Resp{mcdm, cnt})
		}
	}
	data := map[string]interface{}{
		"nsrsbh": nsrsbh,
		"nsrmc":  nsrmc,
		"list":   resps,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 获取所有监控项目异常(未终结审核)数量/尚未解除的疑点数量
func GetJkxmsUnsolved(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		nsrsbh = c.Query("nsrsbh")
		fqrqq  = c.Query("fqrqq")
		fqrqz  = c.Query("fqrqz")
		zjrqq  = c.Query("zjrqq")
		zjrqz  = c.Query("zjrqz")
	)
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
			if fqrqq == "" && fqrqz == "" {
				fqrqq, fqrqz = "2020-01-01", "2099-12-31"
			}
			cond += fmt.Sprintf(
				" and fqrq>='%s 00:00:00' and fqrq<='%s 23:59:59'", fqrqq, fqrqz)
			if zjrqq != "" && zjrqz != "" {
				cond += fmt.Sprintf(
					" and zjrq>='%s 00:00:00' and zjrq<='%s 23:59:59'", zjrqq, zjrqz)
			}
			cnt := models.CountJkxmsUnsolved(mcdm.Dm, cond)
			resps = append(resps, &Resp{mcdm, cnt})
		}
	}
	var flag = false
	if nsrsbh != "" {
		flag = models.IsNsrsbhExist(nsrsbh)
	}
	data := map[string]interface{}{
		"nsrsbh": nsrsbh,
		"nsrmc":  nsrmc,
		"list":   resps,
		"flag":   flag, //是否已启动即办注销流程
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 成果统计
func Zxtj(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		rqq  = c.Query("rqq")
		rqz  = c.Query("rqz")
	)
	if rqq == "" && rqz == "" {
		ldDate := models.NewLdDate("")
		rqq, rqz = ldDate.LdQrqBn, ldDate.LdCurrentDate
	}
	var cond = fmt.Sprintf(
		"qrrq>='%s 00:00:00' and qrrq<='%s 23:59:59'", rqq, rqz)
	jbzxs, err := models.GetJbzxs(cond)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	//todo:internet test
	//sql := fmt.Sprintf(
	//	"select * from jkxm_gt3 where zxsj>='%s 00:00:00' and zxsj<='%s 23:59:59'",
	//	rqq, rqz)
	pd, err := models.GetConfigSql("注销")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR,
			fmt.Sprintf("获取金三注销名单查询语句失败:%v", err))
		return
	}
	sql := fmt.Sprintf(pd.XmSql, rqq, rqz)
	gt3Zxs, err := models.QueryData(sql)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"tjsj":     time.Now().Format("2006-01-02 15:04:05"),
		"zgswjg":   "国家税务总局烟台经济技术开发区税务局",
		"jbzx_cnt": len(jbzxs),
		"gt3_cnt":  len(gt3Zxs),
	})
}

// 确认启动即办注销流程户数明细
func JbzxJkxmMx(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		rqq  = c.Query("rqq")
		rqz  = c.Query("rqz")
	)
	if rqq == "" && rqz == "" {
		ldDate := models.NewLdDate("")
		rqq, rqz = ldDate.LdQrqBn, ldDate.LdCurrentDate
	}
	var cond = fmt.Sprintf(
		"qrrq>='%s 00:00:00' and qrrq<='%s 23:59:59'", rqq, rqz)
	jbzxs, err := models.GetJbzxs(cond)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	if len(jbzxs) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, jbzxs)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 金三即办注销户数明细
func JbzxGT3Mx(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		rqq  = c.Query("rqq")
		rqz  = c.Query("rqz")
	)
	if rqq == "" && rqz == "" {
		ldDate := models.NewLdDate("")
		rqq, rqz = ldDate.LdQrqBn, ldDate.LdCurrentDate
	}
	//todo
	//sql := fmt.Sprintf(
	//	"select * from jkxm_gt3 where zxsj>='%s 00:00:00' and zxsj<='%s 23:59:59'",
	//	rqq, rqz)
	pd, err := models.GetConfigSql("注销")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR,
			fmt.Sprintf("获取金三注销名单查询语句失败:%v", err))
		return
	}
	sql := fmt.Sprintf(pd.XmSql, rqq, rqz)
	gt3Zxs, err := models.QueryData(sql)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	if len(gt3Zxs) > 0 {
		var data = make([]*models.JkxmGt3, len(gt3Zxs))
		for i, gt3Zx := range gt3Zxs {
			t, _ := time.Parse(time.RFC3339, gt3Zx["ZXSJ"])
			data[i] = &models.JkxmGt3{
				Nsrsbh: gt3Zx["NSRSBH"],
				Nsrmc:  gt3Zx["NSRMC"],
				Zxry:   gt3Zx["ZXRY"],
				Zxsj:   t.Format("2006-01-02"),
			}
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 注销运行情况监控
func Zxqkjk(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		rqq  = c.Query("rqq")
		rqz  = c.Query("rqz")
	)
	if rqq == "" && rqz == "" {
		ldDate := models.NewLdDate("")
		rqq, rqz = ldDate.LdQrqBn, ldDate.LdCurrentDate
	}
	var cond = fmt.Sprintf(
		"qrrq>='%s 00:00:00' and qrrq<='%s 23:59:59'", rqq, rqz)
	jbzxs, err := models.GetJbzxs(cond)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	//todo
	//sql := fmt.Sprintf(
	//	"select * from jkxm_gt3 where zxsj>='%s 00:00:00' and zxsj<='%s 23:59:59'",
	//	rqq, rqz)
	pd, err := models.GetConfigSql("注销")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR,
			fmt.Sprintf("获取金三注销名单查询语句失败:%v", err))
		return
	}
	sql := fmt.Sprintf(pd.XmSql, rqq, rqz)
	gt3Zxs, err := models.QueryData(sql)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	var data = make([]*models.JkxmGt3, 0)
	if len(gt3Zxs) > 0 {
		for _, gt3Zx := range gt3Zxs {
			var flag = false
			for i, jbzx := range jbzxs {
				if gt3Zx["NSRSBH"] == jbzx.Nsrsbh {
					jbzxs = append(jbzxs[:i], jbzxs[i+1:]...)
					flag = true
					break
				}
			}
			if flag {
				continue
			} else {
				t, _ := time.Parse(time.RFC3339, gt3Zx["ZXSJ"])
				data = append(data, &models.JkxmGt3{
					Nsrsbh: gt3Zx["NSRSBH"],
					Nsrmc:  gt3Zx["NSRMC"],
					Zxry:   gt3Zx["ZXRY"],
					Zxsj:   t.Format("2006-01-02"),
				})
			}
		}
	}
	if len(data) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 下载确认启动即办注销流程户数明细
func DownloadJbzxJkxmMx(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		rqq  = c.Query("rqq")
		rqz  = c.Query("rqz")
	)
	if rqq == "" && rqz == "" {
		ldDate := models.NewLdDate("")
		rqq, rqz = ldDate.LdQrqBn, ldDate.LdCurrentDate
	}
	sql := fmt.Sprintf(
		"select * from jkxm_jbzx where qrrq>='%s 00:00:00' and qrrq<='%s 23:59:59'",
		rqq, rqz)
	records, err := models.QueryData(sql)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	var url = "即办注销流程名单为空"
	if len(records) > 0 {
		fileName := "即办注销流程名单" + strconv.Itoa(int(time.Now().Unix()))
		url, err = export.WriteIntoExcel(fileName, records)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, url)
}

// 下载金三即办注销流程名单
func DownloadGT3Mx(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		rqq  = c.Query("rqq")
		rqz  = c.Query("rqz")
	)
	if rqq == "" && rqz == "" {
		ldDate := models.NewLdDate("")
		rqq, rqz = ldDate.LdQrqBn, ldDate.LdCurrentDate
	}
	//todo
	//sql := fmt.Sprintf(
	//	"select * from jkxm_gt3 where zxsj>='%s 00:00:00' and zxsj<='%s 23:59:59'",
	//	rqq, rqz)
	pd, err := models.GetConfigSql("注销")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR,
			fmt.Sprintf("获取金三注销名单查询语句失败:%v", err))
		return
	}
	sql := fmt.Sprintf(pd.XmSql, rqq, rqz)
	gt3Zxs, err := models.QueryData(sql)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	var url = "金三即办注销流程名单为空"
	if len(gt3Zxs) > 0 {
		var records = make([]map[string]string, 0)
		for _, gt3Zx := range gt3Zxs {
			t, _ := time.Parse(time.RFC3339, gt3Zx["ZXSJ"])
			records = append(records, map[string]string{
				"nsrsbh": gt3Zx["NSRSBH"],
				"nsrmc":  gt3Zx["NSRMC"],
				"zxry":   gt3Zx["ZXRY"],
				"zxsj":   t.Format("2006-01-02"),
			})
		}
		fileName := "金三即办注销流程名单" + strconv.Itoa(int(time.Now().Unix()))
		url, err = export.WriteIntoExcel(fileName, records)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, url)
}

// 下载注销运行情况监控名单
func DownloadZxqkjk(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		rqq  = c.Query("rqq")
		rqz  = c.Query("rqz")
	)
	if rqq == "" && rqz == "" {
		ldDate := models.NewLdDate("")
		rqq, rqz = ldDate.LdQrqBn, ldDate.LdCurrentDate
	}
	var cond = fmt.Sprintf(
		"qrrq>='%s 00:00:00' and qrrq<='%s 23:59:59'", rqq, rqz)
	jbzxs, err := models.GetJbzxs(cond)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	//todo
	//sql := fmt.Sprintf(
	//	"select * from jkxm_gt3 where zxsj>='%s 00:00:00' and zxsj<='%s 23:59:59'",
	//	rqq, rqz)
	pd, err := models.GetConfigSql("注销")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR,
			fmt.Sprintf("获取金三注销名单查询语句失败:%v", err))
		return
	}
	sql := fmt.Sprintf(pd.XmSql, rqq, rqz)
	gt3Zxs, err := models.QueryData(sql)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	var records = make([]map[string]string, 0)
	if len(gt3Zxs) > 0 {
		for _, gt3Zx := range gt3Zxs {
			var flag = false
			for i, jbzx := range jbzxs {
				if gt3Zx["NSRSBH"] == jbzx.Nsrsbh {
					jbzxs = append(jbzxs[:i], jbzxs[i+1:]...)
					flag = true
					break
				}
			}
			if flag {
				continue
			} else {
				t, _ := time.Parse(time.RFC3339, gt3Zx["ZXSJ"])
				records = append(records, map[string]string{
					"nsrsbh": gt3Zx["NSRSBH"],
					"nsrmc":  gt3Zx["NSRMC"],
					"zxry":   gt3Zx["ZXRY"],
					"zxsj":   t.Format("2006-01-02"),
				})
			}
		}
	}
	var url = "注销运行情况监控名单为空"
	if len(records) > 0 {
		fileName := "注销运行情况监控名单" + strconv.Itoa(int(time.Now().Unix()))
		url, err = export.WriteIntoExcel(fileName, records)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, url)
}

// 手动同步后台监控指标
func SyncJkxm(c *gin.Context) {
	appG := app.Gin{C: c}
	cron.SyncJkxm()
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
