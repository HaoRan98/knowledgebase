package jkxm

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/logging"
	"NULL/knowledgebase/pkg/util"
	"fmt"
	"strings"
	"time"
)

//----------------------------定时同步开发区数据中台指标----------------------------
//欠税
func SyncJkxmQs(ldDate *models.LdDate) {
	var add, upd, end int
	pd, err := models.GetConfigSql("欠税")
	if err != nil {
		logging.Error(fmt.Sprintf(
			"同步欠税数据指标异常!获取查询语句失败:%v", err))
		return
	}
	sql := fmt.Sprintf(pd.XmSql, ldDate.LdCurrentDate, ldDate.LdCurrentDate)
	records, err := models.QueryData(sql)
	if err != nil {
		logging.Error(fmt.Sprintf(
			"同步欠税指标异常!获取金三数据失败:%v", err))
		return
	}
	if len(records) > 0 {
		logging.Info(fmt.Sprintf("需同步欠税数据%d条!", len(records)))
		//新增&更新欠税
		for _, record := range records {
			xm := models.JkxmQs{}
			xm.ID = "QS-" + util.RandomString(17)
			t := time.Now().Format("2006-01-02 15:04:05")
			xm.FqrName = "后台同步"
			xm.Fqrq = t
			xm.ShrName = "后台同步审核"
			xm.Shrq = t
			xm.Shbz = "Y"
			xm.Nsrsbh = record["NSRSBH"]
			xm.Nsrmc = record["NSRMC"]
			updXm := map[string]string{
				"zzs":      "0.00",
				"xfs":      "0.00",
				"qysds":    "0.00",
				"grsds":    "0.00",
				"tdzzs":    "0.00",
				"yys":      "0.00",
				"zys":      "0.00",
				"fcs":      "0.00",
				"yhs":      "0.00",
				"hjbhs":    "0.00",
				"ccs":      "0.00",
				"cswhjss":  "0.00",
				"cztdsys":  "0.00",
				"gdzys":    "0.00",
				"qs":       "0.00",
				"qtsssr":   "0.00",
				"clgzs":    "0.00",
				"whsyjsf":  "0.00",
				"sljszxsr": "0.00",
				"cjrjybzj": "0.00",
				"dfjyfj":   "0.00",
				"swbmfmsr": "0.00",
				"jyffj":    "0.00",
				"qtsr":     "0.00",
			}
			szQs := strings.Split(record["QS"], ",")
			for _, qs := range szQs {
				if strings.Contains(qs, ":") {
					sz := strings.Split(qs, ":")[0]
					je := strings.Split(qs, ":")[1]
					switch sz {
					case "增值税":
						xm.Zzs = je
						updXm["zzs"] = je
					case "消费税":
						xm.Xfs = je
						updXm["xfs"] = je
					case "企业所得税":
						xm.Qysds = je
						updXm["qysds"] = je
					case "个人所得税":
						xm.Grsds = je
						updXm["grsds"] = je
					case "土地增值税":
						xm.Tdzzs = je
						updXm["tdzzs"] = je
					case "营业税":
						xm.Yys = je
						updXm["yys"] = je
					case "资源税":
						xm.Zys = je
						updXm["zys"] = je
					case "房产税":
						xm.Fcs = je
						updXm["fcs"] = je
					case "印花税":
						xm.Yhs = je
						updXm["yhs"] = je
					case "环境保护税":
						xm.Hjbhs = je
						updXm["hjbhs"] = je
					case "车船税":
						xm.Ccs = je
						updXm["ccs"] = je
					case "城市维护建设税":
						xm.Cswhjss = je
						updXm["cswhjss"] = je
					case "城镇土地使用税":
						xm.Cztdsys = je
						updXm["cztdsys"] = je
					case "耕地占用税":
						xm.Gdzys = je
						updXm["gdzys"] = je
					case "契税":
						xm.Qs = je
						updXm["qs"] = je
					case "其他税收收入":
						xm.Qtsssr = je
						updXm["qtsssr"] = je
					case "车辆购置税":
						xm.Clgzs = je
						updXm["clgzs"] = je
					case "文化事业建设费":
						xm.Whsyjsf = je
						updXm["whsyjsf"] = je
					case "水利建设专项收入":
						xm.Sljszxsr = je
						updXm["sljszxsr"] = je
					case "残疾人就业保障金":
						xm.Cjrjybzj = je
						updXm["cjrjybzj"] = je
					case "地方教育附加":
						xm.Dfjyfj = je
						updXm["dfjyfj"] = je
					case "税务部门罚没收入":
						xm.Swbmfmsr = je
						updXm["swbmfmsr"] = je
					case "教育费附加":
						xm.Jyffj = je
						updXm["jyffj"] = je
					case "其他收入":
						xm.Qtsr = je
						updXm["qtsr"] = je
					}
				}
			}
			query := fmt.Sprintf(
				"select * from jkxm_qs where nsrsbh='%s' and zjbz='N'",
				record["NSRSBH"])
			r, err := models.QueryData(query)
			if err != nil {
				logging.Error(fmt.Sprintf(
					"获取源表欠税指标异常!失败记录\n:%v", record))
				continue
			}
			if len(r) == 0 { //add record
				if err := models.AddJkxmData("jkxm_qs", &xm); err != nil {
					logging.Error(fmt.Sprintf(
						"添加欠税数据指标异常!失败记录\n:%v", record))
					continue
				}
				add++
			} else { //update record
				updXm["nsrsbh"] = r[0]["nsrsbh"]
				updXm["nsrmc"] = r[0]["nsrmc"]

				updXm["fqr_account"] = r[0]["fqr_account"]
				updXm["fqr_name"] = r[0]["fqr_name"]
				updXm["fqr_depid"] = r[0]["fqr_depid"]
				updXm["fqr_deptname"] = r[0]["fqr_deptname"]

				updXm["shrq"] = r[0]["shrq"]
				updXm["shr_account"] = r[0]["shr_account"]
				updXm["shr_name"] = r[0]["shr_name"]
				updXm["shr_depid"] = r[0]["shr_depid"]
				updXm["shr_deptname"] = r[0]["shr_deptname"]
				updXm["shrq"] = r[0]["shrq"]
				updXm["shbz"] = r[0]["shbz"]
				if err := models.UpdateJkxmQs(r[0]["id"], updXm); err != nil {
					logging.Error(fmt.Sprintf(
						"修改欠税数据指标异常!\n源id:%v\n失败记录\n:%v",
						r[0]["id"], record))
					continue
				}
				upd++
			}
		}
		//清税终结
		query := "select * from jkxm_qs where zjbz='N'"
		rs, err := models.QueryData(query)
		if err != nil {
			logging.Info(fmt.Sprintf(
				"欠税数据指标同步成功!新增%d条，修改%d条！", add, upd))
			logging.Error("清税终结失败！获取源表欠税指标异常!")
			return
		}
		for _, r := range rs {
			var flag = false
			for i, record := range records {
				if r["nsrsbh"] == record["NSRSBH"] {
					records = append(records[:i], records[i+1:]...)
					flag = true
					break
				}
			}
			if flag {
				continue
			} else {
				t := time.Now().Format("2006-01-02 15:04:05")
				r["zjr_name"] = "后台清税终结"
				r["zjrq"] = t
				r["zjbz"] = "Y"

				r["zjshr_name"] = "后台清税终审"
				r["zjshrq"] = t
				r["zjshbz"] = "Y"
				if err := models.UpdateJkxmQs(r["id"], r); err != nil {
					logging.Error(fmt.Sprintf(
						"清税终结失败!\n源id:%v", r["id"]))
					continue
				}
				end++
			}
		}
		logging.Info(fmt.Sprintf("新增%d条,修改%d条,终结%d条!", add, upd, end))
	}
	logging.Info("欠税数据指标同步成功!清税终结同步成功!")
}

