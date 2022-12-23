package app

import (
	"context"
	"hasher/internal/config"
	"hasher/internal/db"
	"log"

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
	_handler := NewHandler(ctx, _service)

	_handler.login(ctx)
	// var description string
	// fmt.Println("Как тебя зовут?")
	// var tempText string
	// for {
	// 	fmt.Scanf("%s\n", &tempText)
	// 	if tempText == "END" {
	// 		break
	// 	}
	// 	description += fmt.Sprintf("%s\n", tempText)
	// }
	// color.Cyan(description)
}

func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}
