package main

import (
	"TCP_Chat/client"
	"log"
)

func main() {

	cli := client.NewClient()

	err := cli.Dial("localhost:8080")
	if err != nil {
		log.Panic(err)
	}

	defer cli.Close()
	go cli.StartReceive()

	start_ui(cli)

}
