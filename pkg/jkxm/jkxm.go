package jkxm

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/util"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"strconv"
	"time"
)

// 导入监控项目数据
func ImpJkxm(fileName io.Reader, xmDm string, userInfo map[string]string) (impMsg []string) {
	xlsx, err := excelize.OpenReader(fileName)
	if err != nil {
		impMsg = append(impMsg, "导入错误,文件读取失败!请联系管理员!")
		return
	}
	sheetName := xlsx.GetSheetName(1)
	if sheetName != xmDm {
		impMsg = append(impMsg, "导入错误,该用户无此监控项目导入权限!")
		return
	}
	rows := xlsx.GetRows(sheetName)
	switch sheetName {
	case "jkxm_qs":
		erMsg := QsXmlToDB(rows, userInfo)
		impMsg = append(impMsg, erMsg...)
	case "jkxm_jcwbj":
		erMsg := JcwbjXmlToDB(rows, userInfo)
		impMsg = append(impMsg, erMsg...)
	case "jkxm_pgwbj":
		erMsg := PgwbjXmlToDB(rows, userInfo)
		impMsg = append(impMsg, erMsg...)
	case "jkxm_wjxtdhj":
		erMsg := WjxtdhjXmlToDB(rows, userInfo)
		impMsg = append(impMsg, erMsg...)
	case "jkxm_nsxydj":
		erMsg := NsxydjXmlToDB(rows, userInfo)
		impMsg = append(impMsg, erMsg...)
	case "jkxm_fxfpwcl":
		erMsg := FxfpwclXmlToDB(rows, userInfo)
		impMsg = append(impMsg, erMsg...)
	case "jkxm_fc":
		erMsg := FcXmlToDB(rows, userInfo)
		impMsg = append(impMsg, erMsg...)
	case "jkxm_td":
		erMsg := TdXmlToDB(rows, userInfo)
		impMsg = append(impMsg, erMsg...)
	case "jkxm_qt":
		erMsg := QtXmlToDB(rows, userInfo)
		impMsg = append(impMsg, erMsg...)
	}
	return
}

// 欠税
func QsXmlToDB(rows [][]string, userInfo map[string]string) (impMsg []string) {
	//遍历行读取
	for k, row := range rows {
		// 跳过标题行，遍历每行的列读取
		if k == 0 {
			continue
		}
		erMsg := make([]string, 0)
		xm := models.JkxmQs{}
		xm.ID = "QS-" + util.RandomString(17)
		xm.FqrAccount = userInfo["userAccount"]
		xm.FqrName = userInfo["username"]
		xm.FqrDepid = userInfo["departID"]
		xm.FqrDeptname = userInfo["departName"]
		xm.Fqrq = time.Now().Format("2006-01-02 15:04:05")
		for i, cell := range row {
			if cell == "" {
				switch i {
				case 1, 2:
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列存在未录入项！", k+1, i+1))
					continue
				default:
					cell = "0"
				}
			}
			switch {
			case i == 0:
				xm.Nsrsbh = cell
			case i == 1:
				xm.Nsrmc = cell
			case i == 2:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列增值税金额有误！", k+1, i+1))
				} else {
					xm.Zzs = numFormat
				}
			case i == 3:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列消费税金额有误！", k+1, i+1))
				} else {
					xm.Xfs = numFormat
				}
			case i == 4:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列企业所得税金额有误！", k+1, i+1))
				} else {
					xm.Qysds = numFormat
				}
			case i == 5:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列个人所得税金额有误！", k+1, i+1))
				} else {
					xm.Grsds = numFormat
				}
			case i == 6:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列土地增值税金额有误！", k+1, i+1))
				} else {
					xm.Tdzzs = numFormat
				}
			case i == 7:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列营业税金额有误！", k+1, i+1))
				} else {
					xm.Yys = numFormat
				}
			case i == 8:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行%d列导入错误:资源税金额有误！", k+1, i+1))
				} else {
					xm.Zys = numFormat
				}
			case i == 9:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列房产税金额有误！", k+1, i+1))
				} else {
					xm.Fcs = numFormat
				}
			case i == 10:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列印花税金额有误！", k+1, i+1))
				} else {
					xm.Yhs = numFormat
				}
			case i == 11:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列环境保护税金额有误！", k+1, i+1))
				} else {
					xm.Hjbhs = numFormat
				}
			case i == 12:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列车船税金额有误！", k+1, i+1))
				} else {
					xm.Ccs = numFormat
				}
			case i == 13:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列城市维护建设税金额有误！", k+1, i+1))
				} else {
					xm.Cswhjss = numFormat
				}
			case i == 14:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列城镇土地使用税金额有误！", k+1, i+1))
				} else {
					xm.Cztdsys = numFormat
				}
			case i == 15:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列耕地占用税金额有误！", k+1, i+1))
				} else {
					xm.Gdzys = numFormat
				}
			case i == 16:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列契税金额有误！", k+1, i+1))
				} else {
					xm.Qs = numFormat
				}
			case i == 18:
				numFormat, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					erMsg = append(erMsg, fmt.Sprintf(
						"第%d行导入错误:第%d列其他税收收入金额有误！", k+1, i+1))
				} else {
					xm.Qtsssr = numFormat
				}
			}
		}
		//logging.Debug(fmt.Sprintf("欠税: %+v", &xm))
		if len(erMsg) > 0 {
			impMsg = append(impMsg, erMsg...)
			continue
		}
		err := models.AddJkxmData("jkxm_qs", &xm)
		if err != nil {
			impMsg = append(impMsg, fmt.Sprintf("第%d行导入错误:%s", k+1, err.Error()))
		}
	}
	if len(impMsg) > 0 {
		impMsg = append(impMsg, "除上述记录外导入成功!")
	} else {
		impMsg = append(impMsg, fmt.Sprintf("导入成功,共导入%d条!", len(rows)-1))
	}
	return
}

