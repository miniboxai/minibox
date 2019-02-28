package models

type Organization struct {
	ModelByInt
	Name        string
	Namespace   string
	AvatarURL   string
	Description string
	SiteURL     string
	HomeURL     string
	OwnerID     uint
	Owner       User
	Members     []*User `gorm:"many2many:user_organizations;"`
}

type UserOrganization struct {
	UserID         uint
	OrganizationID uint
}
