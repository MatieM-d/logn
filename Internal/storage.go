package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
)

const vaultFile = ".vault"

// Проверяет существует ли .vault файл
func VaultExists() bool {
	_, err := os.Stat(vaultFile)
	return !errors.Is(err, os.ErrNotExist)
}

// Сохраняет зашифрованные данные в .vault файл
func SaveVault(salt []byte, data []byte) error {
	vf := VaultFile{
		Salt: base64.StdEncoding.EncodeToString(salt),
		Data: base64.StdEncoding.EncodeToString(data),
	}

	encoded, err := json.Marshal(vf)
	if err != nil {
		return err
	}

	// Записываем файл с правами 600 — только владелец может читать
	return os.WriteFile(vaultFile, encoded, 0600)
}

// Читает .vault файл с диска
func LoadVault() (salt []byte, data []byte, err error) {
	raw, err := os.ReadFile(vaultFile)
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