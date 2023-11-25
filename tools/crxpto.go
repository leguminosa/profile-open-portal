package tools

//go:generate mockgen -source=tools/crxpto.go -destination=tools/crxpto.mock.gen.go -package=tools

type HashInterface interface {
	HashPassword(password string) ([]byte, error)
	ComparePassword(hashedPassword []byte, password string) error
}
