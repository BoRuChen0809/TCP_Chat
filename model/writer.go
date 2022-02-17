package model

import (
	"fmt"
	"io"

	My_cmd "TCP_Chat/model/command"
)

type CommandWriter struct {
	writer io.Writer
}

func NewCmdWriter(w io.Writer) *CommandWriter {
	return &CommandWriter{w}
}

func (cmdWriter *CommandWriter) Write(cmd interface{}) (err error) {
	switch msg := cmd.(type) {
	case My_cmd.SetNameCommand:
		_, err = cmdWriter.writer.Write([]byte(fmt.Sprintf("SetName %v\n", msg.Name)))
		return
	case My_cmd.SendMsgCommand:
		_, err = cmdWriter.writer.Write([]byte(fmt.Sprintf("Send %v\n", msg.Msg)))
		return
	case My_cmd.BroadcastCommand:
		_, err = cmdWriter.writer.Write([]byte(fmt.Sprintf("Broadcast %v %v\n", msg.Name, msg.Msg)))
		return
	case My_cmd.ChangeRoomCommand:
		_, err = cmdWriter.writer.Write([]byte(fmt.Sprintf("ChangeRoom %v\n", msg.ID)))
		return
	default:
		return fmt.Errorf("unknown msg")
	}

}
