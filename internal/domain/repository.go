package domain

type LinkRepository interface {
	Save(link Link) error
	FindByCode(code string) (Link, error)
}