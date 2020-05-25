package models

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

func CountJkxms(tName, cond string) int {
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
