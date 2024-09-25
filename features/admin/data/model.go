package data

import "gorm.io/gorm"

type Admin struct {
	*gorm.Model
	ID       string `gorm:"primary_key;type:varchar(50);not null;column:id"`
	Username string `gorm:"type:varchar(255);not null;column:username"`
	Email    string `gorm:"type:varchar(255);not null;column:email"`
	Password string `gorm:"type:varchar(255);not null;column:password"`
}

func (Admin) TableName() string {
	return "admins"
}
