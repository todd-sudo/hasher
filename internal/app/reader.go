package app

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	helloMsg    = "Привет!\n\nЛогин и пароль не забывать!\n\nЧтобы сохранить запись, для начала залогинься\n\n"
	usernameMsg = "Username: "
	passwordMsg = "Password: "
	rootMsg     = "1. Показать все записи\n2. Создать запись\n3. Выйти"
)

func loginReader() (string, string, error) {
	var username string
	var password string
	cyan := color.New(color.FgCyan)

	cyan.Print(usernameMsg)
	_, err := fmt.Scan(&username)
	if err != nil {
		return "", "", err
	}

	cyan.Print(passwordMsg)
	_, err = fmt.Scan(&password)
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

func rootReader() string {
	green := color.New(color.FgGreen)
	green.Printf("%s\n\n>> ", rootMsg)
	var state string
	fmt.Scanf("%s", &state)
	return state
}

// allSecretsReader - показ всех записей
func allSecretsReader() {
	color.Green("hello")
}

func createSecretReader() (string, string) {
	var title string
	var content string

	yellow := color.New(color.FgYellow)

	yellow.Print("Название: ")
	fmt.Scanf("%s\n", &title)
	yellow.Print("Контент:")
	var text string
	for {
		fmt.Scanf("%s\n", &text)
		if text == "END" {
			break
		}
		content += fmt.Sprintf("%s\n", text)
	}
	return title, content
}
