package model

import (
	My_cmd "TCP_Chat/model/command"
	"bufio"
	"fmt"
	"io"
)

type CommandReader struct {
	reader *bufio.Reader
}

func NewCmdReader(r io.Reader) *CommandReader {
	return &CommandReader{reader: bufio.NewReader(r)}
}

func (r *CommandReader) Read() (interface{}, error) {
	cmd, err := r.reader.ReadString(' ')
	if err != nil {
		return nil, err
	}

	switch cmd {
	case "ChangeRoom ":
		msg, err := r.reader.ReadString('\n')
		return My_cmd.ChangeRoomCommand{ID: msg[:len(msg)-1]}, err
	case "Send ":
		msg, err := r.reader.ReadString('\n')
		return My_cmd.SendMsgCommand{Msg: msg[:len(msg)-1]}, err
	case "SetName ":
		name, err := r.reader.ReadString('\n')
		return My_cmd.SetNameCommand{Name: name[:len(name)-1]}, err
	case "Broadcast ":
		name, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}
		msg, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return My_cmd.BroadcastCommand{Name: name[:len(name)-1], Msg: msg[:len(msg)-1]}, nil
	default:
		return nil, fmt.Errorf("nuknown msg")
	}

}
