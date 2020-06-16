package models

import "strings"

// 监控项目--基础表
type JkxmBase struct {
	Nsrsbh string `json:"nsrsbh" gorm:"COMMENT:'纳税人识别号'"`
	Nsrmc  string `json:"nsrmc" gorm:"COMMENT:'纳税人名称'"`

	FqrAccount  string `json:"fqr_account" gorm:"COMMENT:'发起人账号'"`
	FqrName     string `json:"fqr_name" gorm:"COMMENT:'发起人姓名'"`
	FqrDepid    string `json:"fqr_depid" gorm:"COMMENT:'发起人部门id'"`
	FqrDeptname string `json:"fqr_deptname" gorm:"COMMENT:'发起人部门名称'"`
	Fqrq        string `json:"fqrq" gorm:"COMMENT:'发起日期'"`

	ShrAccount  string `json:"shr_account" gorm:"COMMENT:'审核人账号'"`
	ShrName     string `json:"shr_name" gorm:"COMMENT:'审核人姓名'"`
	ShrDepid    string `json:"shr_depid" gorm:"COMMENT:'审核人部门id'"`
	ShrDeptname string `json:"shr_deptname" gorm:"COMMENT:'审核人部门名称'"`
	Shrq        string `json:"shrq" gorm:"COMMENT:'审核日期'"`
	Shbz        string `json:"shbz" gorm:"COMMENT:'审核标志';default:'N'"`

	ZjrAccount  string `json:"zjr_account" gorm:"COMMENT:'终结人账号'"`
	ZjrName     string `json:"zjr_name" gorm:"COMMENT:'终结人姓名'"`
	ZjrDepid    string `json:"zjr_depid" gorm:"COMMENT:'终结人部门id'"`
	ZjrDeptname string `json:"zjr_deptname" gorm:"COMMENT:'终结人部门名称'"`
	Zjrq        string `json:"zjrq" gorm:"COMMENT:'终结日期'"`
	Zjbz        string `json:"zjbz" gorm:"COMMENT:'终结标志';default:'N'"`

	ZjshrAccount  string `json:"zjshr_account" gorm:"COMMENT:'终结审核人账号'"`
	ZjshrName     string `json:"zjshr_name" gorm:"COMMENT:'终结审核人姓名'"`
	ZjshrDepid    string `json:"zjshr_depid" gorm:"COMMENT:'终结审核人部门id'"`
	ZjshrDeptname string `json:"zjshr_deptname" gorm:"COMMENT:'终结审核人部门名称'"`
	Zjshrq        string `json:"zjshrq" gorm:"COMMENT:'终结审核日期'"`
	Zjshbz        string `json:"zjshbz" gorm:"COMMENT:'终结审核标志';default:'N'"`
}

func AddJkxmData(tName string, data interface{}) error {
	if err := db.Table(tName).Create(data).Error; err != nil {
		return err
	}
	return nil
}

func ShJkxm(tName, id string, shMap map[string]string) error {
	if err := db.Table(tName).Where("id=?", id).Updates(shMap).Error; err != nil {
		return err
	}
	return nil
}

func ZjJkxm(tName, id string, zjMap map[string]string) error {
	if err := db.Table(tName).Where("id=?", id).Updates(zjMap).Error; err != nil {
		return err
	}
	return nil
}

func CountJkxmsToal(tName, cond string) int {
	var cnt int
	if err := db.Table(tName).
		Where(cond).Where("shbz='Y'").Count(&cnt).Error; err != nil {
		return 0
	}
	return cnt
}

func CountJkxmRsolved(tName, cond string) int {
	var cnt int
	if err := db.Table(tName).
		Where(cond).Where("zjshbz='Y'").Count(&cnt).Error; err != nil {
		return 0
	}
	return cnt
}

func CountJkxmsUnsolved(tName, cond string) int {
	var cnt int
	if err := db.Table(tName).
		Where(cond).Where("shbz='Y' and zjshbz='N'").Count(&cnt).Error; err != nil {
		return 0
	}
	return cnt
}

func QueryData(squery string) ([]map[string]string, error) {
	rows, err := db.Raw(squery).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range cols {
		scans[i] = &vals[i]
	}
	var results []map[string]string
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return results, err
		}
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}
	return results, nil
}

