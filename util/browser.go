package util

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/pterm/pterm"
)

func logInUser(page *rod.Page, user User, shouldAutoSubmit bool) {
	// Navigated to Choose a Username Page
	page.Race().ElementR("h2#header", "/Choose a Username/i").MustHandle(func(e *rod.Element) {
		savedUserElement := page.MustElement("ul#idlist span")
		displayedUsername := savedUserElement.MustText()
		fmt.Println("displayedUsername: ", displayedUsername)
		if displayedUsername == user.username {
			savedUserElement.MustClick()
			page.MustElement("a#clear_link").MustClick()
			handleCredentialsPage(page, user, shouldAutoSubmit)
		} else {
			printErrMessage("Log in unsuccessful, please enter credentials again")
			page.Browser().Close()
			GetCredentialsAndLogin(shouldAutoSubmit)
		}

		// Navigated to Prefilled Login Page
	}).Element("div#idcard-container[style*=\"display: block;\"]").MustHandle(func(e *rod.Element) {
		page.MustElement("a#clear_link").MustClick()
		handleCredentialsPage(page, user, shouldAutoSubmit)

		// Navigated to Login Page
	}).Element("input#username[style*=\"display: block;\"]").MustHandle(func(e *rod.Element) {
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
	page.MustElement("input#username[style*=\"display: block;\"]").MustSelectAllText().MustInput("").MustInput(user.username)
	page.MustElement("input#password").MustSelectAllText().MustInput("").MustInput(user.password)
	rememberMeCheckBox := page.MustElement("input#rememberUn")
	if !rememberMeCheckBox.MustProperty("checked").Bool() {
		rememberMeCheckBox.MustClick()
	}
	page.MustElement("input#Login").MustClick()
}

func GetCredentialsAndLogin(shouldAutoSubmit bool) {
	user, _ := PromptForCredentials()
	spinnerInfo, _ := pterm.DefaultSpinner.Start("Launching browser...")

	url := launcher.New().
		UserDataDir("path").
		Headless(false).
		MustLaunch()

	page := rod.New().ControlURL(url).MustConnect().MustPage("https://inrhythm.my.salesforce.com/")
	spinnerInfo.Success()
	logInUser(page, user, shouldAutoSubmit)
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
		GetCredentialsAndLogin(shouldAutoSubmit)
	}).MustDo()
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
		fmt.Println("When you are finished submitting your timesheet you can press (ctrl + C) in the terminal to exit the program")
		time.Sleep(time.Minute * 10)
	}
}
