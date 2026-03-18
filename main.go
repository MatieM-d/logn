package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/MatieM-d/logn/internal"
	"golang.org/x/term"
)

func main() {
	internal.InitColors()

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		cmdInit()
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Использование: logn add <сервис>")
			return
		}
		cmdAdd(os.Args[2])
	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Использование: logn get <сервис>")
			return
		}
		cmdGet(os.Args[2])
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Использование: logn delete <сервис>")
			return
		}
		cmdDelete(os.Args[2])
	case "list":
		cmdList()
	case "generate":
		cmdGenerate()
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Использование: logn search <запрос>")
			return
		}
		cmdSearch(os.Args[2])
	case "check":
		if len(os.Args) < 3 {
			cmdCheckAll()
		} else {
			cmdCheckOne(os.Args[2])
		}
	case "backup":
		cmdBackup()
	case "restore":
		if len(os.Args) < 3 {
			fmt.Println("Использование: logn restore <путь>")
			return
		}
		cmdRestore(os.Args[2])
	default:
		fmt.Println("Неизвестная команда:", command)
		printHelp()
	}
}

// Ввод мастер-пароля без отображения на экране
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(password), nil
}

func cmdInit() {
	// Выбор пути к хранилищу
	vaultPath, err := internal.SetupVaultPath()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Сохраняем конфиг
	config := &internal.Config{VaultPath: vaultPath}
	if err := internal.SaveConfig(config); err != nil {
		fmt.Println("Ошибка сохранения конфига:", err)
		return
	}

	password, err := readPassword("Введите мастер-пароль: ")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	confirm, err := readPassword("Подтвердите мастер-пароль: ")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	if password != confirm {
		fmt.Println("Пароли не совпадают")
		return
	}

	if err := internal.Init(password); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Хранилище LOGN успешно создано!")
	fmt.Println("Путь:", vaultPath)
}

func cmdAdd(service string) {
	password, err := readPassword("Мастер-пароль: ")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	vault, key, err := internal.Open(password)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Print("Логин: ")
	var login string
	fmt.Scanln(&login)

	fmt.Print("Пароль (оставьте пустым для генерации): ")
	var entryPassword string
	fmt.Scanln(&entryPassword)

	if entryPassword == "" {
		entryPassword, err = internal.Generate(internal.GeneratorOptions{
			Length:    20,
			Uppercase: true,
			Digits:    true,
			Symbols:   true,
		})
		if err != nil {
			fmt.Println("Ошибка генерации пароля:", err)
			return
		}
		fmt.Println("Сгенерирован пароль:", entryPassword)
	}

	fmt.Print("Заметка (необязательно): ")
	var note string
	fmt.Scanln(&note)

	entry := internal.Entry{
		Service:  service,
		Login:    login,
		Password: entryPassword,
		Note:     note,
	}

	if err := internal.Add(vault, key, entry); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Запись для", service, "добавлена!")
}

func cmdGet(service string) {
	password, err := readPassword("Мастер-пароль: ")
	if err != nil {
		internal.Error(err.Error())
		return
	}

	vault, _, err := internal.Open(password)
	if err != nil {
		internal.Error(err.Error())
		return
	}

	entry, err := internal.Get(vault, service)
	if err != nil {
		internal.Error(err.Error())
		return
	}

	fmt.Println(internal.Separator())
	fmt.Printf("Сервис: %s\n", internal.Blue(internal.Bold(entry.Service)))
	fmt.Printf("Логин:  %s\n", internal.White(entry.Login))
	if entry.Note != "" {
		fmt.Printf("Заметка: %s\n", internal.Yellow(entry.Note))
	}
	fmt.Println(internal.Separator())

	if err := internal.CopyToClipboard(entry.Password); err != nil {
		internal.Error("Ошибка копирования: " + err.Error())
		return
	}

	internal.Success("Пароль скопирован в буфер! Очистится через 10 секунд.")
}

func cmdDelete(service string) {
	password, err := readPassword("Мастер-пароль: ")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	vault, key, err := internal.Open(password)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	if err := internal.Delete(vault, key, service); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Запись для", service, "удалена!")
}

func cmdList() {
	password, err := readPassword("Мастер-пароль: ")
	if err != nil {
		internal.Error(err.Error())
		return
	}

	vault, _, err := internal.Open(password)
	if err != nil {
		internal.Error(err.Error())
		return
	}

	if len(vault.Entries) == 0 {
		fmt.Println(internal.Yellow("Хранилище пусто"))
		return
	}

	fmt.Println(internal.Bold("\nСохранённые записи:"))
	for i, entry := range vault.Entries {
		fmt.Println(internal.Separator())
		fmt.Printf("%s %s\n", internal.Gray(fmt.Sprintf("%d.", i+1)), internal.Blue(internal.Bold(entry.Service)))
		fmt.Printf("   Логин:   %s\n", internal.White(entry.Login))
		fmt.Printf("   Пароль:  %s\n", colorPassword(entry.Password))
		if entry.Note != "" {
			fmt.Printf("   Заметка: %s\n", internal.Yellow(entry.Note))
		}
	}
	fmt.Println(internal.Separator())
}

