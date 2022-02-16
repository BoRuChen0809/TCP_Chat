package main

import (
	"TCP_Chat/client"
	"fmt"
	"log"

	"github.com/marcusolsson/tui-go"
)

func start_ui(cli *client.Chat_Client) {
	//對話紀錄
	history := tui.NewVBox()
	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)
	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)
	//使用者名稱
	name := tui.NewEntry()
	name.SetSizePolicy(tui.Expanding, tui.Maximum)
	name.SetFocused(true)
	nameBox := tui.NewHBox(tui.NewLabel("Set your name : "), name)
	nameBox.SetBorder(true)
	nameBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	name.OnSubmit(func(entry *tui.Entry) {
		if entry.Text() != "" {
			cli.SetName(entry.Text())
			entry.SetText("")
		}
	})
	//輸入框
	input := tui.NewEntry()
	input.SetFocused(false)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox := tui.NewHBox(tui.NewLabel("Msg : "), input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	input.OnSubmit(func(entry *tui.Entry) {
		if entry.Text() != "" {
			cli.SendMsg(entry.Text())
			entry.SetText("")
		}
	})

	chat := tui.NewVBox(historyBox, inputBox, nameBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	ui, err := tui.New(chat)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Ctrl+n", func() {
		name.SetFocused(true)
		input.SetFocused(false)
	})
	ui.SetKeybinding("Ctrl+t", func() {
		name.SetFocused(false)
		input.SetFocused(true)
	})

	go func() {
		for c := range cli.Msgs() {
			ui.Update(func() {
				history.Append(tui.NewHBox(
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>:", c.Name))),
					tui.NewLabel(c.Msg),
					tui.NewSpacer(),
				))
			})
		}
	}()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
