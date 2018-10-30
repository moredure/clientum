package main

import (
	"fmt"
	"github.com/marcusolsson/tui-go"
	"github.com/mikefaraponov/chatum"
	"time"
)

const TimeFormat = "15:04:05"

func NewUI(user User, client chatum.Chatum_CommunicateClient) (tui.UI, error) {
	history := tui.NewVBox()
	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)
	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)
	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)
	input.OnSubmit(func(e *tui.Entry) {
		txt := e.Text()
		msg := &chatum.ClientSideEvent{
			Message: txt,
		}
		if err := client.Send(msg); err != nil {
			history.Append(tui.NewHBox(
				tui.NewLabel(time.Now().Format(TimeFormat)),
				tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", string(user)))),
				tui.NewLabel("Message not delivered!"),
				tui.NewSpacer(),
			))
		} else {
			history.Append(tui.NewHBox(
				tui.NewLabel(time.Now().Format(TimeFormat)),
				tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", string(user)))),
				tui.NewLabel(txt),
				tui.NewSpacer(),
			))
		}
		input.SetText("")
	})
	root := tui.NewHBox(chat)
	ui, err := tui.New(root)
	go func() {
		for {
			msg, err := client.Recv()
			if err != nil {
				history.Append(tui.NewHBox(
					tui.NewLabel(time.Now().Format(TimeFormat)),
					tui.NewPadder(1, 0, tui.NewLabel("$error$")),
					tui.NewLabel("Failed to receive message"),
					tui.NewSpacer(),
				))
			} else {
				history.Append(tui.NewHBox(
					tui.NewLabel(time.Now().Format(TimeFormat)),
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", msg.GetUsername()))),
					tui.NewLabel(msg.GetMessage()),
					tui.NewSpacer(),
				))
			}
			ui.Update(func() {})
		}
	}()
	return ui, err
}