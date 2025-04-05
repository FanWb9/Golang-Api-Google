package models

import "time"
type User struct{
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email string `json:"email" gorm:"unique;not null"`
	Password string `json:"password"`
	Name string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth *time.Time `json:"-"` // this field which insert to the database
	DateBirtStr string `json:"date_of_birth" gorm:"-"` // this field is not insert to the database
}

func (User) TableName() string{
	return "users"
}