func cmdGenerate() {
	password, err := internal.Generate(internal.GeneratorOptions{
		Length:    20,
		Uppercase: true,
		Digits:    true,
		Symbols:   true,
	})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Сгенерированный пароль:", password)
}

func cmdSearch(query string) {
	password, err := readPassword("Мастер-пароль: ")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	vault, _, err := internal.Open(password)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	results := internal.Search(vault, query)
	if len(results) == 0 {
		fmt.Println("Ничего не найдено по запросу:", query)
		return
	}

	fmt.Printf("\nНайдено записей: %d\n", len(results))
	fmt.Println("─────────────────────────────────────")
	for i, entry := range results {
		fmt.Printf("%d. Сервис:  %s\n", i+1, entry.Service)
		fmt.Printf("   Логин:   %s\n", entry.Login)
		fmt.Printf("   Пароль:  %s\n", entry.Password)
		if entry.Note != "" {
			fmt.Printf("   Заметка: %s\n", entry.Note)
		}
		fmt.Println("─────────────────────────────────────")
	}
}

func cmdCheckOne(service string) {
	password, err := readPassword("Мастер-пароль: ")
	if err != nil {
		internal.Error(err.Error())
		return
	}

	vault, _, err := internal.Open(password)
	if err != nil {
		internal.Error(err.Error())
		return
	}

	result, err := internal.CheckOne(vault, service)
	if err != nil {
		internal.Error(err.Error())
		return
	}

	fmt.Println(internal.Separator())
	fmt.Printf("Сервис: %s\n", internal.Blue(internal.Bold(result.Service)))
	fmt.Printf("Пароль: %s\n", colorPassword(result.Password))

	if len(result.Failed) == 0 {
		internal.Success("Пароль прошёл проверку")
	} else {
		internal.Error("Пароль не прошёл проверку:")
		for _, f := range result.Failed {
			fmt.Println(internal.Red("  — " + f))
		}
	}
	fmt.Println(internal.Separator())
}

func cmdCheckAll() {
	password, err := readPassword("Мастер-пароль: ")
	if err != nil {
		internal.Error(err.Error())
		return
	}

	vault, _, err := internal.Open(password)
	if err != nil {
		internal.Error(err.Error())
		return
	}

	results := internal.CheckAll(vault)
	if len(results) == 0 {
		internal.Success("Все пароли прошли проверку!")
		return
	}

	fmt.Println(internal.Red(internal.Bold(fmt.Sprintf("\nНайдено слабых паролей: %d", len(results)))))
	for _, result := range results {
		fmt.Println(internal.Separator())
		fmt.Printf("Сервис: %s\n", internal.Blue(internal.Bold(result.Service)))
		fmt.Printf("Пароль: %s\n", colorPassword(result.Password))
		fmt.Println(internal.Red("✗ Проблемы:"))
		for _, f := range result.Failed {
			fmt.Println(internal.Red("  — " + f))
		}
	}
	fmt.Println(internal.Separator())
}

func cmdBackup() {
	// Генерируем имя файла с датой
	backupName := fmt.Sprintf("logn-backup-%s.vault", time.Now().Format("2006-01-02_15-04-05"))
	backupPath := filepath.Join(`D:\Projects\logn\backups`, backupName)

	if err := internal.BackupVault(backupPath); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Резервная копия создана:", backupPath)
}

func cmdRestore(backupPath string) {
	fmt.Print("Вы уверены? Текущее хранилище будет заменено (да/нет): ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "да" {
		fmt.Println("Отменено")
		return
	}

	if err := internal.RestoreVault(backupPath); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Хранилище восстановлено из:", backupPath)
}

func colorPassword(password string) string {
	failed := internal.CheckPassword(password)
	if len(failed) == 0 {
		return internal.Green(password)
	}
	return internal.Red(password)
}

func printHelp() {
	fmt.Print(`
LOGN — менеджер паролей

Команды:
  logn init                Создать новое хранилище
  logn add <сервис>        Добавить запись
  logn get <сервис>        Получить пароль (копирует в буфер)
  logn list                Список всех записей
  logn search <запрос>     Поиск по названию сервиса
  logn check               Проверить все пароли
  logn check <сервис>      Проверить пароль сервиса
  logn backup              Создать резервную копию
  logn restore <путь>      Восстановить из резервной копии
  logn delete <сервис>     Удалить запись
  logn generate            Сгенерировать пароль
`)
}

// ТЕСТОВЫЙ КОММЕНТАРИЙ
