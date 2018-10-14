package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/server-may-cry/bubble-go/models"
)

// DB is a wrapper to dicrease dependency to gorm
type DB struct {
	gorm *gorm.DB
}

// New create DB
func New(db *gorm.DB) *DB {
	return &DB{
		gorm: db,
	}
}

// FindUser search user in DB
func (db *DB) FindUser(platformID int, externalID string) models.User {
	var user models.User
	db.gorm.Where("sys_id = ? AND ext_id = ?", platformID, externalID).First(&user)
	return user
}

// SaveUser in DB only if user exists
func (db *DB) SaveUser(user *models.User) error {
	return db.gorm.Save(user).Error
}

// CreateUser new record in DB
func (db *DB) CreateUser(user *models.User) error {
	return db.gorm.Create(user).Error
}

// CreateTransaction new record in DB
func (db *DB) CreateTransaction(transaction *models.Transaction) error {
	return db.gorm.Create(transaction).Error
}
