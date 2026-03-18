package internal

import (
	"time"

	"github.com/atotto/clipboard"
)

const clearDelay = 10 * time.Second

// Копирует пароль в буфер и очищает через 10 секунд
func CopyToClipboard(password string) error {
	if err := clipboard.WriteAll(password); err != nil {
		return err
	}

	// Очищаем буфер в фоне
	go func() {
		time.Sleep(clearDelay)
		clipboard.WriteAll("")
	}()

	return nil
}
