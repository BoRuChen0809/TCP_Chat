package model_test

import (
	model "TCP_Chat/model"
	My_cmd "TCP_Chat/model/command"
	"bytes"
	"reflect"
	"testing"
)

func Test_Writer(t *testing.T) {
	type test_case struct {
		cmd    interface{}
		expect string
	}

	tt := []test_case{
		{cmd: My_cmd.SetNameCommand{Name: "BoBo"}, expect: "SetName BoBo\n"},
		{cmd: My_cmd.SendMsgCommand{Msg: "Test Message"}, expect: "Send Test Message\n"},
		{cmd: My_cmd.BroadcastCommand{Name: "BoBo", Msg: "Test Message"}, expect: "Broadcast BoBo Test Message\n"},
	}

	buf := new(bytes.Buffer)
	for _, tc := range tt {
		buf.Reset()
		w := model.NewCmdWriter(buf)
		if w.Write(tc.cmd) != nil {
			t.Errorf("write error")
		}
		if !reflect.DeepEqual(buf.String(), tc.expect) {
			t.Errorf("output and expect are not same")
		}
	}

}
