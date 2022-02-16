package server

type My_Server interface {
	Listen(address string) error     //監聽
	Broadcast(cmd interface{}) error //廣播訊息
	StartProcess()                   //開始工作
	Close()                          //關閉
}
