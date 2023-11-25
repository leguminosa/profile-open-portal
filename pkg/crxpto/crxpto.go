package crxpto

//go:generate mockgen -source=pkg/crxpto/crxpto.go -destination=pkg/crxpto/crxpto.mock.gen.go -package=crxpto

type HashInterface interface {
	HashPassword(password string) ([]byte, error)
}
