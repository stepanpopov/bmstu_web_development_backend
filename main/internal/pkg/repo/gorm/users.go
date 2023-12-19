package repo

import (
	"fmt"
	"rip/internal/pkg/repo"

	"github.com/google/uuid"
)

func (r *Repository) CreateUser(username, passwordHash string, isModerator bool) (uuid.UUID, error) {
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	user := &repo.User{
		UserID:      randomUUID,
		Username:    username,
		Password:    passwordHash,
		IsModerator: isModerator,
	}

	if err := r.db.Create(user).Error; err != nil {
		return uuid.Nil, err
	}

	return user.UserID, nil
}

func (r *Repository) CheckUser(username, passwordHash string) (uuid.UUID, bool, error) {
	user := &repo.User{
		Username: username,
		Password: passwordHash,
	}

	if err := r.db.Where(user).Take(user).Error; err != nil {
		return uuid.Nil, false, err
	}
	fmt.Println("repo: ", user.IsModerator)
	return user.UserID, user.IsModerator, nil
}
