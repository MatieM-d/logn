package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/MatieM-d/logn/internal"
	"golang.org/x/term"
)

func main() {
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
		fmt.Println("Ошибка:", err)
		return
	}

	vault, _, err := internal.Open(password)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	entry, err := internal.Get(vault, service)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Сервис:", entry.Service)
	fmt.Println("Логин: ", entry.Login)
	if entry.Note != "" {
		fmt.Println("Заметка:", entry.Note)
	}

	if err := internal.CopyToClipboard(entry.Password); err != nil {
		fmt.Println("Ошибка копирования:", err)
		return
	}

	fmt.Println("Пароль скопирован в буфер! Очистится через 10 секунд.")
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
		fmt.Println("Ошибка:", err)
		return
	}

	vault, _, err := internal.Open(password)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	services := internal.List(vault)
	if len(services) == 0 {
		fmt.Println("Хранилище пусто")
		return
	}

	fmt.Println("Сохранённые сервисы:")
	for i, service := range services {
		fmt.Println(" ", strconv.Itoa(i+1)+".", service)
	}
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

func printHelp() {
	fmt.Println(`
LOGN — менеджер паролей

Команды:
  logn init              Создать новое хранилище
  logn add <сервис>      Добавить запись
  logn get <сервис>      Получить пароль (копирует в буфер)
  logn list              Список всех сервисов
  logn delete <сервис>   Удалить запись
  logn generate          Сгенерировать пароль
	`)
}