package models

// 监控项目--欠税
type JkxmQs struct {
	ID string `json:"id" gorm:"primary_key"`
	JkxmBase
	Zzs     float64 `json:"zzs" gorm:"COMMENT:'增值税';default:0.00"`
	Xfs     float64 `json:"xfs" gorm:"COMMENT:'消费税';default:0.00"`
	Qysds   float64 `json:"qysds" gorm:"COMMENT:'企业所得税';default:0.00"`
	Grsds   float64 `json:"grsds" gorm:"COMMENT:'个人所得税';default:0.00"`
	Tdzzs   float64 `json:"tdzzs" gorm:"COMMENT:'土地增值税';default:0.00"`
	Yys     float64 `json:"yys" gorm:"COMMENT:'营业税';default:0.00"`
	Zys     float64 `json:"zys" gorm:"COMMENT:'资源税';default:0.00"`
	Fcs     float64 `json:"fcs" gorm:"COMMENT:'房产税';default:0.00"`
	Yhs     float64 `json:"yhs" gorm:"COMMENT:'印花税';default:0.00"`
	Hjbhs   float64 `json:"hjbhs" gorm:"COMMENT:'环境保护税';default:0.00"`
	Ccs     float64 `json:"ccs" gorm:"COMMENT:'车船税';default:0.00"`
	Cswhjss float64 `json:"cswhjss" gorm:"COMMENT:'城市维护建设税';default:0.00"`
	Cztdsys float64 `json:"cztdsys" gorm:"COMMENT:'城镇土地使用税';default:0.00"`
	Gdzys   float64 `json:"gdzys" gorm:"COMMENT:'耕地占用税';default:0.00"`
	Qs      float64 `json:"qs" gorm:"COMMENT:'契税';default:0.00"`
	Qtsssr  float64 `json:"qtsssr" gorm:"COMMENT:'其他税收收入';default:0.00"`
}

// 监控项目--稽查未办结
type JkxmJcwbj struct {
	ID string `json:"id" gorm:"primary_key"`
	JkxmBase
	Jcajbh string `json:"jcajbh" gorm:"COMMENT:'稽查案件编号'"`
	Jcyr   string `json:"jcyr" gorm:"COMMENT:'检查人员'"`
}

// 监控项目--评估未办结
type JkxmPgwbj struct {
	ID string `json:"id" gorm:"primary_key"`
	JkxmBase
	Pgajbh string `json:"pgajbh" gorm:"COMMENT:'评估案件编号'"`
	Pgry   string `json:"pgyr" gorm:"COMMENT:'评估人员'"`
}

// 监控项目--未进行土增汇缴
type JkxmWjxtdhj struct {
	ID string `json:"id" gorm:"primary_key"`
	JkxmBase
	Sbhjbz string `json:"sbhjbz" gorm:"COMMENT:'土地增值税申报汇缴标志'"`
	Xmmc   string `json:"jcyr" gorm:"COMMENT:'项目名称'"`
}

// 监控项目--纳税信用等级
type JkxmNsxydj struct {
	ID string `json:"id" gorm:"primary_key"`
	JkxmBase
	Nsxydj string `json:"nsxydj" gorm:"COMMENT:'纳税信用等级为D';default:'N'"`
}

// 监控项目--风险发票未处理
type JkxmFxfpwcl struct {
	ID string `json:"id" gorm:"primary_key"`
	JkxmBase
	Fpdm string `json:"fpdm" gorm:"COMMENT:'发票代码'"`
	Fphm string `json:"fphm" gorm:"COMMENT:'发票号码'"`
	Hsry string `json:"hsyr" gorm:"COMMENT:'核实人员'"`
}

// 监控项目--房产
type JkxmFc struct {
	ID string `json:"id" gorm:"primary_key"`
	JkxmBase
	Fcdz string `json:"fcdz" gorm:"COMMENT:'房产地址'"`
	Fcbh string `json:"fcbh" gorm:"COMMENT:'房产编号'"`
}

// 监控项目--土地
type JkxmTd struct {
	ID string `json:"id" gorm:"primary_key"`
	JkxmBase
	Tddz string `json:"tddz" gorm:"COMMENT:'土地地址'"`
	Tdbh string `json:"tdbh" gorm:"COMMENT:'土地编号'"`
}

// 监控项目--其他
type JkxmQt struct {
	ID string `json:"id" gorm:"primary_key"`
	JkxmBase
	Qtxzxx string `json:"qtxzxx" gorm:"COMMENT:'其他限制信息'"`
}
