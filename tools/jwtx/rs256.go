package jwtx

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// SigningMethodRSA implements the RSASSA-PKCS1-v1_5 signature algorithm
type SigningMethodRS256 struct {
	publicKey  []byte
	privateKey []byte
	timeNow    func() time.Time
	ttl        time.Duration
}

type NewSigningMethodRS256Options struct {
	PrivateKey []byte
	PublicKey  []byte
}

func NewSigningMethodRS256(opts NewSigningMethodRS256Options) *SigningMethodRS256 {
	return &SigningMethodRS256{
		privateKey: opts.PrivateKey,
		publicKey:  opts.PublicKey,
		timeNow:    time.Now,
		ttl:        time.Hour * 24,
	}
}

func (j *SigningMethodRS256) Generate(content interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", err
	}

	claims := make(jwt.MapClaims)
	claims["dat"] = content
	claims["iat"] = j.timeNow().Unix()
	claims["exp"] = j.timeNow().Add(j.ttl).Unix()

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
}

var (
	// ErrInvalidToken obscures the underlying error
	// from the jwt library to avoid brute force attack.
	ErrInvalidToken = errors.New("invalid token")
)

func (j *SigningMethodRS256) Validate(tokenString string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return nil, ErrInvalidToken
	}

	token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidToken
		}
		return key, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return nil, ErrInvalidToken
	}

	return claims["dat"], nil
}
