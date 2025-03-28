package models


type User struct{
	ID       uint   `gorm:"primaryKey"`
	Email string `gorm:"unique;not null"`
	Password string
	Name string
	GoogleID string `gorm:"unique"`

}

func (User) TableName() string{
	return "users"
}