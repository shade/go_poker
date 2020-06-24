/*package main

import (
	"github.com/rivo/tview"
)

var app *tview.Application

func login() {
	flex := tview.NewFlex()

	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Login", nil).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)

	flex.AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
	flex.AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Welcome to GoPoker"), 0, 1, false).
			AddItem(form, 0, 3, true).
			AddItem(tview.NewBox().SetBorder(false), 0, 1, false),
		0, 2, true)
	flex.AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

func signup() {
	flex := tview.NewFlex()

	form := tview.NewForm().
		AddInputField("First name", "", 20, nil, nil).
		AddInputField("Last name", "", 20, nil, nil).
		AddInputField("Username", "", 20, nil, nil).
		AddCheckbox("Age 21+", false, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Save", nil).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)

	flex.AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
	flex.AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Welcome to GoPoker"), 0, 1, false).
			AddItem(form, 0, 3, true),
		0, 2, true)
	flex.AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

func main() {
	app = tview.NewApplication()

	flex := tview.NewFlex()

	flex.AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
	flex.AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText("Welcome to GoPoker"), 0, 1, false).
			AddItem(tview.NewList().
				AddItem("Log in", "", 'a', login).
				AddItem("Signup", "", 'b', signup).
				AddItem("Quit", "", 'd', nil), 0, 3, true).
			AddItem(tview.NewBox().SetBorder(false), 0, 1, false),
		0, 2, true)
	flex.AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
*/