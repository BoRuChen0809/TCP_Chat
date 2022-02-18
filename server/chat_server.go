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
	conn    net.Conn
	name    string
	writer  model.CommandWriter
	room_id string
}

type Chat_Server struct {
	listner net.Listener
	mutex   *sync.Mutex
	count     uint64
	chat_room map[string][]*client
}

func NewServer() *Chat_Server {
	return &Chat_Server{mutex: &sync.Mutex{}, count: 0, chat_room: make(map[string][]*client)}
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
func (s *Chat_Server) Broadcast(cmd interface{}, room_id string) error {
	for _, c := range s.chat_room[room_id] {
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
			s.joinRoom(cli.room_id, cli)
			go s.Broadcast(My_cmd.BroadcastCommand{Name: "SYSTEM", Msg: fmt.Sprintf("%v join chat room [%v]", cli.name, "DEFAULT")}, "DEFAULT")
			go s.process(cli) //處理client
		}
	}
}

//加入client
func (s *Chat_Server) accept(conn net.Conn) *client {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.count += 1
	cli := &client{name: fmt.Sprintf("guest_%d", s.count),
		conn: conn, writer: *model.NewCmdWriter(conn), room_id: "DEFAULT"}

	log.Printf("Accepting new connection [%v] from %v", cli.name, cli.conn.RemoteAddr().String())

	return cli
}

//登出client
func (s *Chat_Server) remove(cli *client) {
	log.Printf("Closing connection [%v] from %v", cli.name, cli.conn.RemoteAddr().String())
	go s.Broadcast(My_cmd.BroadcastCommand{Name: "SYSTEM", Msg: fmt.Sprintf("%v leave chat room [%v]", cli.name, cli.room_id)}, cli.room_id)
	cli.conn.Close()
}

//處理client
func (s *Chat_Server) process(cli *client) {
	cmdReader := model.NewCmdReader(cli.conn)
	defer s.leaveRoom(cli)
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
				go s.Broadcast(My_cmd.BroadcastCommand{Name: cli.name, Msg: msg.Msg}, cli.room_id)
			case My_cmd.SetNameCommand:
				old := cli.name
				cli.name = msg.Name
				go s.Broadcast(My_cmd.BroadcastCommand{Name: "SYSTEM", Msg: fmt.Sprintf("user [%v] change name to [%v]", old, cli.name)}, cli.room_id)
			case My_cmd.ChangeRoomCommand:
				pre_room := cli.room_id
				s.leaveRoom(cli)
				go s.Broadcast(My_cmd.BroadcastCommand{Name: "SYSTEM", Msg: fmt.Sprintf("%v leave chat room [%v]", cli.name, pre_room)}, pre_room)
				if len(s.chat_room[pre_room]) == 0 {
					s.deleteRoom(pre_room)
				}
				cli.room_id = msg.ID
				go s.Broadcast(My_cmd.BroadcastCommand{Name: "SYSTEM", Msg: fmt.Sprintf("%v join chat room [%v]", cli.name, msg.ID)}, msg.ID)
				s.joinRoom(msg.ID, cli)

			}
		}
	}
}

func (s *Chat_Server) joinRoom(room_id string, cli *client) {
	s.mutex.Lock()

	clis := s.chat_room[room_id]
	clis = append(clis, cli)
	s.chat_room[room_id] = clis
	cli.room_id = room_id

	log.Printf("[%v] join chat room [%v]\n", cli.name, cli.room_id)

	s.mutex.Unlock()
}

func (s *Chat_Server) leaveRoom(cli *client) {
	s.mutex.Lock()

	clis := s.chat_room[cli.room_id]
	for i, c := range clis {
		if c == cli {
			clis = append(clis[:i], clis[i+1:]...)
		}
	}
	s.chat_room[cli.room_id] = clis

	log.Printf("[%v] leave chat room [%v]\n", cli.name, cli.room_id)
	cli.room_id = ""

	s.mutex.Unlock()
}

func (s *Chat_Server) deleteRoom(room_id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.chat_room, room_id)
}

//關閉
func (s *Chat_Server) Close() {
	s.listner.Close()
}
