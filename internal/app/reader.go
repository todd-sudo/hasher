package app

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

const (
	helloMsg    = "Привет!\n\nЛогин и пароль не забывать!\n\nЧтобы сохранить запись, для начала залогинься\n\n"
	usernameMsg = "Username: "
	passwordMsg = "\nPassword: "
)

func loginHandler() {
	var username string
	var password string
	cyan := color.New(color.FgCyan)

	cyan.Print(usernameMsg)
	_, err := fmt.Scan(&username)
	if err != nil {
		log.Fatalln(err)
	}

	cyan.Print(passwordMsg)
	_, err = fmt.Scan(&password)
	if err != nil {
		log.Fatalln(err)
	}

	color.Red("\n%s:%s", username, password)
}

func startReader() {
	color.Green(helloMsg)
	loginHandler()
}