// 稽查未办结
func JcwbjXmlToDB(rows [][]string, userInfo map[string]string) (impMsg []string) {
	//遍历行读取
	for k, row := range rows {
		// 跳过标题行，遍历每行的列读取
		if k == 0 {
			continue
		}
		erMsg := make([]string, 0)
		xm := models.JkxmJcwbj{}
		xm.ID = "JCWBJ-" + util.RandomString(14)
		xm.FqrAccount = userInfo["userAccount"]
		xm.FqrName = userInfo["username"]
		xm.FqrDepid = userInfo["departID"]
		xm.FqrDeptname = userInfo["departName"]
		xm.Fqrq = time.Now().Format("2006-01-02 15:04:05")
		for i, cell := range row {
			if cell == "" {
				erMsg = append(erMsg, fmt.Sprintf(
					"第%d行导入错误:第%d列存在未录入项！", k+1, i+1))
				continue
			}
			switch {
			case i == 0:
				xm.Nsrsbh = cell
			case i == 1:
				xm.Nsrmc = cell
			case i == 2:
				xm.Jcajbh = cell
			case i == 3:
				xm.Jcyr = cell
			}
		}
		//logging.Debug(fmt.Sprintf("稽查未办结: %+v", &xm))
		if len(erMsg) > 0 {
			impMsg = append(impMsg, erMsg...)
			continue
		}
		err := models.AddJkxmData("jkxm_jcwbj", &xm)
		if err != nil {
			impMsg = append(impMsg, fmt.Sprintf("第%d行导入错误:%s", k+1, err.Error()))
		}
	}
	if len(impMsg) > 0 {
		impMsg = append(impMsg, "除上述记录外导入成功!")
	} else {
		impMsg = append(impMsg, fmt.Sprintf("导入成功,共导入%d条!", len(rows)-1))
	}
	return
}

