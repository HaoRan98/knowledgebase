package models

type SysUser struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	UserAccount string  `json:"user_account"`
	DepID       float64 `json:"dep_id"`
	DepartName  string  `json:"depart_name"`
	ParentID    float64 `json:"parent_id"`
}
