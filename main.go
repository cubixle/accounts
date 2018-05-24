package main

import (
	"fmt"
	"log"

	"github.com/cubixle/accounts/internal/api"
	"github.com/cubixle/accounts/internal/context"
	"github.com/cubixle/accounts/internal/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	c := initAppContext()

	if err := api.New(c).Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

func initDB(dsn string, verbose bool) *gorm.DB {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.LogMode(verbose)

	db.AutoMigrate(
		models.Team{},
	)

	return db
}

func initAppContext() *context.AppContext {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("./")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %v \n", err))
	}

	v.SetDefault("debug", false)

	db := initDB(v.GetString("db.dsn"), v.GetBool("debug"))

	logger := logrus.New()
	if v.GetBool("debug") {
		logger.SetLevel(logrus.DebugLevel)
	}

	return &context.AppContext{DB: db, Logger: logger}
}
