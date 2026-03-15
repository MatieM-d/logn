package internal

// Одна запись в хранилище
type Entry struct {
	Service  string `json:"service"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Note     string `json:"note"`
}

// Структура .vault файла на диске
type VaultFile struct {
	Salt string `json:"salt"` // base64 соль для Argon2id
	Data string `json:"data"` // base64 зашифрованный blob
}

// Расшифрованное содержимое хранилища
type Vault struct {
	Entries []Entry `json:"entries"`
}