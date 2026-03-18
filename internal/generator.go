package internal

import (
	"crypto/rand"
	"math/big"
)

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits    = "0123456789"
	symbols   = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

type GeneratorOptions struct {
	Length     int
	Uppercase  bool
	Digits     bool
	Symbols    bool
}

// Генерация пароля с заданными параметрами
func Generate(opts GeneratorOptions) (string, error) {
	charset := lowercase

	if opts.Uppercase {
		charset += uppercase
	}
	if opts.Digits {
		charset += digits
	}
	if opts.Symbols {
		charset += symbols
	}

	password := make([]byte, opts.Length)
	for i := range password {
		// Криптографически случайный индекс
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[index.Int64()]
	}

	return string(password), nil
}