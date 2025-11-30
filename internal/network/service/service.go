package network

import (
	"gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/internal/network/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) Service {
	return Service{
		repo,
	}
}
