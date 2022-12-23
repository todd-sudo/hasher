package app

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	helloMsg    = "Привет!\n\nЛогин и пароль не забывать!\n\nЧтобы сохранить запись, для начала залогинься\n\n"
	usernameMsg = "Username: "
	passwordMsg = "Password: "
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
