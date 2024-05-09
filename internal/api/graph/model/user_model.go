package model

type User struct {
	ID       string `json:"id" gorm:"ype:varchar(255);primarykey"`
	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	Email    string `json:"email" gorm:"type:varchar(100);not null;unique;index"`
	Password string `json:"password" gorm:"type:varchar(100);not null"`
}