//出口退（免）税备案
func SyncJkxmCktsba() {
	pd, err := models.GetConfigSql("出口退（免）税备案")
	if err != nil {
		logging.Error(fmt.Sprintf(
			"同步出口退（免）税备案数据指标异常!获取查询语句失败:%v", err))
		return
	}
	records, err := models.QueryData(pd.XmSql)
	if err != nil {
		logging.Error(fmt.Sprintf(
			"同步出口退（免）税备案数据指标异常!获取金三数据失败:%v", err))
		return
	}
	var total = len(records)
	if len(records) > 0 {
		for _, record := range records {
			xm := models.JkxmCktsba{}
			xm.ID = "CKTSBA-" + util.RandomString(13)
			t := time.Now().Format("2006-01-02 15:04:05")
			xm.FqrName = "后台同步"
			xm.Fqrq = t
			xm.ShrName = "后台同步审核"
			xm.Shrq = t
			xm.Shbz = "Y"
			xm.Nsrsbh = record["NSRSBH"]
			xm.Nsrmc = record["NSRMC"]
			query := fmt.Sprintf(
				"select * from jkxm_cktsba where nsrsbh='%s'", record["NSRSBH"])
			r, err := models.QueryData(query)
			if err != nil {
				logging.Error(fmt.Sprintf(
					"同步出口退（免）税备案数据指标异常!失败记录\n:%v", record))
				total--
				continue
			}
			if len(r) > 0 {
				total--
				continue
			}
			err = models.AddJkxmData("jkxm_cktsba", &xm)
			if err != nil {
				logging.Error(fmt.Sprintf(
					"同步出口退（免）税备案数据指标异常!失败记录\n:%v", record))
				total--
				continue
			}
			//备案资格终结
			squery := "select * from jkxm_cktsba where zjbz='N'"
			rs, err := models.QueryData(squery)
			if err != nil {
				logging.Error("备案资格终结失败！获取源表出口备案指标异常!")
				return
			}
			for _, r := range rs {
				var flag = false
				for i, record := range records {
					if r["nsrsbh"] == record["NSRSBH"] {
						records = append(records[:i], records[i+1:]...)
						flag = true
						break
					}
				}
				if flag {
					continue
				} else {
					t := time.Now().Format("2006-01-02 15:04:05")
					r["zjr_name"] = "后台终结"
					r["zjrq"] = t
					r["zjbz"] = "Y"

					r["zjshr_name"] = "后台终审"
					r["zjshrq"] = t
					r["zjshbz"] = "Y"
					if err := models.UpdateJkxmQs(r["id"], r); err != nil {
						logging.Error(fmt.Sprintf(
							"备案资格终结失败!\n源id:%v", r["id"]))
						continue
					}
				}
			}
		}
	}
	logging.Info(fmt.Sprintf("出口退（免）税备案数据指标同步成功,共同步%d条数据!", total))
}

