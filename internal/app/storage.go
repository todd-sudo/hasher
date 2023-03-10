package app

import (
	"context"

	"gorm.io/gorm"
)

type storage interface {
	saveUser(ctx context.Context, user *User) (*User, error)
	getUser(ctx context.Context, username string) (*User, error)
	getAllSecrets(ctx context.Context, limit int, lastID int) ([]*Secret, error)
	insertSecret(ctx context.Context, secret *Secret) error
	getSecretByID(ctx context.Context, secretID int) (*Secret, error)
}

type Storage struct {
	ctx        context.Context
	connection *gorm.DB
}

func NewStorage(ctx context.Context, connection *gorm.DB) storage {
	return &Storage{ctx: ctx, connection: connection}
}

func (db *Storage) saveUser(ctx context.Context, user *User) (*User, error) {
	tx := db.connection.WithContext(ctx)
	res := tx.Save(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (db *Storage) getUser(ctx context.Context, username string) (*User, error) {
	tx := db.connection.WithContext(ctx)
	var user *User
	res := tx.Take(&user, map[string]string{"username": username})
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (db *Storage) getAllSecrets(ctx context.Context, limit int, lastID int) ([]*Secret, error) {
	tx := db.connection.WithContext(ctx)
	var secrets []*Secret

	if err := tx.Where(`id >= ?`, lastID).Limit(limit).Find(&secrets).Error; err != nil {
		return nil, err
	}
	return secrets, nil
}

func (db *Storage) insertSecret(ctx context.Context, secret *Secret) error {
	tx := db.connection.WithContext(ctx)
	if err := tx.Save(&secret).Error; err != nil {
		return err
	}
	return nil
}

func (db *Storage) getSecretByID(ctx context.Context, secretID int) (*Secret, error) {
	tx := db.connection.WithContext(ctx)
	var secret *Secret
	res := tx.Where(`id = ?`, secretID).Find(&secret)
	if res.Error != nil {
		return nil, res.Error
	}
	return secret, nil
}
