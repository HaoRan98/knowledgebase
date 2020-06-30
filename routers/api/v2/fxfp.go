package v2

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/app"
	"NULL/knowledgebase/pkg/e"
	"NULL/knowledgebase/pkg/export"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//风险发票超XX份或税额超XX万元
func GetJkxmFxfp(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		se   = c.Query("se")
		cnt  = c.Query("cnt")
		nsrs []*models.NsrInfo
	)
	n1, err := models.GetJkxmFxfpOverSe(se)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if n1 != nil {
		nsrs = append(nsrs, n1...)
	}
	n2, err := models.GetJkxmFxfpOverNum(cnt)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if n2 != nil {
		nsrs = append(nsrs, n2...)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nsrs)
}

//获取风险发票明细
func GetJkxmFxfpByNsrsbh(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		se     = c.Query("se")
		cnt    = c.Query("cnt")
		nsrsbh = c.Query("nsrsbh") //不传取全部，传取某一户
		nsrs   []*models.NsrInfo
	)
	if len(nsrsbh) > 0 {
		fxfpwcls, err := models.GetJkxmFxfpByNsrsbh(nsrsbh)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
		if len(fxfpwcls) > 0 {
			appG.Response(http.StatusOK, e.SUCCESS, fxfpwcls)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
	n1, err := models.GetJkxmFxfpOverSe(se)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if n1 != nil {
		nsrs = append(nsrs, n1...)
	}
	n2, err := models.GetJkxmFxfpOverNum(cnt)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if n2 != nil {
		nsrs = append(nsrs, n2...)
	}
	data := make([]*models.JkxmFxfpwcl, 0)
	if len(nsrs) > 0 {
		for _, nsr := range nsrs {
			fxfpwcls, err := models.GetJkxmFxfpByNsrsbh(nsr.Nsrsbh)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR, err)
				return
			}
			if len(fxfpwcls) > 0 {
				data = append(data, fxfpwcls...)
			}
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//下载风险发票明细
func DownloadJkxmFxfp(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		err     error
		se      = c.Query("se")
		cnt     = c.Query("cnt")
		nsrsbh  = c.Query("nsrsbh")
		nsrs    []*models.NsrInfo
		records = make([]map[string]string, 0)
	)
	var url = "该监控项目没有异常数据"
	if len(nsrsbh) > 0 {
		squery := fmt.Sprintf(
			`select * from jkxm_fxfpwcl where nsrsbh='%s'`, nsrsbh)
		record, err := models.QueryData(squery)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
		if len(record) > 0 {
			records = append(records, record...)
		}
	} else {
		n1, err := models.GetJkxmFxfpOverSe(se)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
		if n1 != nil {
			nsrs = append(nsrs, n1...)
		}
		n2, err := models.GetJkxmFxfpOverNum(cnt)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
		if n2 != nil {
			nsrs = append(nsrs, n2...)
		}
		if len(nsrs) > 0 {
			for _, nsr := range nsrs {
				squery := fmt.Sprintf(
					`select * from jkxm_fxfpwcl where nsrsbh='%s'`, nsr.Nsrsbh)
				record, err := models.QueryData(squery)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR, err)
					return
				}
				if len(record) > 0 {
					records = append(records, record...)
				}
			}
		}
	}
	if len(records) > 0 {
		fileName := "风险发票明细" + strconv.Itoa(int(time.Now().Unix()))
		url, err = export.WriteIntoExcel(fileName, records)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, url)
}