//评估未办结
func SyncJkxmPgwbj(ldDate *models.LdDate) {
	pd, err := models.GetConfigSql("评估未办结")
	if err != nil {
		logging.Error(fmt.Sprintf(
			"同步评估未办结数据指标异常!获取查询语句失败:%v", err))
		return
	}
	//todo
	sql := fmt.Sprintf(pd.XmSql)
	records, err := models.QueryData(sql)
	if err != nil {
		logging.Error(fmt.Sprintf(
			"同步评估未办结数据指标异常!获取金三数据失败:%v", err))
		return
	}
	var total = len(records)
	if len(records) > 0 {
		for _, record := range records {
			xm := models.JkxmPgwbj{}
			xm.ID = "PGWBJ-" + util.RandomString(14)
			t := time.Now().Format("2006-01-02 15:04:05")
			xm.FqrName = "后台同步"
			xm.Fqrq = t
			xm.ShrName = "后台同步审核"
			xm.Shrq = t
			xm.Nsrsbh = record["NSRSBH"]
			xm.Nsrmc = record["NSRMC"]
			xm.Pgajbh = record["PGAJBH"]
			xm.Pgry = record["PGRY"]
			err := models.AddJkxmData("jkxm_pgwbj", &xm)
			if err != nil {
				logging.Error(fmt.Sprintf(
					"同步评估未办结数据指标异常!失败记录\n:%v", record))
				total--
				continue
			}
		}
	}
	logging.Info(fmt.Sprintf("评估未办结数据指标同步成功,共同步%d条数据!", total))
}

//纳税信用等级
func SyncJkxmNsxydj(ldDate *models.LdDate) {
	pd, err := models.GetConfigSql("纳税信用等级")
	if err != nil {
		logging.Error(fmt.Sprintf(
			"同步纳税信用等级数据指标异常!获取查询语句失败:%v", err))
		return
	}
	//todo
	sql := fmt.Sprintf(pd.XmSql)
	records, err := models.QueryData(sql)
	if err != nil {
		logging.Error(fmt.Sprintf(
			"同步纳税信用等级数据指标异常!获取金三数据失败:%v", err))
		return
	}
	var total = len(records)
	if len(records) > 0 {
		for _, record := range records {
			xm := models.JkxmNsxydj{}
			xm.ID = "NSXYDJ-" + util.RandomString(14)
			t := time.Now().Format("2006-01-02 15:04:05")
			xm.FqrName = "后台同步"
			xm.Fqrq = t
			xm.ShrName = "后台同步审核"
			xm.Shrq = t
			xm.Nsrsbh = record["NSRSBH"]
			xm.Nsrmc = record["NSRMC"]
			err = models.AddJkxmData("jkxm_nsxydj", &xm)
			if err != nil {
				logging.Error(fmt.Sprintf(
					"同步纳税信用等级数据指标异常!失败记录\n:%v", record))
				total--
				continue
			}
		}
	}
	logging.Info(fmt.Sprintf("纳税信用等级数据指标同步成功,共同步%d条数据!", total))
}
