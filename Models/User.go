package models
type User struct{
	ID       uint   `gorm:"primaryKey"`
	Email string `gorm:"unique;not null"`
	Password string
	Name string
}

func (User) TableName() string{
	return "users"
}