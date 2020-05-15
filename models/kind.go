package models

type Kind struct {
	ID uint   `gorm:"primary_key"`
	Mc string `json:"mc"`
}

func AddKind(kind *Kind) error {
	if err := db.Create(kind).Error; err != nil {
		return err
	}
	return nil
}
func DelKind(id uint) error {
	if err := db.Where("id=?", id).Delete(Kind{}).Error; err != nil {
		return err
	}
	return nil
}
func GetKinds() ([]*Kind, error) {
	var kinds []*Kind
	if err := db.Find(&kinds).Error; err != nil {
		return nil, err
	}
	return kinds, nil
}
func IsKindExist(mc string) bool {
	var kind Kind
	if err := db.Where("mc=?", mc).First(&kind).Error; err != nil {
		return false
	}
	if kind.ID > 0 {
		return true
	}
	return false
}
