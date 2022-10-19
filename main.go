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

func promptCredentials() (User, error) {
	reader := bufio.NewReader(os.Stdin)

	// TODO: add tests here
	username, err := getInput("Enter username: ", false, reader)
	if err != nil {
		fmt.Println("Please enter valid username")
		promptCredentials()
	}
	password, err := getInput("Enter password: ", true, reader)

	// TODO: only prompt for password again
	if err != nil {
		fmt.Println("Please enter valid password")
		promptCredentials()
	}

	newUser := createUser(username, password)
	// TODO: remove
	fmt.Println("user", newUser)
	return *newUser, nil
}

func promptVerificationCode() string {
	fmt.Println("You should have received a one time verification code in your email")
	reader := bufio.NewReader(os.Stdin)
	code, err := getInput("Enter one time code: ", false, reader)

	if err != nil {
		fmt.Println("Please enter a valid code")
		promptVerificationCode()
	}

	return code
}

func createUser(username string, password string) *User {
	newUser := User{
		username: username,
		password: password,
	}

	fmt.Println("Login credentials added")
	return &newUser
}

func logInUser(page *rod.Page, user User) {
	page.MustElement("input[name=\"username\"").MustInput(user.username)
	page.MustElement("input[name=\"pw\"]").MustInput(user.password)
	// rememberMeCheckBox := page.Locator("input#rememberUn")
	// rememberMeCheckBox.Check()
	page.MustElementR("input", "/Log In/i").MustClick()
}

func main() {
	user, _ := promptCredentials()

	url := launcher.New().
		UserDataDir("path").
		Headless(false).
		MustLaunch()

	page := rod.New().ControlURL(url).MustConnect().MustPage("https://inrhythm.my.salesforce.com/")
	logInUser(page, user)

	// TODO: message if login successful or not
	// fmt.Println(page.MustElement(".slds-global-header__logo"))
	// page.Race().Element("input[name=\"emc\"]").MustHandle(func(e *rod.Element) {
	// 	verificationCode := promptVerificationCode()

	// 	e.MustInput(verificationCode)
	// 	page.MustElementR("input", "/Verify/i").MustClick()

	// }).Element("body > div.desktop.container.forceStyle.oneOne.navexDesktopLayoutContainer.lafAppLayoutHost.forceAccess.tablet > div.viewport > section > div.none.navexStandardManager > div.slds-no-print.oneAppNavContainer > one-appnav > div > one-app-nav-bar > nav > div > one-app-nav-bar-item-root:nth-child(4) > a > span").MustHandle(func(e *rod.Element) {
	// 	e.MustClick()
	// }).Element("ERROR-CASE").MustHandle(func(e *rod.Element) {
	// 	// print message to user for error case
	// 	fmt.Println("Log in unsuccessful, please enter credentials again")
	// 	page.Browser().Close()
	// 	main()
	// }).MustDo()
	/********** Verification Code Submit ********** Conditional
	verificationCode := promptVerificationCode()

	verificationCodeInput, _ := page.Locator("input[name=\"emc\"]")
	verificationCodeInput.Fill(verificationCode)
	verificationCodeSubmitBtn, err := page.Locator("input:has-text(\"Verify\")")
	if err != nil {
		log.Fatalf("could not find code submit button: %v", err)
	}
	verificationCodeSubmitBtn.Click()
	// verificationCodeInput.Press("Enter")
	***************************************************/
	timeExpenseTab := page.MustSearch("one-app-nav-bar-item-root:nth-child(4) > a > span")
	timeExpenseTab.MustClick()

	calendarBtn := page.MustSearch("#TimesheetHeader > div > span.clickable.left > a > span.linkText")
	calendarBtn.MustClick()

	x := page.MustSearch("div#TimesheetBody")
	y := x.MustElements(".actualise-time")

	for i := 0; i < len(y); i++ {
		checkbox := y[i]
		checkbox.MustClick()
	}

	convertActualTimeBtn := page.MustSearch(".actualise-selected-btn")
	convertActualTimeBtn.MustClick()

	// submitAllBtn := page.MustSearch(".submit-all-btn")
	// submitAllBtn.MustClick()

	// TODO: determine the amount of time to sleep
	time.Sleep(time.Hour)
}
