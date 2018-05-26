package api_test

import (
	"os"

	"github.com/cubixle/accounts/internal/context"
	"github.com/cubixle/accounts/internal/models"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func appContext() (*context.AppContext, func()) {
	db, tearDown := setupTestDB()
	return &context.AppContext{
		DB:     db,
		Logger: logrus.New(),
	}, tearDown
}
func setupTestDB() (*gorm.DB, func()) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		models.Team{},
		models.User{},
	)

	return db, func() {
		os.Remove("test.db")
	}
}
