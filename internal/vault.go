package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode"
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

// Результат проверки пароля
type CheckResult struct {
	Service  string
	Password string
	Failed   []string // список проваленных критериев
}

// Проверка одного пароля
func CheckPassword(password string) []string {
	var failed []string

	if len(password) < 8 {
		failed = append(failed, fmt.Sprintf("слишком короткий (%d из 8 символов)", len(password)))
	}

	hasUpper := false
	for _, c := range password {
		if unicode.IsUpper(c) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		failed = append(failed, "нет заглавных букв")
	}

	hasDigit := false
	for _, c := range password {
		if unicode.IsDigit(c) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		failed = append(failed, "нет цифр")
	}

	hasSymbol := false
	symbols := "!@#$%^&*()-_=+[]{}|;:,.<>?"
	for _, c := range password {
		if strings.ContainsRune(symbols, c) {
			hasSymbol = true
			break
		}
	}
	if !hasSymbol {
		failed = append(failed, "нет спецсимволов")
	}

	return failed
}

// Проверка конкретного сервиса
func CheckOne(vault *Vault, service string) (*CheckResult, error) {
	entry, err := Get(vault, service)
	if err != nil {
		return nil, err
	}

	return &CheckResult{
		Service:  entry.Service,
		Password: entry.Password,
		Failed:   CheckPassword(entry.Password),
	}, nil
}

// Проверка всех паролей — возвращает только не прошедшие
func CheckAll(vault *Vault) []CheckResult {
	var results []CheckResult
	for _, entry := range vault.Entries {
		failed := CheckPassword(entry.Password)
		if len(failed) > 0 {
			results = append(results, CheckResult{
				Service:  entry.Service,
				Password: entry.Password,
				Failed:   failed,
			})
		}
	}
	return results
}

// Редактирование существующей записи
func Edit(vault *Vault, key []byte, service string, updated Entry) error {
	for i, e := range vault.Entries {
		if e.Service == service {
			vault.Entries[i] = updated
			return Save(vault, key)
		}
	}
	return errors.New("запись для " + service + " не найдена")
}
