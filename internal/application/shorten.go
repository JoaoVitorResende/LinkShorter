package application

import "shortapi/internal/domain"

type CodeGenerator interface {
	Generate() string
}

type ShortenUseCase struct {
	repo domain.LinkRepository
	gen  CodeGenerator
}

func NewShortenUseCase(repo domain.LinkRepository, gen CodeGenerator) *ShortenUseCase {
	return &ShortenUseCase{repo: repo, gen: gen}
}

func (u *ShortenUseCase) Execute(rawURL string) (string, error) {
	link, err := domain.NewLink(u.gen.Generate(), rawURL)

	if err != nil {
		return "", err
	}

	if err := u.repo.Save(link); err != nil {
		return "", err
	}
	
	return link.Code(), nil
}