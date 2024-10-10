package service

import (
	"github.com/Clymba/testTask/internal/repository"
)

type Service struct {
	repository *repository.Repository
}

func New(repository *repository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}
