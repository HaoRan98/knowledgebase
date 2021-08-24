package models

type Member struct {
	ID        string `json:"id"`
	GroupID   string `json:"groupId"`
	GroupName string `json:"groupName"`
	UserName  string `json:"userName"`
	Account   string `json:"account"`
	DeptID    string `json:"deptid"`
	deptname  string `json:"deptname"`
	JGMC      string `json:"jgmc"`
	JGDM      string `json:"jgdm"`
	Status    int    `json:"status"`
}

// 创建群成员
func CreateGroupMem(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

// 成员是否存在
func MemberIsExit(groupid, account string) (bool, error) {
	var member []Member
	if err := db.
		Where("groupid = ?", groupid).
		Where("account = ?", account).
		Find(&member).Error; err != nil {
		return false, err
	}

	if len(member) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// 退群
func DropOut(groupid, account string) error {
	if err := db.Table("member").Debug().
		Where("groupid = ?", groupid).
		Where("account = ?", account).
		Delete(&Member{}).Error; err != nil {
		return err
	}

	return nil
}

// 加入的群
type Qun struct {
	GroupID   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

// 我加入的团队
func MyJoinedGroup(account string) ([]Qun, error) {

	var qun []Qun

	if err := db.Debug().Table("member").
		Where("account = ?", account).
		Where("status = ?", 0).Find(&qun).Error; err != nil {
		return qun, err
	}
	return qun, nil
}

func GetMembers(groupid string) ([]Member, error) {
	var member []Member
	if err := db.
		Where("groupid = ?", groupid).
		Where("status = ?", 0).
		Find(&member).Error; err != nil {
		return member, err
	}

	return member, nil
}
