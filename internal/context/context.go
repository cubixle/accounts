package context

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type AppContext struct {
	DB     *gorm.DB
	Logger *logrus.Logger
}
