package internal

import (
	"encoding/json"
	"errors"
	"strings"
)

// Инициализация нового хранилища
func Init(masterPassword string) error {
	if VaultExists() {
		return errors.New("хранилище уже существует")
	}

	salt, err := GenerateSalt()
	if err != nil {
		return err
	}

	key := DeriveKey(masterPassword, salt)

	vault := Vault{Entries: []Entry{}}

	data, err := json.Marshal(vault)
	if err != nil {
		return err
	}

	encrypted, err := Encrypt(data, key)
	if err != nil {
		return err
	}

	return SaveVault(salt, encrypted)
}

// Открытие хранилища — возвращает расшифрованный Vault
func Open(masterPassword string) (*Vault, []byte, error) {
	salt, encrypted, err := LoadVault()
	if err != nil {
		return nil, nil, err
	}

	key := DeriveKey(masterPassword, salt)

	data, err := Decrypt(encrypted, key)
	if err != nil {
		return nil, nil, err
	}

	var vault Vault
	if err := json.Unmarshal(data, &vault); err != nil {
		return nil, nil, errors.New("данные хранилища повреждены")
	}

	return &vault, key, nil
}

// Сохранение изменений в хранилище
func Save(vault *Vault, key []byte) error {
	salt, _, err := LoadVault()
	if err != nil {
		return err
	}

	data, err := json.Marshal(vault)
	if err != nil {
		return err
	}

	encrypted, err := Encrypt(data, key)
	if err != nil {
		return err
	}

	return SaveVault(salt, encrypted)
}

// Добавление новой записи
func Add(vault *Vault, key []byte, entry Entry) error {
	for _, e := range vault.Entries {
		if e.Service == entry.Service {
			return errors.New("запись для " + entry.Service + " уже существует")
		}
	}

	vault.Entries = append(vault.Entries, entry)
	return Save(vault, key)
}

// Получение записи по сервису
func Get(vault *Vault, service string) (*Entry, error) {
	for _, e := range vault.Entries {
		if e.Service == service {
			return &e, nil
		}
	}
	return nil, errors.New("запись для " + service + " не найдена")
}

// Удаление записи по сервису
func Delete(vault *Vault, key []byte, service string) error {
	for i, e := range vault.Entries {
		if e.Service == service {
			vault.Entries = append(vault.Entries[:i], vault.Entries[i+1:]...)
			return Save(vault, key)
		}
	}
	return errors.New("запись для " + service + " не найдена")
}

// Список всех сервисов
func List(vault *Vault) []string {
	services := make([]string, len(vault.Entries))
	for i, e := range vault.Entries {
		services[i] = e.Service
	}
	return services
}

// Поиск записей по названию сервиса
func Search(vault *Vault, query string) []Entry {
	var results []Entry
	query = strings.ToLower(query)
	for _, e := range vault.Entries {
		if strings.Contains(strings.ToLower(e.Service), query) {
			results = append(results, e)
		}
	}
	return results
}
