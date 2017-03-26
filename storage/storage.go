package storage

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
)

var MongoDB *mgo.Session
var Gorm *gorm.DB