// 评估未办结
func PgwbjXmlToDB(rows [][]string, userInfo map[string]string) (impMsg []string) {
	//遍历行读取
	for k, row := range rows {
		// 跳过标题行，遍历每行的列读取
		if k == 0 {
			continue
		}
		erMsg := make([]string, 0)
		xm := models.JkxmPgwbj{}
		xm.ID = "PGWBJ-" + util.RandomString(14)
		xm.FqrAccount = userInfo["userAccount"]
		xm.FqrName = userInfo["username"]
		xm.FqrDepid = userInfo["departID"]
		xm.FqrDeptname = userInfo["departName"]
		xm.Fqrq = time.Now().Format("2006-01-02 15:04:05")
		for i, cell := range row {
			if cell == "" {
				erMsg = append(erMsg, fmt.Sprintf(
					"第%d行导入错误:第%d列存在未录入项！", k+1, i+1))
				continue
			}
			switch {
			case i == 0:
				xm.Nsrsbh = cell
			case i == 1:
				xm.Nsrmc = cell
			case i == 2:
				xm.Pgajbh = cell
			case i == 3:
				xm.Pgry = cell
			}
		}
		//logging.Debug(fmt.Sprintf("评估未办结: %+v", &xm))
		if len(erMsg) > 0 {
			impMsg = append(impMsg, erMsg...)
			continue
		}
		err := models.AddJkxmData("jkxm_pgwbj", &xm)
		if err != nil {
			impMsg = append(impMsg, fmt.Sprintf("第%d行导入错误:%s", k+1, err.Error()))
		}
	}
	if len(impMsg) > 0 {
		impMsg = append(impMsg, "除上述记录外导入成功!")
	} else {
		impMsg = append(impMsg, fmt.Sprintf("导入成功,共导入%d条!", len(rows)-1))
	}
	return
}

// 未进行土增汇缴
func WjxtdhjXmlToDB(rows [][]string, userInfo map[string]string) (impMsg []string) {
	//遍历行读取
	for k, row := range rows {
		// 跳过标题行，遍历每行的列读取
		if k == 0 {
			continue
		}
		erMsg := make([]string, 0)
		xm := models.JkxmWjxtdhj{}
		xm.ID = "WJXTZHJ-" + util.RandomString(12)
		xm.FqrAccount = userInfo["userAccount"]
		xm.FqrName = userInfo["username"]
		xm.FqrDepid = userInfo["departID"]
		xm.FqrDeptname = userInfo["departName"]
		xm.Fqrq = time.Now().Format("2006-01-02 15:04:05")
		for i, cell := range row {
			if cell == "" {
				erMsg = append(erMsg, fmt.Sprintf(
					"第%d行导入错误:第%d列存在未录入项！", k+1, i+1))
				continue
			}
			switch {
			case i == 0:
				xm.Nsrsbh = cell
			case i == 1:
				xm.Nsrmc = cell
			case i == 2:
				xm.Sbhjbz = cell
			case i == 3:
				xm.Xmmc = cell
			}
		}
		//logging.Debug(fmt.Sprintf("未进行土增汇缴: %+v", &xm))
		if len(erMsg) > 0 {
			impMsg = append(impMsg, erMsg...)
			continue
		}
		err := models.AddJkxmData("jkxm_wjxtdhj", &xm)
		if err != nil {
			impMsg = append(impMsg, fmt.Sprintf("第%d行导入错误:%s", k+1, err.Error()))
		}
	}
	if len(impMsg) > 0 {
		impMsg = append(impMsg, "除上述记录外导入成功!")
	} else {
		impMsg = append(impMsg, fmt.Sprintf("导入成功,共导入%d条!", len(rows)-1))
	}
	return
}

