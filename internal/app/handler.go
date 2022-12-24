package app

import (
	"context"
	"hasher/internal/config"
	"log"

	"github.com/fatih/color"
)

type handler interface {
	login(ctx context.Context)
	getAllSecrets(ctx context.Context)
	insertSecret(ctx context.Context)
}

type Handler struct {
	ctx     context.Context
	cfg     config.Config
	service service
}

func NewHandler(ctx context.Context, service service, cfg config.Config) handler {
	return &Handler{ctx: ctx, service: service, cfg: cfg}
}

func (h *Handler) login(ctx context.Context) {
	username, password, err := loginReader()
	if err != nil {
		log.Panicln(err)
	}
	user, err := h.service.getUser(ctx, username, password)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Вы вошли!")
	log.Println(user)
}

func (h *Handler) getAllSecrets(ctx context.Context) {
	var externalID string
	var createdAt string

	externalIDCtx := ctx.Value("externalID")
	createdAtCtx := ctx.Value("createdAt")

	if externalIDCtx != nil && createdAtCtx != nil {
		externalID = externalIDCtx.(string)
		createdAt = createdAtCtx.(string)
	}

	secrets, err := h.service.getAllSecrets(ctx, h.cfg.Limit, externalID, createdAt)
	if err != nil {
		log.Panicln("Ошибка при получении данных")
	}
	log.Println(secrets)
	if secrets != nil {
		for _, secret := range secrets {
			content := secret.Content
			color.Yellow("ID: %d | %s\n\n%s\nCreate at: %s", secret.ID, secret.Title, content, secret.CreatedAt)
		}
		lastSecret := secrets[len(secrets)-1]
		_ = context.WithValue(ctx, "externalID", lastSecret.ExternalID)
		_ = context.WithValue(ctx, "createdAt", lastSecret.CreatedAt)
	}

	color.Red("У вас нет секретов!")
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
