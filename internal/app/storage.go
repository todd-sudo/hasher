package app

import (
	"context"

	"gorm.io/gorm"
)

type storage interface {
	saveUser(ctx context.Context, user *User) (*User, error)
	getUser(ctx context.Context, username string, password string) (*User, error)
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

func (db *Storage) getUser(ctx context.Context, username string, password string) (*User, error) {
	tx := db.connection.WithContext(ctx).Debug()
	var user *User
	res := tx.Where(`username = ? AND password = ?`, username, password).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
