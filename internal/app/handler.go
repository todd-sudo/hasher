package app

import (
	"context"
	"crypto/rsa"
	"hasher/internal/config"
	"log"

	"github.com/fatih/color"
)

type handler interface {
	login(ctx context.Context)
	getAllSecrets(ctx context.Context, lastID int) int
	insertSecret(ctx context.Context)
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
		if len(secrets) < h.cfg.Limit {
			return 0
		}
		for _, secret := range secrets {
			color.Yellow("ID: %d | %s\n", secret.ID, secret.Title)
		}
		lastSecret := secrets[len(secrets)-1]
		return int(lastSecret.ID)
	}
	color.Red("У вас нет секретов!")
	return 0

}

func (h *Handler) insertSecret(ctx context.Context) {
	title, content := createSecretReader()
	secretDTO := CreateSecretDTO{
		Title:   title,
		Content: content,
	}
	if err := h.service.insertSecret(ctx, &secretDTO); err != nil {
		log.Panicln("Ошибка сохранения секрета")
	}
	color.Blue("Секрет сохранен")
}
