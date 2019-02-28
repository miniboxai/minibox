package models

type Project struct {
	Model
	Name      string `gorm:"not null"`
	Summary   string
	Namespace string `gorm:"not null"`
	Private   bool
	UserID    uint
	User      User
}
