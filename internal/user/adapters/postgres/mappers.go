package postgres

import (
	"github.com/yourusername/go-scaffolding/internal/user/domain"
)

// ToUserModel converts a domain.User to a UserModel
func ToUserModel(user *domain.User) *UserModel {
	if user == nil {
		return nil
	}

	return &UserModel{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToDomainUser converts a UserModel to a domain.User
func ToDomainUser(model *UserModel) *domain.User {
	if model == nil {
		return nil
	}

	return &domain.User{
		ID:        model.ID,
		Email:     model.Email,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// ToDomainUsers converts a slice of UserModel to a slice of domain.User
func ToDomainUsers(models []*UserModel) []*domain.User {
	if models == nil {
		return nil
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = ToDomainUser(model)
	}

	return users
}
