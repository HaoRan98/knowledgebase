package models

import (
	"NULL/knowledgebase/pkg/setting"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	"log"
)

// 监控项目--项目代码表
type JkxmMcdm struct {
	ID string `json:"id" gorm:"primary_key"`
	Mc string `json:"mc" gorm:"COMMENT:'类别名称'"`
	Dm string `json:"dm" gorm:"COMMENT:'类别代码'"`
}

func GetJkxmMcdms() ([]*JkxmMcdm, error) {
	var mcDms []*JkxmMcdm
	err := db.Find(&mcDms).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return mcDms, nil
}

func GetJkxmMcByDm(dm string) string {
	var mcDm JkxmMcdm
	if err := db.First(&mcDm).Error; err != nil {
		return dm
	}
	return mcDm.Mc
}

func InitXmdm() {
	inFile := setting.AppSetting.RuntimeRootPath +
		setting.AppSetting.ExportSavePath + "jkxm_mcdm.xlsx"
	xlsx, err := excelize.OpenFile(inFile)
	if err != nil {
		log.Println("OpenFile err: ", err)
		return
	}
	sheetName := xlsx.GetSheetName(1)
	log.Println("sheet name: ", sheetName)
	rows := xlsx.GetRows(sheetName)
	for k, row := range rows {
		if k == 0 || k == 6 {
			continue
		}
		xm := JkxmMcdm{}
		for i, cell := range row {
			switch {
			case i == 0:
				xm.ID = cell
			case i == 1:
				xm.Mc = cell
			case i == 2:
				xm.Dm = cell
			}
		}
		//logging.Debug(fmt.Sprintf("*: %+v", xm))
		err := AddJkxmData("jkxm_mcdm", &xm)
		if err != nil {
			log.Println("AddJkxmMcdm err: ", err)
		}
	}
}
