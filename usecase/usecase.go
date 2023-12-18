package usecase

import (
	"D/Works/GO/31/models"
	"context"
)

type Storage interface {
	CreateNewUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, id string) (models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	MakeFriends(ctx context.Context, id1 string, id2 string) (error, error)
	UpdateUser(ctx context.Context, id string, mu models.User) error
	DeleteUser(ctx context.Context, id string) error
	GetUserFriends(ctx context.Context, id string) ([]string, error)
}

func CreateNewUser(ctx context.Context, s Storage, mu models.User) error {
	return s.CreateNewUser(ctx, mu)
}

func MakeFriends(ctx context.Context, s Storage, id1 string, id2 string) (error, error) {
	return s.MakeFriends(ctx, id1, id2)
}

func GetUser(ctx context.Context, s Storage, id string) (models.User, error) {
	return s.GetUser(ctx, id)
}

func GetAllUsers(ctx context.Context, s Storage) ([]*models.User, error) {
	return s.GetAllUsers(ctx)
}

func GetUserFriends(ctx context.Context, s Storage, id string) ([]string, error) {
	return s.GetUserFriends(ctx, id)
}

func UpdateUser(ctx context.Context, s Storage, id string, mu models.User) error {
	return s.UpdateUser(ctx, id, mu)
}

func DeleteUser(ctx context.Context, s Storage, id string) error {
	return s.DeleteUser(ctx, id)
}
