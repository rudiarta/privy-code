package model

type BaseUser struct {
	ID int `gorm:"column:id;AUTO_INCREMENT;primary_key;type:int"`
}

type User struct {
	BaseUser
	Name     string `gorm:"column:name;"`
	Email    string `gorm:"column:email;"`
	Password string `gorm:"column:password;"`
}

func (User) TableName() string {
	return "user"
}
