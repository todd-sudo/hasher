package app

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type service interface {
	saveUser(ctx context.Context, user *User) (*User, error)
	getUser(ctx context.Context, username string, hashedPassword string) (*User, error)
	getAllSecrets(ctx context.Context, limit int, lastID int) ([]*Secret, error)
	insertSecret(ctx context.Context, secretDTO *CreateSecretDTO) error
}

type Service struct {
	ctx context.Context
	st  storage
}

func NewService(ctx context.Context, st storage) service {
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
	return user, nil
}

func (s *Service) getUser(
	ctx context.Context,
	username string,
	password string,
) (*User, error) {
	user, err := s.st.getUser(ctx, username)
	if err != nil || user == nil {
		user, err = s.saveUser(ctx, &User{Username: username, Password: password})
		if err != nil {
			return nil, err
		}
	}
	if err := compareHashAndPassword(user.Password, password); err != nil {
		log.Fatalf("Неверный пароль для пользователя с username: %s", username)
		return nil, err
	}
	return user, nil
}

func (s *Service) getAllSecrets(ctx context.Context, limit int, lastID int) ([]*Secret, error) {
	secrets, err := s.st.getAllSecrets(ctx, limit, lastID)
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

func (s *Service) insertSecret(ctx context.Context, secretDTO *CreateSecretDTO) error {
	secretDB := Secret{}
	if err := smapping.FillStruct(&secretDB, smapping.MapFields(secretDTO)); err != nil {
		return err
	}
	secretDB.ExternalID = uuid.NewString()
	secretDB.CreatedAt = time.Now()
	if err := s.st.insertSecret(ctx, &secretDB); err != nil {
		return err
	}
	return nil
}
