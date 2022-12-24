package app

import (
	"context"
	"hasher/internal/config"
	"hasher/internal/db"
	"log"
	"os"

	"gorm.io/gorm"
)

// import (
// 	"fmt"

// 	"github.com/fatih/color"
// )

func RunApplication() {
	ctx := context.Background()
	cfg := config.GetConfig()
	_db, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Panicln("Ошибка подключения базы данных")
	}
	if err := migrate(_db); err != nil {
		log.Panicln("Ошибка миграции базы данных")
	}
	_storage := NewStorage(ctx, _db)
	_service := NewService(ctx, _storage)
	_handler := NewHandler(ctx, _service, *cfg)
	runner(ctx, _handler)
}

func runner(ctx context.Context, _handler handler) {
	_handler.login(ctx)
	for {
		state := rootReader()
		switch state {
		case "1":
			_handler.getAllSecrets(ctx)
			return
		case "2":
			_handler.insertSecret(ctx)
			return
		case "3":
			os.Exit(0)
		}
	}
}

func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Secret{}); err != nil {
		return err
	}
	return nil
}
