package main

import (
	"github.com/gcbest/timesheet-automation/util"
)

func main() {
	util.PrintWelcomeMessage()
	shouldAutoSubmit := util.CheckIfAutoSubmit()
	util.Start(shouldAutoSubmit)
}
