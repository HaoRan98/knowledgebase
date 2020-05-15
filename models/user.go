package models

import (
	"fmt"
	"log"
)

type SysUser struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	UserAccount string `json:"user_account"`
	DepID       string `json:"dep_id"`
	DepartName  string `json:"depart_name"`
	ParentID    string `json:"parent_id"`
}

func GetUserByID(id string) (*SysUser, error) {
	var user SysUser
	qu := `select sys_user.id,sys_user.username,sys_user.user_account,sys_user_depart.dep_id 
           from sys_user 
           left join sys_user_depart 
           on sys_user_depart.user_id=sys_user.id 
           where sys_user.id='%s'`
	squery := fmt.Sprintf(qu, id)
	rows, err := db.Raw(squery).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.UserAccount, &user.DepID)
		if err != nil {
			log.Printf("Rows Scan %v.", err)
		}
	}
	return &user, nil
}

func GetUserListByDepartmentID(deptID string) ([]*SysUser, error) {
	users := make([]*SysUser, 0)
	qu := `select sys_user.id,sys_user.username,sys_user.user_account,sys_user_depart.dep_id 
           from sys_user 
           left join sys_user_depart 
           on sys_user_depart.user_id=sys_user.id 
           where sys_user_depart.dep_id='%s'`
	squery := fmt.Sprintf(qu, deptID)
	rows, err := db.Raw(squery).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user SysUser
		err := rows.Scan(&user.ID, &user.Username, &user.UserAccount, &user.DepID)
		if err != nil {
			log.Printf("Rows Scan %v.", err)
		}
		users = append(users, &user)
	}
	return users, nil
}
