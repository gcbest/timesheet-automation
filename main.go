package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/playwright-community/playwright-go"
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
	// promptCredentials()
	err := playwright.Install()

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://news.ycombinator.com"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	entries, err := page.QuerySelectorAll(".athing")
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}
	for i, entry := range entries {
		titleElement, err := entry.QuerySelector("td.title > a")
		if err != nil {
			log.Fatalf("could not get title element: %v", err)
		}
		title, err := titleElement.TextContent()
		if err != nil {
			log.Fatalf("could not get text content: %v", err)
		}
		fmt.Printf("%d: %s\n", i+1, title)
	}
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
