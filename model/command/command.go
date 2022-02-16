package My_cmd

// client use
type SetNameCommand struct {
	Name string
}

type SendMsgCommand struct {
	Msg string
}

// server use
type BroadcastCommand struct {
	Name string
	Msg  string
}
