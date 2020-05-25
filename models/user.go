package models

type SysUser struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	UserAccount string `json:"user_account"`
	DepID       string `json:"dep_id"`
	DepartName  string `json:"depart_name"`
	ParentID    string `json:"parent_id"`
}
