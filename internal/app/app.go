package app

import (
	"context"
	"crypto/rsa"
	"fmt"
	"hasher/internal/config"
	"hasher/internal/db"
	"log"
	"os/exec"

	"github.com/fatih/color"

	"os"

	"gorm.io/gorm"
)

func RunApplication() {
	cfg := config.GetConfig()

	privatKey, err := checkPrivateKeyFile(*cfg)
	if err != nil {
		log.Panicln("Ошибка создания приватного ключа", err)
	}

	ctx := context.Background()

	_db, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Panicln("Ошибка подключения базы данных")
	}
	if err := migrate(_db); err != nil {
		log.Panicln("Ошибка миграции базы данных")
	}
	_storage := NewStorage(ctx, _db)
	_service := NewService(ctx, _storage)
	_handler := NewHandler(ctx, _service, *cfg, privatKey)
	runner(ctx, _handler)
}

func runner(ctx context.Context, _handler handler) {
	_handler.login(ctx)
	for {
		clearTerminal()
		state := rootReader()
		switch state {
		case "1":
			viewSecrets(ctx, _handler)
			continue
		case "2":
			_handler.insertSecret(ctx)
			continue
		case "3":
			os.Exit(0)
		}

	}
}

func viewSecrets(ctx context.Context, _handler handler) {
	green := color.New(color.FgGreen)
	var lastID int = 1

	for {
		newLastID := _handler.getAllSecrets(ctx, lastID)

		green.Print("\n\n1. Далее\n0. Выход\n\n>> ")
		var state string
		fmt.Scan(&state)

		switch state {
		case "1":

			lastID = newLastID
		case "0":
			return
		}
	}
}

func clearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func checkPrivateKeyFile(cfg config.Config) (*rsa.PrivateKey, error) {
	var privateKey *rsa.PrivateKey
	if _, err := os.Stat(cfg.PrivatePemFileName); os.IsNotExist(err) {
		privateKey, err = generateKeys(cfg)
		if err != nil {
			return nil, err
		}
		savePrivateKeyToFile(privateKey, cfg.PrivatePemFileName)
		color.Green(fmt.Sprintf("RSA ключ создан и сохранен в файл %s", cfg.PrivatePemFileName))
	} else {
		privateKey, err = uploadPrivateKey(cfg.PrivatePemFileName)
		if err != nil {
			color.Red("Ошибка чтения файла с ключом")
		}
	}
	return privateKey, nil
}

func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Secret{}); err != nil {
		return err
	}
	return nil
}
