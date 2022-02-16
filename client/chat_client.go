package client

import (
	"TCP_Chat/model"
	My_cmd "TCP_Chat/model/command"
	"io"
	"log"
	"net"
)

type Chat_Client struct {
	conn      net.Conn
	cmdWriter *model.CommandWriter
	cmdReader *model.CommandReader
	name      string
	msgs      chan My_cmd.BroadcastCommand
}

func NewClient() *Chat_Client {
	return &Chat_Client{msgs: make(chan My_cmd.BroadcastCommand)}
}

//連線
func (cli *Chat_Client) Dial(address string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	cli.conn = conn

	cli.cmdReader = model.NewCmdReader(conn)
	cli.cmdWriter = model.NewCmdWriter(conn)

	return nil
}

//設定名字
func (cli *Chat_Client) SetName(name string) error {
	err := cli.cmdWriter.Write(My_cmd.SetNameCommand{Name: name})
	if err == nil {
		cli.name = name
		return nil
	}
	return err
}

//傳訊息
func (cli *Chat_Client) SendMsg(msg string) error {
	return cli.cmdWriter.Write(My_cmd.SendMsgCommand{Msg: msg})
}

//接收server廣播
func (cli *Chat_Client) StartReceive() {
	for {
		cmd, err := cli.cmdReader.Read()
		//if server無回應
		if err == io.EOF {
			break
		} else if err != nil { //其他錯誤
			log.Println(err)
			break
		}

		if cmd != nil {
			switch msg := cmd.(type) {
			case My_cmd.BroadcastCommand:
				cli.msgs <- msg
			default:
				log.Printf("unknown type : %v", msg)
			}
		}
	}
}

//關閉連線
func (cli *Chat_Client) Close() {
	cli.conn.Close()
}

func (cli *Chat_Client) Msgs() chan My_cmd.BroadcastCommand {
	return cli.msgs
}
