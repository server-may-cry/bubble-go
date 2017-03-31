package storage

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
)

// MongoDB mgo mongodb connection
var MongoDB *mgo.Session

// Gorm orm
var Gorm *gorm.DB
