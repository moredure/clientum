package main

import (
	"fmt"
	"github.com/marcusolsson/tui-go"
	"github.com/mikefaraponov/chatum"
	"time"
)

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
		if err := client.Send(NewMessage(e.Text())); err != nil {
			history.Append(tui.NewHBox(
				tui.NewLabel(time.Now().Format(TimeFormat)),
				tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf(UserLabelTemplate, string(user)))),
				tui.NewLabel(MessageNotDeliveredPrompt),
				tui.NewSpacer(),
			))
		}
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format(TimeFormat)),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf(UserLabelTemplate, string(user)))),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		input.SetText(CleanInput)
	})
	ui, err := tui.New(tui.NewHBox(chat))
	go func() {
		for {
			if msg, err := client.Recv(); err != nil {
				history.Append(tui.NewHBox(
					tui.NewLabel(time.Now().Format(TimeFormat)),
					tui.NewPadder(1, 0, tui.NewLabel(ErrorLabel)),
					tui.NewLabel(PleaseRestartPrompt),
					tui.NewSpacer(),
				))
				ui.Update(DoNothing)
				return
			} else if msg.GetType() == chatum.EventType_DEFAULT {
				history.Append(tui.NewHBox(
					tui.NewLabel(time.Now().Format(TimeFormat)),
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf(UserLabelTemplate, msg.GetUsername()))),
					tui.NewLabel(msg.GetMessage()),
					tui.NewSpacer(),
				))
				ui.Update(DoNothing)
			} else if msg.GetType() == chatum.EventType_PING {
				client.Send(NewPongMessage())
			}
		}
	}()
	return ui, err
}
