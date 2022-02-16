package server

import (
	model "TCP_Chat/model"
	My_cmd "TCP_Chat/model/command"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type client struct {
	conn   net.Conn
	name   string
	writer model.CommandWriter
}

type Chat_Server struct {
	listner net.Listener
	mutex   *sync.Mutex
	clients []*client
	count   uint64
}

func NewServer() *Chat_Server {
	return &Chat_Server{mutex: &sync.Mutex{}, count: 0}
}

//開始監聽
func (s *Chat_Server) Listen(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.listner = ln
	log.Printf("Server run on %v\n", ln.Addr())
	return nil
}

//廣播訊息
func (s *Chat_Server) Broadcast(cmd interface{}) error {
	for _, c := range s.clients {
		err := c.writer.Write(cmd)
		if err != nil {
			log.Printf("Broadcast to %v %v error : %v\n", c.name, c.conn.RemoteAddr().String(), err)
		}
	}
	return nil
}

//開始工作
func (s *Chat_Server) StartProcess() {
	for {
		conn, err := s.listner.Accept()
		if err != nil {
			log.Print(err)
		} else {
			cli := s.accept(conn) //加入client
			go s.process(cli)     //處理client
		}

	}
}

//加入client
func (s *Chat_Server) accept(conn net.Conn) *client {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.count += 1
	cli := &client{name: fmt.Sprintf("guest_%d", s.count),
		conn: conn, writer: *model.NewCmdWriter(conn)}

	s.clients = append(s.clients, cli)
	log.Printf("Accepting new connection [%v] from %v", cli.name, cli.conn.RemoteAddr().String())
	go s.Broadcast(My_cmd.BroadcastCommand{Name: "SYSTEM", Msg: fmt.Sprintf("%v join chat room", cli.name)})
	return cli
}

//登出client
func (s *Chat_Server) remove(cli *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, c := range s.clients {
		if c == cli {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}

	log.Printf("Closing connection [%v] from %v", cli.name, cli.conn.RemoteAddr().String())
	go s.Broadcast(My_cmd.BroadcastCommand{Name: "SYSTEM", Msg: fmt.Sprintf("%v leave chat room", cli.name)})
	cli.conn.Close()
}

//處理client
func (s *Chat_Server) process(cli *client) {
	cmdReader := model.NewCmdReader(cli.conn)

	defer s.remove(cli)

	for {
		cmd, err := cmdReader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		if cmd != nil {
			switch msg := cmd.(type) {
			case My_cmd.SendMsgCommand:
				go s.Broadcast(My_cmd.BroadcastCommand{Name: cli.name, Msg: msg.Msg})
			case My_cmd.SetNameCommand:
				old := cli.name
				cli.name = msg.Name
				go s.Broadcast(My_cmd.BroadcastCommand{Name: "SYSTEM", Msg: fmt.Sprintf("user [%v] change name to [%v]", old, cli.name)})
			}
		}
	}
}

//關閉
func (s *Chat_Server) Close() {
	s.listner.Close()
}
