package main

import "TCP_Chat/server"

func main() {
	s := server.NewServer()
	s.Listen(":8080")
	s.StartProcess()
}
