package app

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type storage interface {
	saveUser(ctx context.Context, user *User) (*User, error)
	getUser(ctx context.Context, username string) (*User, error)
	getAllSecrets(ctx context.Context, limit int, externalID, createdAt string) ([]*Secret, error)
	insertSecret(ctx context.Context, secret *Secret) error
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
	tx := db.connection.WithContext(ctx).Debug()
	var user *User
	res := tx.Where(`username = ?`, username).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (db *Storage) getAllSecrets(ctx context.Context, limit int, externalID, createdAt string) ([]*Secret, error) {
	tx := db.connection.WithContext(ctx).Debug()
	var secrets []*Secret

	if externalID == "" && createdAt == "" {
		if err := tx.Limit(limit).Order("id DESC").Find(&secrets).Error; err != nil {
			return nil, err
		}
		return secrets, nil
	}

	if err := tx.Where(
		`(external_id, created_at) < (?, ?)`,
		externalID,
		createdAt,
	).Limit(limit).Order("id DESC").Find(&secrets).Error; err != nil {
		return nil, err
	}
	return secrets, nil
}

func (db *Storage) insertSecret(ctx context.Context, secret *Secret) error {
	tx := db.connection.WithContext(ctx)
	log.Println(secret)
	if err := tx.Save(&secret).Error; err != nil {
		return err
	}
	return nil
}
