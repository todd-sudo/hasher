package app

import (
	"context"
	"log"
)

type service interface {
	saveUser(ctx context.Context, user *User) (*User, error)
	getUser(ctx context.Context, username string, hashedPassword string) (*User, error)
}

type Service struct {
	ctx context.Context
	st  *Storage
}

func NewService(ctx context.Context, st *Storage) service {
	return &Service{ctx: ctx, st: st}
}

func (s *Service) saveUser(ctx context.Context, user *User) (*User, error) {
	hashPasswd, err := hashedPassword(user.Password)
	if err != nil {
		log.Fatalln("Ошибка хеширования пароля")
	}
	user.Password = hashPasswd
	user, err = s.st.saveUser(ctx, user)
	if err != nil {
		log.Fatalln("Ошибка сохранения пользователя")
	}
}

func (s *Service) getUser(ctx context.Context, username string, password string) (*User, error) {
	// TODO: Сделать получение пользователя и сравнение пароля
	hashPasswd := compareHashAndPassword()
	user, err := s.st.getUser(ctx, username, hashPasswd)
	return user, nil
}
