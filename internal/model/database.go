package model

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = initSQLite()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&User{},
		&UserDevice{},
		&Friend{},
		&FriendGroup{},
		&Group{},
		&GroupMember{},
		&GroupAdmin{},
		&GroupJoinRequest{},
		&Message{},
		&Favorite{},
		&FriendRequest{},
		&GroupRequest{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	initOfficialGroup()

	log.Println("Database initialized successfully")
}

func initOfficialGroup() {
	var officialGroup Group
	if err := DB.Where("group_number = ?", "2182344375").First(&officialGroup).Error; err != nil {
		var officialUser User
		if err := DB.Where("account = ?", "2182344375").First(&officialUser).Error; err != nil {
			officialUser = User{
				QQNumber: "2182344375",
				Password: "",
				Nickname: "ReCh官方",
				Avatar:   "",
			}
			DB.Create(&officialUser)
		}

		officialGroup = Group{
			GroupNumber: "2182344375",
			Name:        "ReCh官方交流群",
			Description: "ReCh官方交流群，欢迎加入！",
			OwnerID:     officialUser.ID,
			Avatar:      "",
		}
		DB.Create(&officialGroup)
	}
}

func initSQLite() (*gorm.DB, error) {
	dbDir := "./data"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}
	dbPath := filepath.Join(dbDir, "qq_like.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	log.Printf("Using SQLite database at: %s", dbPath)
	return db, nil
}
