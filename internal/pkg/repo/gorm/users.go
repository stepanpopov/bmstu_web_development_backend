package repo

import (
	"rip/internal/pkg/repo"

	"github.com/google/uuid"
)

func (r *Repository) CreateUser(username, passwordHash string, isModerator bool) (uuid.UUID, error) {
	user := &repo.User{
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

	if err := r.db.First(user).Error; err != nil {
		return uuid.Nil, false, err
	}

	return user.UserID, user.IsModerator, nil
}
