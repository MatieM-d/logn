package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Экспорт в CSV файл
func ExportCSV(vault *Vault, exportPath string) error {
	file, err := os.Create(exportPath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Заголовок
	if err := writer.Write([]string{"service", "login", "password", "note"}); err != nil {
		return err
	}

	// Записи
	for _, entry := range vault.Entries {
		row := []string{entry.Service, entry.Login, entry.Password, entry.Note}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// Генерация пути для экспорта с датой
func DefaultExportPath() string {
	name := fmt.Sprintf("logn-export-%s.csv", time.Now().Format("2006-01-02_15-04-05"))
	return filepath.Join(`D:\Projects\logn`, name)
}