// 纳税信用等级
func NsxydjXmlToDB(rows [][]string, userInfo map[string]string) (impMsg []string) {
	//遍历行读取
	for k, row := range rows {
		// 跳过标题行，遍历每行的列读取
		if k == 0 {
			continue
		}
		erMsg := make([]string, 0)
		xm := models.JkxmNsxydj{}
		xm.ID = "NSXYDJ-" + util.RandomString(13)
		xm.FqrAccount = userInfo["userAccount"]
		xm.FqrName = userInfo["username"]
		xm.FqrDepid = userInfo["departID"]
		xm.FqrDeptname = userInfo["departName"]
		xm.Fqrq = time.Now().Format("2006-01-02 15:04:05")
		for i, cell := range row {
			if cell == "" {
				erMsg = append(erMsg, fmt.Sprintf(
					"第%d行导入错误:第%d列存在未录入项！", k+1, i+1))
				continue
			}
			switch {
			case i == 0:
				xm.Nsrsbh = cell
			case i == 1:
				xm.Nsrmc = cell
			case i == 2:
				xm.Nsxydj = cell
			}
		}
		//logging.Debug(fmt.Sprintf("纳税信用等级: %+v", &xm))
		if len(erMsg) > 0 {
			impMsg = append(impMsg, erMsg...)
			continue
		}
		err := models.AddJkxmData("jkxm_nsxydj", &xm)
		if err != nil {
			impMsg = append(impMsg, fmt.Sprintf("第%d行导入错误:%s", k+1, err.Error()))
		}
	}
	if len(impMsg) > 0 {
		impMsg = append(impMsg, "除上述记录外导入成功!")
	} else {
		impMsg = append(impMsg, fmt.Sprintf("导入成功,共导入%d条!", len(rows)-1))
	}
	return
}

// 风险发票未处理
func FxfpwclXmlToDB(rows [][]string, userInfo map[string]string) (impMsg []string) {
	//遍历行读取
	for k, row := range rows {
		// 跳过标题行，遍历每行的列读取
		if k == 0 {
			continue
		}
		erMsg := make([]string, 0)
		xm := models.JkxmFxfpwcl{}
		xm.ID = "FXFPWCL-" + util.RandomString(12)
		xm.FqrAccount = userInfo["userAccount"]
		xm.FqrName = userInfo["username"]
		xm.FqrDepid = userInfo["departID"]
		xm.FqrDeptname = userInfo["departName"]
		xm.Fqrq = time.Now().Format("2006-01-02 15:04:05")
		for i, cell := range row {
			if cell == "" {
				impMsg = append(impMsg, fmt.Sprintf(
					"第%d行导入错误:第%d列存在未录入项！", k+1, i+1))
				continue
			}
			switch {
			case i == 0:
				xm.Nsrsbh = cell
			case i == 1:
				xm.Nsrmc = cell
			case i == 2:
				xm.Fpdm = cell
			case i == 3:
				xm.Fphm = cell
			case i == 4:
				xm.Hsry = cell
			}
		}
		//logging.Debug(fmt.Sprintf("风险发票未处理: %+v", &xm))
		if len(erMsg) > 0 {
			impMsg = append(impMsg, erMsg...)
			continue
		}
		err := models.AddJkxmData("jkxm_fxfpwcl", &xm)
		if err != nil {
			impMsg = append(impMsg, fmt.Sprintf("第%d行导入错误:%s", k+1, err.Error()))
		}
	}
	if len(impMsg) > 0 {
		impMsg = append(impMsg, "除上述记录外导入成功!")
	} else {
		impMsg = append(impMsg, fmt.Sprintf("导入成功,共导入%d条!", len(rows)-1))
	}
	return
}

// 房产
func FcXmlToDB(rows [][]string, userInfo map[string]string) (impMsg []string) {
	//遍历行读取
	for k, row := range rows {
		// 跳过标题行，遍历每行的列读取
		if k == 0 {
			continue
		}
		erMsg := make([]string, 0)
		xm := models.JkxmFc{}
		xm.ID = "FC-" + util.RandomString(17)
		xm.FqrAccount = userInfo["userAccount"]
		xm.FqrName = userInfo["username"]
		xm.FqrDepid = userInfo["departID"]
		xm.FqrDeptname = userInfo["departName"]
		xm.Fqrq = time.Now().Format("2006-01-02 15:04:05")
		for i, cell := range row {
			if cell == "" {
				erMsg = append(erMsg, fmt.Sprintf(
					"第%d行导入错误:第%d列存在未录入项！", k+1, i+1))
				continue
			}
			switch {
			case i == 0:
				xm.Nsrsbh = cell
			case i == 1:
				xm.Nsrmc = cell
			case i == 2:
				xm.Fcdz = cell
			case i == 3:
				xm.Fcbh = cell
			}
		}
		//logging.Debug(fmt.Sprintf("房产: %+v", &xm))
		if len(erMsg) > 0 {
			impMsg = append(impMsg, erMsg...)
			continue
		}
		err := models.AddJkxmData("jkxm_fc", &xm)
		if err != nil {
			impMsg = append(impMsg, fmt.Sprintf("第%d行导入错误:%s", k+1, err.Error()))
		}
	}
	if len(impMsg) > 0 {
		impMsg = append(impMsg, "除上述记录外导入成功!")
	} else {
		impMsg = append(impMsg, fmt.Sprintf("导入成功,共导入%d条!", len(rows)-1))
	}
	return
}

