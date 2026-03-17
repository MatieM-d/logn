package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

// Параметры Argon2id
const (
	argonTime    = 2          // количество итераций
	argonMemory  = 64 * 1024  // 64MB памяти
	argonThreads = 4          // количество потоков
	argonKeyLen  = 32         // длина ключа 32 байта = 256 бит
	saltLen      = 16         // длина соли в байтах
)

// Генерация случайной соли
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// Деривация ключа из мастер-пароля через Argon2id
func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)
}

// Шифрование данных через AES-256-GCM
func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Генерируем случайный nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Шифруем и добавляем nonce в начало результата
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Дешифрование данных
func Decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Извлекаем nonce из начала данных
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("данные повреждены")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("неверный мастер-пароль или данные повреждены")
	}

	return plaintext, nil
}