package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const configName = "config.json"

// Возвращает путь к папке конфига в домашней директории
func configPath() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".logn")
	os.MkdirAll(dir, 0700)
	return filepath.Join(dir, configName)
}

// Загружает конфиг
func LoadConfig() (*Config, error) {
	raw, err := os.ReadFile(configPath())
	if err != nil {
		return nil, errors.New("конфиг не найден")
	}

	var config Config
	if err := json.Unmarshal(raw, &config); err != nil {
		return nil, errors.New("конфиг повреждён")
	}

	return &config, nil
}

// Сохраняет конфиг
func SaveConfig(config *Config) error {
	encoded, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), encoded, 0600)
}

// Возвращает путь к .vault файлу из конфига
func vaultPath() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", errors.New("хранилище не настроено, выполните logn init")
	}
	return config.VaultPath, nil
}

// Проверяет существует ли .vault файл
func VaultExists() bool {
	path, err := vaultPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// Сохраняет зашифрованные данные в .vault файл
func SaveVault(salt []byte, data []byte) error {
	path, err := vaultPath()
	if err != nil {
		return err
	}

	vf := VaultFile{
		Salt: base64.StdEncoding.EncodeToString(salt),
		Data: base64.StdEncoding.EncodeToString(data),
	}

	encoded, err := json.Marshal(vf)
	if err != nil {
		return err
	}

	return os.WriteFile(path, encoded, 0600)
}

// Читает .vault файл с диска
func LoadVault() (salt []byte, data []byte, err error) {
	path, err := vaultPath()
	if err != nil {
		return nil, nil, err
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, errors.New("хранилище не найдено, выполните logn init")
	}

	var vf VaultFile
	if err := json.Unmarshal(raw, &vf); err != nil {
		return nil, nil, errors.New("файл хранилища повреждён")
	}

	salt, err = base64.StdEncoding.DecodeString(vf.Salt)
	if err != nil {
		return nil, nil, errors.New("соль повреждена")
	}

	data, err = base64.StdEncoding.DecodeString(vf.Data)
	if err != nil {
		return nil, nil, errors.New("данные повреждены")
	}

	return salt, data, nil
}

// Спрашивает пользователя где хранить .vault файл
func SetupVaultPath() (string, error) {
	currentDir, _ := os.Getwd()
	defaultPath := filepath.Join(currentDir, ".vault")

	fmt.Println("\nГде хранить файл хранилища?")
	fmt.Println("  1. Текущая папка —", defaultPath)
	fmt.Println("  2. Указать свой путь")
	fmt.Print("\nВыберите (1/2): ")

	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "1":
		return defaultPath, nil
	case "2":
		fmt.Print("Введите путь (например D:\\Projects\\logn\\.vault): ")
		var customPath string
		fmt.Scanln(&customPath)
		if customPath == "" {
			return "", errors.New("путь не может быть пустым")
		}
		// Создаём папку если не существует
		dir := filepath.Dir(customPath)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return "", errors.New("не удалось создать папку: " + err.Error())
		}
		return customPath, nil
	default:
		return defaultPath, nil
	}
}

// Создание резервной копии .vault файла
func BackupVault(backupPath string) error {
	path, err := vaultPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return errors.New("хранилище не найдено")
	}

	// Создаём папку если не существует
	dir := filepath.Dir(backupPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return os.WriteFile(backupPath, data, 0600)
}

// Восстановление из резервной копии
func RestoreVault(backupPath string) error {
	// Проверяем что файл резервной копии существует
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return errors.New("файл резервной копии не найден: " + backupPath)
	}

	// Проверяем что это валидный .vault файл
	var vf VaultFile
	if err := json.Unmarshal(data, &vf); err != nil {
		return errors.New("файл не является валидным хранилищем LOGN")
	}

	path, err := vaultPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
