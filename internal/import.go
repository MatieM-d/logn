package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type ImportResult struct {
	Added   int
	Skipped int
	Errors  int
}

// Определяет формат CSV по заголовкам
func detectFormat(headers []string) string {
	headerStr := strings.Join(headers, ",")

	if strings.Contains(headerStr, "login_username") {
		return "bitwarden"
	}
	if strings.Contains(headerStr, "grouping") {
		return "lastpass"
	}
	if strings.Contains(headerStr, "Account") {
		return "keepass"
	}
	return "logn"
}

// Конвертирует строку CSV в Entry в зависимости от формата
func rowToEntry(row []string, headers []string, format string) (*Entry, error) {
	// Создаём map заголовок → значение
	data := make(map[string]string)
	for i, header := range headers {
		if i < len(row) {
			data[header] = row[i]
		}
	}

	entry := &Entry{}

	switch format {
	case "bitwarden":
		entry.Service = data["name"]
		entry.Login = data["login_username"]
		entry.Password = data["login_password"]
		entry.Note = data["notes"]

	case "lastpass":
		entry.Service = data["name"]
		entry.Login = data["username"]
		entry.Password = data["password"]
		entry.Note = data["extra"]

	case "keepass":
		entry.Service = data["Account"]
		entry.Login = data["Login Name"]
		entry.Password = data["Password"]
		entry.Note = data["Comments"]

	case "logn":
		entry.Service = data["service"]
		entry.Login = data["login"]
		entry.Password = data["password"]
		entry.Note = data["note"]
	}

	if entry.Service == "" {
		return nil, fmt.Errorf("пустое название сервиса")
	}
	if entry.Password == "" {
		return nil, fmt.Errorf("пустой пароль для %s", entry.Service)
	}

	return entry, nil
}

// Импорт из CSV файла
func ImportCSV(vault *Vault, key []byte, filePath string, overwrite bool) (*ImportResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	// Читаем заголовки
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать заголовки: %w", err)
	}

	format := detectFormat(headers)
	result := &ImportResult{}

	// Читаем строки
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл: %w", err)
	}

	for _, row := range rows {
		entry, err := rowToEntry(row, headers, format)
		if err != nil {
			result.Errors++
			continue
		}

		// Проверяем существует ли уже такая запись
		existing, _ := Get(vault, entry.Service)
		if existing != nil {
			if !overwrite {
				result.Skipped++
				continue
			}
			// Удаляем старую запись перед добавлением новой
			Delete(vault, key, entry.Service)
		}

		if err := Add(vault, key, *entry); err != nil {
			result.Errors++
			continue
		}

		result.Added++
	}

	return result, nil
}
