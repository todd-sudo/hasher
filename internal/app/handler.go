package app

import (
	"context"
	"crypto/rsa"
	"fmt"
	"hasher/internal/config"
	"log"

	"github.com/fatih/color"
)

type handler interface {
	login(ctx context.Context)
	getAllSecrets(ctx context.Context, lastID int) int
	insertSecret(ctx context.Context)
	getSecretByID(ctx context.Context, secretID int)
}

type Handler struct {
	privateKey *rsa.PrivateKey
	ctx        context.Context
	cfg        config.Config
	service    service
}

func NewHandler(
	ctx context.Context,
	service service,
	cfg config.Config,
	privateKey *rsa.PrivateKey,
) handler {
	return &Handler{ctx: ctx, service: service, cfg: cfg, privateKey: privateKey}
}

func (h *Handler) login(ctx context.Context) {
	username, password, err := loginReader()
	if err != nil {
		log.Panicln(err)
	}
	_, err = h.service.getUser(ctx, username, password)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Вы вошли!")
}

func (h *Handler) getAllSecrets(ctx context.Context, lastID int) int {
	clearTerminal()
	color.Yellow("Перемещайся по секретам с помощью кнопок.\nЧтобы выбрать секрет введи: $ID_SECRET\n\n")
	secrets, err := h.service.getAllSecrets(ctx, h.cfg.Limit, lastID)
	if err != nil {
		log.Panicln("Ошибка при получении данных")
	}
	if secrets != nil {
		for _, secret := range secrets {
			color.Yellow("ID: %d | %s\n", secret.ID, secret.Title)
		}
		lastSecret := secrets[len(secrets)-1]
		if len(secrets) < h.cfg.Limit {
			return 0
		}
		return int(lastSecret.ID)
	}
	color.Red("У вас нет секретов!")
	return 0

}

func (h *Handler) insertSecret(ctx context.Context) {
	title, content := createSecretReader()
	contentEncrypt, err := rsaEncrypt(content, h.privateKey.PublicKey)
	if err != nil {
		color.Red("Ошибка зашифровки контента. Слишком длинный контент")
		return
	}
	secretDTO := CreateSecretDTO{
		Title:   title,
		Content: contentEncrypt,
	}
	if err := h.service.insertSecret(ctx, &secretDTO); err != nil {
		log.Panicln("Ошибка сохранения секрета")
	}
	color.Blue("Секрет сохранен")
}

func (h *Handler) getSecretByID(ctx context.Context, secretID int) {
	for {
		mag := color.New(color.FgMagenta)
		secret, err := h.service.getSecretByID(ctx, secretID)
		if err != nil {
			log.Panicln("Ошибка при получение секрета", err)
		}
		secretId := secret.ID
		secretTitle := secret.Title
		secretCreateAt := secret.CreatedAt.String()

		secretContentEncode := secret.Content

		secretContent, err := rsaDecrypt(secretContentEncode, *h.privateKey)
		if err != nil {
			log.Panicln("Ошибка расшифровки контента", err)
		}
		// secretContent := secret.Content
		clearTerminal()
		mag.Printf("%d | %s\n\n%v\n\n%s\n", secretId, secretTitle, secretContent, secretCreateAt)
		var state string
		fmt.Print("\n\nНажми q чтобы перейти назад ...\n\n>> ")
		fmt.Scan(&state)
		if state == "q" {
			return
		}
		clearTerminal()
		continue
	}

}
