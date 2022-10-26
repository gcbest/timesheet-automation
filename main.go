package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"golang.org/x/term"
)

type User struct {
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

func getPassword(r *bufio.Reader) string {
	password, err := getInput("Enter Salesforce Password: ", true, r)

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

func promptCredentials() (User, error) {
	reader := bufio.NewReader(os.Stdin)

	// TODO: add tests here
	username := getUsername(reader)
	password := getPassword(reader)

	newUser := createUser(username, password)
	// TODO: remove
	fmt.Println("user", newUser)
	return *newUser, nil
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

func logInUser(page *rod.Page, user User, shouldAutoSubmit bool) {
	// Navigated to Choose a Username Page
	page.Race().Element("ul#idlist span").MustHandle(func(e *rod.Element) {
		displayedUsername, _ := e.Text()
		if displayedUsername == user.username {
			e.MustClick()
			handleCredentialsPage(page, user, shouldAutoSubmit)
		} else {
			printErrMessage("Log in unsuccessful, please enter credentials again")
			page.Browser().Close()
			main()
		}
		// Navigated to Login Page
	}).Element("input#username").MustHandle(func(e *rod.Element) {
		handleCredentialsPage(page, user, shouldAutoSubmit)
		// Navigated to Homepage
	}).Element("div.slds-no-print.oneAppNavContainer").MustHandle(func(e *rod.Element) {
		handleSubmittingTime(page, shouldAutoSubmit)
	}).MustDo()
}

func handleCredentialsPage(page *rod.Page, user User, shouldAutoSubmit bool) {
	submitUserCredentials(page, user)
	handleLoginNavigation(page, user, shouldAutoSubmit)
}

func submitUserCredentials(page *rod.Page, user User) {
	page.MustElement("input#username").MustSelectAllText().MustInput("").MustInput(user.username)
	page.MustElement("input#password").MustSelectAllText().MustInput("").MustInput(user.password)
	rememberMeCheckBox := page.MustElement("input#rememberUn")
	if !rememberMeCheckBox.MustProperty("checked").Bool() {
		rememberMeCheckBox.MustClick()
	}
	page.MustElement("input#Login").MustClick()
}

func printWelcomeMessage() {
	var INRHYTHM_COLOR pterm.RGB = pterm.NewRGB(241, 91, 33)
	bigText, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromStringWithRGB("InRhythm", INRHYTHM_COLOR)).Srender()
	pterm.DefaultCenter.Println(bigText)
	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).WithTextStyle(pterm.NewStyle(pterm.Bold)).Println(" Salesforce Timesheet Automation")
}

func main() {
	printWelcomeMessage()
	fmt.Println("")
	shouldAutoSubmit := checkIfAutoSubmit()
	getCredentialsAndLogin(shouldAutoSubmit)
}

func getCredentialsAndLogin(shouldAutoSubmit bool) {
	user, _ := promptCredentials()

	spinnerInfo, _ := pterm.DefaultSpinner.Start("Launching browser...")

	url := launcher.New().
		UserDataDir("path").
		Headless(false).
		MustLaunch()

	page := rod.New().ControlURL(url).MustConnect().MustPage("https://inrhythm.my.salesforce.com/")
	spinnerInfo.Success()
	logInUser(page, user, shouldAutoSubmit)
	// handleLoginNavigation(page, user, shouldAutoSubmit)
}

func handleLoginNavigation(page *rod.Page, user User, shouldAutoSubmit bool) {
	// Navigated to Homepage
	page.Race().Element("div.slds-no-print.oneAppNavContainer").MustHandle(func(e *rod.Element) {
		handleSubmittingTime(page, shouldAutoSubmit)
		// Navigated to Verification Code Page
	}).ElementR("label", "/Verification Code/i").MustHandle(func(e *rod.Element) {
		fmt.Println("You should have received a one time verification code in your email")
		fmt.Println("You can enter that code into the input in the browser")

		handleSubmittingTime(page, shouldAutoSubmit)

		// Incorrect User Credentials
	}).Element(".loginError").MustHandle(func(e *rod.Element) {
		fmt.Println("Log in unsuccessful, please enter credentials again")
		page.Browser().Close()
		main()
	}).MustDo()
}

func checkIfAutoSubmit() bool {
	fmt.Println("Would you like to auto-submit the timesheet? (Best for regular work week i.e. no PTO or holidays)")
	shouldAutoSubmit, _ := pterm.DefaultInteractiveConfirm.Show()
	if !shouldAutoSubmit {
		fmt.Println("")
		fmt.Println("After logging in, you will have 10 minutes to update and manually submit your timesheet")
		fmt.Println("")
	}
	return shouldAutoSubmit
}

func handleSubmittingTime(page *rod.Page, shouldAutoSubmit bool) {
	timeExpenseTabSelector := "body > div.desktop.container.forceStyle.oneOne.navexDesktopLayoutContainer.lafAppLayoutHost.forceAccess.tablet > div.viewport > section > div.none.navexStandardManager > div.slds-no-print.oneAppNavContainer > one-appnav > div > one-app-nav-bar > nav > div > one-app-nav-bar-item-root:nth-child(4) > a"
	timeExpenseTab := page.MustSearch(timeExpenseTabSelector)
	timeExpenseTab.MustClick()

	calendarBtnSelector := "#TimesheetHeader > div > span.clickable.left > a > span.linkText"
	calendarBtn := page.MustSearch(calendarBtnSelector)
	calendarBtn.MustClick()

	// Give them 10 mins to add custom times
	if shouldAutoSubmit {
		timesheetBodySelector := "div#TimesheetBody"
		timesheetBody := page.MustSearch(timesheetBodySelector)
		actualTimeCheckboxSelector := ".actualise-time"
		actualTimeCheckboxes := timesheetBody.MustElements(actualTimeCheckboxSelector)

		// Check all checkboxes
		for i := 0; i < len(actualTimeCheckboxes); i++ {
			checkbox := actualTimeCheckboxes[i]
			checkbox.MustClick()
		}

		convertActualTimeBtnSelector := ".actualise-selected-btn"
		convertActualTimeBtn := page.MustSearch(convertActualTimeBtnSelector)
		convertActualTimeBtn.MustClick()
		// Auto-submit
		submitAllBtnSelector := ".submit-all-btn"
		submitAllBtn := page.MustSearch(submitAllBtnSelector)
		submitAllBtn.MustClick()
		time.Sleep(time.Second * 5)
	} else {
		fmt.Println("You have 10 minutes to submit your timesheet before the browser window closes")
		fmt.Println("When you are finished submitting your timesheet you can close the browser window")
		time.Sleep(time.Minute * 10)
	}
}
