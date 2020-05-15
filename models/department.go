package models

type SysDepart struct {
	ID         string `json:"id"`
	ParentID   string `json:"parent_id"`
	DepartName string `json:"depart_name"`
}
