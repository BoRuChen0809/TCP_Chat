# TCP_Chat
用socket練習寫tcp聊天室

主要分四部分:
- [TCP_Chat](#tcp_chat)
	- [共用model](#共用model)
	- [Client](#client)
	- [Server](#server)
	- [TUI](#tui)


## 共用model

* ### Command(傳送的訊息)

  server端要用到的功能有:
  
  - 廣播訊息
  
    ```go
    // server use
    type BroadcastCommand struct {
    	Name string
    	Msg  string
    }
    ```
  
  client端要用到的功能有:
  
  - 傳送訊息:
  
    ```go
    // client use
    type SendMsgCommand struct {
    	Msg string
    }
    ```
  
  - 修改名字:
  
    ```go
    // client use
    type SetNameCommand struct {
    	Name string
    }
    ```

  - 變更聊天室
	```go
	//client use
	type ChangeRoomCommand struct {
		ID string
	}
	```
  
    [command程式碼](https://github.com/BoRuChen0809/TCP_Chat/blob/main/model/command/command.go)
  
* ### CmdWriter

  功用:將command轉成指定格式string再將其以byte array格式輸出

  	```go
  	type CommandWriter struct {
  		writer io.Writer
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
  
	```

  [writert程式碼](https://github.com/BoRuChen0809/TCP_Chat/blob/main/model/writer.go)

* ### CmdReader

  功用:讀取CmdWriter輸出的byte array並轉為command物件

  ```go
	type CommandReader struct {
		reader *bufio.Reader
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
  ```

  [reader程式碼](https://github.com/BoRuChen0809/TCP_Chat/blob/main/model/reader.go)

  ------

## Client

- ### 定義client有的動作

  ```go
  	type My_Client interface {
		Dial(address string) error       //連線
		SetName(name string) error       //設定名字
		SendMsg(msg string) error        //傳訊息
		ChangeRoom(room_id string) error //切換聊天室
		StartReceive()                   //接收server廣播
		Close()                          //關閉連線
	}
  ```

- ### 實現Chat_Client

  ```go
	type Chat_Client struct {
		conn      net.Conn
		cmdWriter *model.CommandWriter
		cmdReader *model.CommandReader
		name      string
		room_id   string
		msgs      chan My_cmd.BroadcastCommand
	}
  ```

  - 連線:

    ```go
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
    ```

  - 設定名字:

    ```go
    func (cli *Chat_Client) SetName(name string) error {
    	err := cli.cmdWriter.Write(My_cmd.SetNameCommand{Name: name})
    	if err == nil {
    		cli.name = name
    		return nil
    	}
    	return err
    }
    ```

  - 傳訊息:

    ```go
    func (cli *Chat_Client) SendMsg(msg string) error {
    	return cli.cmdWriter.Write(My_cmd.SendMsgCommand{Msg: msg})
    }
    ```

  - 切換聊天室
	```go
	func (cli *Chat_Client) ChangeRoom(room_id string) error {
		err := cli.cmdWriter.Write(My_cmd.ChangeRoomCommand{ID: room_id})
		if err == nil {
			cli.room_id = room_id
			return nil
		}
		return err
	}
	```

  - 接收server訊息:

    ```go
    func (cli *Chat_Client) StartReceive() {
    	for {
    		cmd, err := cli.cmdReader.Read()
    		
    		if err == io.EOF {		//if server無回應
    			break
    		} else if err != nil { 	//其他錯誤
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
    ```

  - 關閉連線:

    ```go
    func (cli *Chat_Client) Close() {
    	cli.conn.Close()
    }
    ```

  - 取出要顯示的訊息:

    ```go
    func (cli *Chat_Client) Msgs() chan My_cmd.BroadcastCommand {
    	return cli.msgs
    }
    ```

    [client程式碼](https://github.com/BoRuChen0809/TCP_Chat/tree/main/client)

    ------

## Server

- ### 定義server有的動作

  ```go
  type My_Server interface {
  	Listen(address string) error     //監聽port
  	Broadcast(cmd interface{}) error //廣播訊息
  	StartProcess()                   //開始工作
  	Close()                          //關閉
  }
  ```

- ### 實現Chat_server

  ```go
  	type Chat_Server struct {
		listner net.Listener
		mutex   *sync.Mutex
		count     uint64
		chat_room map[string][]*client
	}
  
  	//用來存client資訊及傳訊息給client
  	type client struct {
		conn    net.Conn
		name    string
		writer  model.CommandWriter
		room_id string
	}
  ```

  - 監聽port:

    ```go
    func (s *Chat_Server) Listen(address string) error {
    	ln, err := net.Listen("tcp", address)
    	if err != nil {
    		return err
    	}
    
    	s.listner = ln
    	log.Printf("Server run on %v\n", ln.Addr())
    	return nil
    }
    ```

  - 廣播訊息:

    ```go
    func (s *Chat_Server) Broadcast(cmd interface{}) error {
    	for _, c := range s.clients {
    		err := c.writer.Write(cmd)
    		if err != nil {
    			log.Printf("Broadcast to %v %v error : %v\n", 
                           c.name, c.conn.RemoteAddr().String(), err)
    		}
    	}
    	return nil
    }
    ```

  - 開始工作(處理client的連線、訊息、離線以及加入聊天室):

    ```go
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
    ```

    server處理client的流程主要是

    1. accept client 連線
    2. handle client Message
    3. handle client 離線

    分別對應以下三個function:

    - 加入client:

      ```go
		func (s *Chat_Server) accept(conn net.Conn) *client {
			s.mutex.Lock()
			defer s.mutex.Unlock()

			s.count += 1
			cli := &client{name: fmt.Sprintf("guest_%d", s.count),
				conn: conn, writer: *model.NewCmdWriter(conn), room_id: "DEFAULT"}

			log.Printf("Accepting new connection [%v] from %v", cli.name, cli.conn.RemoteAddr().String())

			return cli
		}
      ```

    - 處理client訊息:

      ```go
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
      ```
	  //client的聊天室相關功能有三個function
	  ```go
	  	//加入聊天室
		func (s *Chat_Server) joinRoom(room_id string, cli *client) {
			s.mutex.Lock()

			clis := s.chat_room[room_id]
			clis = append(clis, cli)
			s.chat_room[room_id] = clis
			cli.room_id = room_id

			log.Printf("[%v] join chat room [%v]\n", cli.name, cli.room_id)

			s.mutex.Unlock()
		}
		//離開聊天室
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
		//當聊天室人數為0時觸發的移除聊天室
		func (s *Chat_Server) deleteRoom(room_id string) {
			s.mutex.Lock()
			defer s.mutex.Unlock()

			delete(s.chat_room, room_id)
		}
	  ```

    - client離線:

      ```go
      	func (s *Chat_Server) remove(cli *client) {
			log.Printf("Closing connection [%v] from %v", cli.name, cli.conn.RemoteAddr().String())
			go s.Broadcast(My_cmd.BroadcastCommand{Name: "SYSTEM", Msg: fmt.Sprintf("%v leave chat room [%v]", cli.name, cli.room_id)}, cli.room_id)
			cli.conn.Close()
		}
      ```

  - server關閉:

    ```go
    func (s *Chat_Server) Close() {
    	s.listner.Close()
    }
    ```

    [server程式碼](https://github.com/BoRuChen0809/TCP_Chat/tree/main/server)

    ------

## TUI

參考範例:https://github.com/yuuki0xff/tui-go/tree/master/example/chat

[tui程式碼](https://github.com/BoRuChen0809/TCP_Chat/blob/main/client/tui/tui.go)

