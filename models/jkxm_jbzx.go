package models

import (
	"github.com/jinzhu/gorm"
	"strings"
)

// 监控项目--即办注销企业名单
type JkxmJbzx struct {
	ID          string `json:"id" gorm:"primary_key"`
	Nsrsbh      string `json:"nsrsbh" gorm:"COMMENT:'纳税人识别号'"`
	Nsrmc       string `json:"nsrmc" gorm:"COMMENT:'纳税人名称'"`
	JbrAccount  string `json:"jbr_account" gorm:"COMMENT:'经办人员账号'"`
	JbrName     string `json:"jbr_name" gorm:"COMMENT:'经办人员姓名'"`
	JbrDepid    string `json:"jbr_depid" gorm:"COMMENT:'经办人员部门id'"`
	JbrDeptname string `json:"jbr_deptname" gorm:"COMMENT:'经办人员部门名称'"`
	Qrrq        string `json:"qrrq" gorm:"COMMENT:'确认日期'"`
}

// 监控项目--金三注销企业名单
type JkxmGt3 struct {
	Nsrsbh string `json:"nsrsbh" gorm:"COMMENT:'纳税人识别号'"`
	Nsrmc  string `json:"nsrmc" gorm:"COMMENT:'纳税人名称'"`
	Zxry   string `json:"zxry" gorm:"COMMENT:'注销人员'"`
	Zxsj   string `json:"zxsj" gorm:"COMMENT:'金税三期注销时间'"`
}

func IsNsrsbhExist(nsrsbh string) bool {
	var jbzx JkxmJbzx
	if err := db.Table("jkxm_jbzx").
		Where("nsrsbh=?", nsrsbh).First(&jbzx).Error; err != nil {
		return false
	}
	return true
}

func JbZxJkxm(data interface{}) error {
	if err := db.Table("jkxm_jbzx").Create(data).Error; err != nil {
		return err
	}
	return nil
}

func GetJbzxs(cond string) ([]*JkxmJbzx, error) {
	var jbzxs []*JkxmJbzx
	err := db.Table("jkxm_jbzx").Where(cond).Find(&jbzxs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return jbzxs, nil
}

func SortFields(tHeader map[string]string) []string {
	sorted_keys, tempMap := make([]string, 5), make([]string, 5)
	for filed := range tHeader {
		filed = strings.ToLower(filed)
		if filed == "id" ||
			strings.Contains(filed, "_account") ||
			strings.Contains(filed, "_depid") {
			continue
		}
		switch filed {
		case "nsrsbh":
			sorted_keys[0] = filed
		case "nsrmc":
			sorted_keys[1] = filed

		case "jbr_name":
			sorted_keys[2] = filed
		case "jbr_deptname":
			sorted_keys[3] = filed
		case "qrrq":
			sorted_keys[4] = filed

		case "zxry":
			sorted_keys[2] = filed
		case "zxsj":
			sorted_keys[3] = filed

		default:
			tempMap = append(tempMap, filed)
		}
	}
	sorted_keys = append(sorted_keys, tempMap...)
	return sorted_keys
}
