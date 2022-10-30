package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"golang.org/x/term"
)

type User struct {
	username string
	password string
}

func createUser(username string, password string) *User {
	newUser := User{
		username: username,
		password: password,
	}

	fmt.Println("")
	fmt.Println("Login credentials added")
	return &newUser
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

func getPassword(r *bufio.Reader) string {
	password, err := getInput("Enter Salesforce Password (text will not be displayed): ", true, r)

	if err != nil {
		printErrMessage("Please enter valid password")
		getPassword(r)
	}

	return password
}

func getUsername(r *bufio.Reader) string {
	username, err := getInput("Enter Salesforce Username: ", false, r)
	if err != nil {
		printErrMessage("Please enter valid username")
		getUsername(r)
	}

	return username
}

func printErrMessage(message string) {
	var ERROR_COLOR pterm.RGB = pterm.NewRGB(255, 0, 0)
	ERROR_COLOR.Println(message)
}

func CheckIfAutoSubmit() bool {
	fmt.Println("Would you like to auto-submit the timesheet? (Best for regular work week i.e. no PTO or holidays)")
	shouldAutoSubmit, _ := pterm.DefaultInteractiveConfirm.Show()
	if !shouldAutoSubmit {
		fmt.Println("")
		fmt.Println("After logging in, you will have 10 minutes to update and manually submit your timesheet")
		fmt.Println("")
	}
	return shouldAutoSubmit
}

func PrintWelcomeMessage() {
	var INRHYTHM_COLOR pterm.RGB = pterm.NewRGB(241, 91, 33)
	bigText, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromStringWithRGB("InRhythm", INRHYTHM_COLOR)).Srender()
	pterm.DefaultCenter.Println(bigText)
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).WithTextStyle(pterm.NewStyle(pterm.Bold)).Println(" Salesforce Timesheet Automation")
	fmt.Println("")
}

func PromptForCredentials() (User, error) {
	reader := bufio.NewReader(os.Stdin)

	// TODO: add tests here
	username := getUsername(reader)
	password := getPassword(reader)

	newUser := createUser(username, password)
	return *newUser, nil
}
