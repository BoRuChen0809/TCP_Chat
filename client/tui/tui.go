package main

import (
	"TCP_Chat/client"
	"fmt"
	"log"

	"github.com/marcusolsson/tui-go"
)

func start_ui(cli *client.Chat_Client) {
	//使用者名稱
	name := tui.NewEntry()
	name.SetSizePolicy(tui.Expanding, tui.Maximum)
	name.SetFocused(true)
	nameBox := tui.NewVBox(tui.NewLabel("Input your nickname : "), name)
	nameBox.SetBorder(true)
	nameBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	name.OnSubmit(func(entry *tui.Entry) {
		if entry.Text() != "" {
			cli.SetName(entry.Text())
		}
	})

	//切換聊天室
	room := tui.NewEntry()
	room.SetSizePolicy(tui.Expanding, tui.Maximum)
	room.SetFocused(false)
	roomBox := tui.NewVBox(tui.NewLabel("Input Room ID : "), room)
	roomBox.SetBorder(true)
	roomBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	room.OnSubmit(func(entry *tui.Entry) {
		if entry.Text() != "" {
			cli.ChangeRoom(entry.Text())
		}
	})
	//左邊區域
	sidebar := tui.NewVBox(
		nameBox,
		tui.NewLabel(""),
		roomBox,
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	//對話紀錄
	history := tui.NewVBox()
	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)
	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	//訊息輸入框
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
	//右邊區域
	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	whole_frame := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(whole_frame)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Ctrl+n", func() {
		room.SetFocused(false)
		name.SetFocused(true)
		input.SetFocused(false)
	})
	ui.SetKeybinding("Ctrl+t", func() {
		room.SetFocused(false)
		name.SetFocused(false)
		input.SetFocused(true)
	})
	ui.SetKeybinding("Ctrl+r", func() {
		room.SetFocused(true)
		name.SetFocused(false)
		input.SetFocused(false)
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
