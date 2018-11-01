package main

import (
	"fmt"
	"github.com/marcusolsson/tui-go"
	"github.com/mikefaraponov/chatum"
	"time"
)

type terminalUI struct {
	User
	chatum.Chatum_CommunicateClient
	tui.UI
	history    *tui.Box
	historyBox *tui.Box
	inputBox   *tui.Box
	chat       *tui.Box
	input      *tui.Entry
}

func (ui *terminalUI) addHistory() *terminalUI {
	ui.history = tui.NewVBox()
	historyScroll := tui.NewScrollArea(ui.history)
	historyScroll.SetAutoscrollToBottom(true)
	ui.historyBox = tui.NewVBox(historyScroll)
	ui.historyBox.SetBorder(true)
	return ui
}

func (ui *terminalUI) addInput() *terminalUI {
	ui.input = tui.NewEntry()
	ui.input.SetFocused(true)
	ui.input.SetSizePolicy(tui.Expanding, tui.Maximum)
	ui.input.OnSubmit(ui.onSubmit)
	ui.inputBox = tui.NewHBox(ui.input)
	ui.inputBox.SetBorder(true)
	ui.inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	return ui
}

func (ui *terminalUI) addChat() *terminalUI {
	ui.chat = tui.NewVBox(ui.historyBox, ui.inputBox)
	ui.chat.SetSizePolicy(tui.Expanding, tui.Expanding)
	return ui
}

func (ui *terminalUI) onSubmit(e *tui.Entry) {
	if err := ui.Send(NewMessage(e.Text())); err != nil {
		ui.history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format(TimeFormat)),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf(UserLabelTemplate, string(ui.User)))),
			tui.NewLabel(MessageNotDeliveredPrompt),
			tui.NewSpacer(),
		))
	}
	ui.history.Append(tui.NewHBox(
		tui.NewLabel(time.Now().Format(TimeFormat)),
		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf(UserLabelTemplate, string(ui.User)))),
		tui.NewLabel(e.Text()),
		tui.NewSpacer(),
	))
	ui.input.SetText(CleanInput)
}

func (ui *terminalUI) buildUI() (terminal tui.UI, err error) {
	terminal, err = tui.New(tui.NewHBox(ui.chat))
	if err != nil {
		return
	}
	ui.UI = terminal
	go ui.onNewMessage()
	return
}

func (ui *terminalUI) onNewMessage() {
	for {
		if msg, err := ui.Recv(); err != nil {
			ui.history.Append(tui.NewHBox(
				tui.NewLabel(time.Now().Format(TimeFormat)),
				tui.NewPadder(1, 0, tui.NewLabel(ErrorLabel)),
				tui.NewLabel(PleaseRestartPrompt),
				tui.NewSpacer(),
			))
			ui.Update(DoNothing)
			return
		} else if msg.GetType() == chatum.EventType_DEFAULT {
			ui.history.Append(tui.NewHBox(
				tui.NewLabel(time.Now().Format(TimeFormat)),
				tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf(UserLabelTemplate, msg.GetUsername()))),
				tui.NewLabel(msg.GetMessage()),
				tui.NewSpacer(),
			))
			ui.Update(DoNothing)
		} else if msg.GetType() == chatum.EventType_PING {
			ui.Send(NewPongMessage())
		}
	}
}

func newTerminalUI(user User, client chatum.Chatum_CommunicateClient) *terminalUI {
	return &terminalUI{
		User:                     user,
		Chatum_CommunicateClient: client,
	}
}

func NewUI(user User, client chatum.Chatum_CommunicateClient) (tui.UI, error) {
	return newTerminalUI(user, client).
		addHistory().
		addInput().
		addChat().
		buildUI()
}