// 土地
func TdXmlToDB(rows [][]string, userInfo map[string]string) (impMsg []string) {
	//遍历行读取
	for k, row := range rows {
		// 跳过标题行，遍历每行的列读取
		if k == 0 {
			continue
		}
		erMsg := make([]string, 0)
		xm := models.JkxmTd{}
		xm.ID = "TD-" + util.RandomString(17)
		xm.FqrAccount = userInfo["userAccount"]
		xm.FqrName = userInfo["username"]
		xm.FqrDepid = userInfo["departID"]
		xm.FqrDeptname = userInfo["departName"]
		xm.Fqrq = time.Now().Format("2006-01-02 15:04:05")
		for i, cell := range row {
			if cell == "" {
				erMsg = append(erMsg, fmt.Sprintf(
					"第%d行导入错误:第%d列存在未录入项！", k+1, i+1))
				continue
			}
			switch {
			case i == 0:
				xm.Nsrsbh = cell
			case i == 1:
				xm.Nsrmc = cell
			case i == 2:
				xm.Tddz = cell
			case i == 3:
				xm.Tdbh = cell
			}
		}
		//logging.Debug(fmt.Sprintf("土地: %+v", &xm))
		if len(erMsg) > 0 {
			impMsg = append(impMsg, erMsg...)
			continue
		}
		err := models.AddJkxmData("jkxm_td", &xm)
		if err != nil {
			impMsg = append(impMsg, fmt.Sprintf("第%d行导入错误:%s", k+1, err.Error()))
		}
	}
	if len(impMsg) > 0 {
		impMsg = append(impMsg, "除上述记录外导入成功!")
	} else {
		impMsg = append(impMsg, fmt.Sprintf("导入成功,共导入%d条!", len(rows)-1))
	}
	return
}

// 其他
func QtXmlToDB(rows [][]string, userInfo map[string]string) (impMsg []string) {
	//遍历行读取
	for k, row := range rows {
		// 跳过标题行，遍历每行的列读取
		if k == 0 {
			continue
		}
		erMsg := make([]string, 0)
		xm := models.JkxmQt{}
		xm.ID = "QT-" + util.RandomString(17)
		xm.FqrAccount = userInfo["userAccount"]
		xm.FqrName = userInfo["username"]
		xm.FqrDepid = userInfo["departID"]
		xm.FqrDeptname = userInfo["departName"]
		xm.Fqrq = time.Now().Format("2006-01-02 15:04:05")
		for i, cell := range row {
			if cell == "" {
				erMsg = append(erMsg, fmt.Sprintf(
					"第%d行导入错误:第%d列存在未录入项！", k+1, i+1))
				continue
			}
			switch {
			case i == 0:
				xm.Nsrsbh = cell
			case i == 1:
				xm.Nsrmc = cell
			case i == 2:
				xm.Qtxzxx = cell
			}
		}
		//logging.Debug(fmt.Sprintf("其他: %+v", &xm))
		if len(erMsg) > 0 {
			impMsg = append(impMsg, erMsg...)
			continue
		}
		err := models.AddJkxmData("jkxm_qt", &xm)
		if err != nil {
			impMsg = append(impMsg, fmt.Sprintf("第%d行导入错误:%s", k+1, err.Error()))
		}
	}
	if len(impMsg) > 0 {
		impMsg = append(impMsg, "除上述记录外导入成功!")
	} else {
		impMsg = append(impMsg, fmt.Sprintf("导入成功,共导入%d条!", len(rows)-1))
	}
	return
}
