package app

import (
	"context"
	"log"
)

type handler interface {
	login(ctx context.Context)
}

type Handler struct {
	ctx     context.Context
	service *Service
}

func NewHandler(ctx context.Context, service *Service) handler {
	return &Handler{ctx: ctx, service: service}
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
