package application

import (
	"shortapi/internal/domain"
	"shortapi/internal/exception"
)

type ResolveUseCase struct {
	repo domain.LinkRepository
}

func NewResolveUseCase(repo domain.LinkRepository) *ResolveUseCase {
	return &ResolveUseCase{repo: repo}
}

func (u *ResolveUseCase) Execute(code string) (string, error) {
	if code == "" {
		return "", exception.ErrEmptyCode
	}
	link, err := u.repo.FindByCode(code)
	if err != nil {
		return "", err
	}
	return link.OriginalURL(), nil
}
