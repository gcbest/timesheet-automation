package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type user struct {
	username string
	password string
}

func getInput(prompt string, isPassword bool, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	if isPassword {
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}

		password := string(bytePassword)
		return password, nil
	} else {
		input, err := r.ReadString('\n')
		return strings.TrimSpace(input), err
	}

}

func promptCredentials() (*user, error) {
	reader := bufio.NewReader(os.Stdin)

	// TODO: add tests here
	username, error := getInput("Enter username: ", false, reader)
	if error != nil {
		fmt.Println("Please enter valid username")
		promptCredentials()
	}
	password, error := getInput("Enter password: ", true, reader)

	// TODO: only prompt for password again
	if error != nil {
		fmt.Println("Please enter valid password")
		promptCredentials()
	}

	newUser := createUser(username, password)
	fmt.Println("user", newUser)
	return newUser, nil
}

func createUser(username string, password string) *user {
	newUser := user{
		username: username,
		password: password,
	}

	fmt.Println("Login credentials added")
	return &newUser
}

func main() {
	promptCredentials()
}
