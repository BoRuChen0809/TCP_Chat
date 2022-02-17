package client

type My_Client interface {
	Dial(address string) error       //連線
	SetName(name string) error       //設定名字
	SendMsg(msg string) error        //傳訊息
	ChangeRoom(room_id string) error //切換聊天室
	StartReceive()                   //接收server廣播
	Close()                          //關閉連線
}
