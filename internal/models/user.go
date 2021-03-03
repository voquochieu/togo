package models

type User struct {
	ID       string `gorm:"primary_key" json:"id"`
	Password string `gorm:"size:100;not null;" json:"password"`
	MaxTodo  int    `gorm:"default:5" json:"max_todo"`
}