func SortJkxmFields(tHeader map[string]string) []string {
	sorted_keys, tempMap := make([]string, 12), make([]string, 10)
	for filed := range tHeader {
		filed = strings.ToLower(filed)
		if filed == "id" ||
			strings.Contains(filed, "_account") ||
			strings.Contains(filed, "_depid") ||
			strings.Contains(filed, "shbz") ||
			strings.Contains(filed, "zjbz") ||
			strings.Contains(filed, "zjshr_name") ||
			strings.Contains(filed, "zjshr_deptname") ||
			strings.Contains(filed, "zjshrq") {
			continue
		}
		switch filed {
		case "nsrsbh":
			sorted_keys[0] = filed
		case "nsrmc":
			sorted_keys[1] = filed
		case "fqr_name":
			sorted_keys[2] = filed
		case "fqr_deptname":
			sorted_keys[3] = filed
		case "fqrq":
			sorted_keys[4] = filed

		case "shr_name":
			sorted_keys[5] = filed
		case "shr_deptname":
			sorted_keys[6] = filed
		case "shrq":
			sorted_keys[7] = filed

		case "zjr_name":
			sorted_keys[8] = filed
		case "zjr_deptname":
			sorted_keys[9] = filed
		case "zjrq":
			sorted_keys[10] = filed

		case "zjshbz":
			sorted_keys[11] = filed

		default:
			tempMap = append(tempMap, filed)
		}
	}
	sorted_keys = append(sorted_keys, tempMap...)
	return sorted_keys
}

func ReplaceTableFileds(field string) string {
	switch field {
	case "nsrsbh":
		field = "纳税人识别号"
	case "nsrmc":
		field = "纳税人名称"

	case "fqr_name":
		field = "发起人"
	case "fqr_deptname":
		field = "发起人部门"
	case "fqrq":
		field = "发起日期"

	case "shr_name":
		field = "审核人"
	case "shr_deptname":
		field = "审核人部门"
	case "shrq":
		field = "审核日期"
	case "shbz":
		field = "审核标志"

	case "zjr_name":
		field = "终结人"
	case "zjr_deptname":
		field = "终结人部门"
	case "zjrq":
		field = "终结日期"

	case "zjshbz":
		field = "终审标志"

	case "zzs":
		field = "增值税"
	case "xfs":
		field = "消费税"
	case "qysds":
		field = "企业所得税"
	case "grsds":
		field = "个人所得税"
	case "tdzzs":
		field = "土地增值税"
	case "yys":
		field = "营业税"
	case "fcs":
		field = "房产税"
	case "yhs":
		field = "印花税"
	case "hjbhs":
		field = "环境保护税"
	case "ccs":
		field = "车船税"
	case "cswhjss":
		field = "城市维护建设税"
	case "zys":
		field = "资源税"
	case "cztdsys":
		field = "城镇土地使用税"
	case "gdzys":
		field = "耕地占用税"
	case "qs":
		field = "契税"
	case "qtsssr":
		field = "其他税收收入"

	case "jcajbh":
		field = "稽查案件编号"
	case "jcyr":
		field = "检查人员"

	case "pgajbh":
		field = "评估案件编号"
	case "pgyr":
		field = "评估案件编号"

	case "sbhjbz":
		field = "土地增值税申报汇缴标志"
	case "xmmc":
		field = "项目名称"

	case "nsxydj":
		field = "纳税信用等级为D"

	case "fpdm":
		field = "发票代码"
	case "fphm":
		field = "发票号码"
	case "hsyr":
		field = "核实人员"

	case "fcdz":
		field = "房产地址"
	case "fcbh":
		field = "房产编号"

	case "tddz":
		field = "土地地址"
	case "tdbh":
		field = "土地编号"

	case "qtxzxx":
		field = "其他限制信息"

	case "jbr_name":
		field = "经办人"
	case "jbr_deptname":
		field = "经办人部门"
	case "qrrq":
		field = "确认日期"

	case "zxry":
		field = "注销人员"
	case "zxsj":
		field = "金三注销时间"

	default:
		field = field
	}

	return field
}
