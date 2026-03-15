package crypto

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/argon2"
)

type KDFParams struct {
	Memory      uint32 `json:"mem"`
	Iterations  uint32 `json:"iter"`
	Parallelism uint8  `json:"par"`
	Salt        []byte `json:"salt"`
}

func GenerateKDFParams() (*KDFParams, error) {
	salt := make([]byte, 16)

	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	return &KDFParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		Salt:        salt,
	}, nil
}

func DeriveKey(password string, p *KDFParams) []byte {
	return argon2.IDKey(
		[]byte(password),
		p.Salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		32,
	)
}
