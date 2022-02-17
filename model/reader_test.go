package model_test

import (
	"TCP_Chat/model"
	My_cmd "TCP_Chat/model/command"
	"reflect"
	"strings"
	"testing"
)

func Test_Reader(t *testing.T) {
	type test_case struct {
		cmd_str string
		expect  interface{}
	}

	tt := []test_case{
		{expect: My_cmd.BroadcastCommand{Name: "BoBo", Msg: "Test Message"},
			cmd_str: "Broadcast BoBo Test Message\n"},
		{expect: My_cmd.SetNameCommand{Name: "BoBo"}, cmd_str: "SetName BoBo\n"},
		{expect: My_cmd.SendMsgCommand{Msg: "Test Message"}, cmd_str: "Send Test Message\n"},
		{expect: My_cmd.ChangeRoomCommand{ID: "Default"}, cmd_str: "ChangeRoom Default\n"},
	}

	for _, tc := range tt {
		r := model.NewCmdReader(strings.NewReader(tc.cmd_str))
		cmd, err := r.Read()
		t.Log(cmd)
		if err != nil {
			t.Errorf("read error : %v\n", err)
		}
		if !reflect.DeepEqual(cmd, tc.expect) {
			t.Errorf("output and expect are not same")
		}
	}
}
