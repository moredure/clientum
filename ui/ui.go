package ui

import (
	"fmt"
	"github.com/marcusolsson/tui-go"
	"github.com/mikefaraponov/chatum"
	"github.com/mikefaraponov/clientum/common"
	"time"
)

type terminalUI struct {
	chatum.Chatum_CommunicateClient
	tui.UI
	user          string
	historyScroll *tui.ScrollArea
	history       *tui.Box
	historyBox    *tui.Box
	inputBox      *tui.Box
	chat          *tui.Box
	input         *tui.Entry
}

func (ui *terminalUI) addHistory() *terminalUI {
	ui.history = tui.NewVBox()
	ui.historyScroll = tui.NewScrollArea(ui.history)
	ui.historyScroll.SetAutoscrollToBottom(true)
	ui.historyBox = tui.NewVBox(ui.historyScroll)
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
	if err := ui.Send(common.NewMessage(e.Text())); err != nil {
		ui.appendToHistory(common.MessageNotDeliveredPrompt)
	}
	ui.appendToHistory(e.Text())
	ui.input.SetText(common.CleanInput)
}

func (ui *terminalUI) getUserLabel() *tui.Label {
	return tui.NewLabel(fmt.Sprintf(common.UserLabelTemplate, ui.user))
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
			ui.appendNetworkError()
			ui.Update(common.DoNothing)
			return
		} else if msg.GetType() == chatum.EventType_DEFAULT {
			ui.appendNewMessage(msg)
			ui.Update(common.DoNothing)
		} else if msg.GetType() == chatum.EventType_PING {
			ui.Send(common.NewPongMessage())
		}
	}
}

func (ui *terminalUI) appendToHistory(msg string) {
	ui.history.Append(tui.NewHBox(
		tui.NewLabel(time.Now().Format(common.TimeFormat)),
		tui.NewPadder(1, 0, ui.getUserLabel()),
		tui.NewLabel(msg),
		tui.NewSpacer(),
	))
}

func (ui *terminalUI) appendNewMessage(msg *chatum.ServerSideEvent) {
	ui.history.Append(tui.NewHBox(
		tui.NewLabel(time.Now().Format(common.TimeFormat)),
		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf(common.UserLabelTemplate, msg.GetUsername()))),
		tui.NewLabel(msg.GetMessage()),
		tui.NewSpacer(),
	))
}

func (ui *terminalUI) appendNetworkError() {
	ui.history.Append(tui.NewHBox(
		tui.NewLabel(time.Now().Format(common.TimeFormat)),
		tui.NewPadder(1, 0, tui.NewLabel(common.ErrorLabel)),
		tui.NewLabel(common.PleaseRestartPrompt),
		tui.NewSpacer(),
	))
}

func newTerminalUI(user string, client chatum.Chatum_CommunicateClient) *terminalUI {
	return &terminalUI{
		user:                     user,
		Chatum_CommunicateClient: client,
	}
}

func NewUI(env *common.Environment, client chatum.Chatum_CommunicateClient) (tui.UI, error) {
	return newTerminalUI(env.User, client).
		addHistory().
		addInput().
		addChat().
		buildUI()
}
