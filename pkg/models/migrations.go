package models

func AutoMigrations(db *Database) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&OAuth2Client{})
	db.AutoMigrate(&OAuth2Token{})
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&Dataset{})
	db.AutoMigrate(&DatasetLabel{})
	db.AutoMigrate(&DatasetResource{})
	db.AutoMigrate(&DatasetStorage{})
	// db.AutoMigrate(&OAuth2Endpoint{})

}
