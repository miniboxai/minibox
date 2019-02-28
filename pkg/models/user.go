package models

type User struct {
	ModelByInt
	Name          string `gorm:"unique;not null"`
	Avatar        string
	Email         string `gorm:"unique;not null"`
	Mobile        string
	DisplayName   string
	Provider      Provider
	Namespace     string `gorm:"unique;not null"`
	Projects      []Project
	Organizations []*Organization `gorm:"many2many:user_organizations;"`
	// Gender
}

type Provider struct {
	ModelByInt
	Name  string
	SubID string
}

var host = "http://localhost:14000"

func (u *User) Profile() string {
	return host + "/" + u.Namespace
}
