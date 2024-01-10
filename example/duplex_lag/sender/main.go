package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joe-at-startupmedia/posix_mq"
	pmq_sender "github.com/joe-at-startupmedia/posix_mq/duplex/sender"
)

const maxSendTickNum = 10

func main() {
	err := pmq_sender.New("posix_mq_example_duplex", posix_mq.POSIX_MQ_DIR, posix_mq.Ownership{})
	if err != nil {
		log.Fatal("Sender: could not initialize: ", err)
	}
	defer func() {
		pmq_sender.Close()
		fmt.Println("Sender: finished and closed")
	}()
	count := 0
	c := make(chan pmqResponse)
	for {
		count++
		request := fmt.Sprintf("Hello, World : %d\n", count)
		go requestResponse(request, c)

		if count >= maxSendTickNum {
			break
		}

		time.Sleep(1 * time.Second)
	}

	result := make([]pmqResponse, maxSendTickNum)
	for i := range result {
		result[i] = <-c
		if result[i].status {
			fmt.Println(result[i].response)
		} else {
			fmt.Printf("Sender: Got error: %s \n", result[i].response)
		}
	}
}

func requestResponse(msg string, c chan pmqResponse) {
	err := pmq_sender.Send([]byte(msg), 0)
	if err != nil {
		c <- pmqResponse{fmt.Sprintf("%s", err), false}
		return
	}
	fmt.Printf("Sender: sent a new request: %s", msg)

	resp, _, err := pmq_sender.WaitForResponse(time.Second)

	if err != nil {
		c <- pmqResponse{fmt.Sprintf("%s", err), false}
		return
	}

	c <- pmqResponse{fmt.Sprintf("Sender: got a response: %s\n", resp), true}
}

type pmqResponse struct {
	response string
	status   bool
}
