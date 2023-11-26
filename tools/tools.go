package tools

//go:generate mockgen -source=tools/tools.go -destination=tools/tools.mock.gen.go -package=tools

type HashInterface interface {
	HashPassword(password string) ([]byte, error)
	ComparePassword(hashedPassword []byte, password string) error
}

type JWTInterface interface {
	Generate(content interface{}) (string, error)
	Validate(tokenString string) (interface{}, error)
}
