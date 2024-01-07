package main

import (
	"fmt"
	"github.com/joe-at-startupmedia/posix_mq"
	pmq_responder "github.com/joe-at-startupmedia/posix_mq/duplex/responder"
	"log"
	"time"
)

const maxSendTickNum = 10

func handleMessage(request []byte) (processed []byte, err error) {
	return []byte(fmt.Sprintf("I recieved request: %s\n", request)), nil
}

func main() {
	flags := posix_mq.O_RDWR | posix_mq.O_CREAT
	err := pmq_responder.New("posix_mq_example_duplex", posix_mq.POSIX_MQ_DIR, posix_mq.Ownership{}, flags)
	if err != nil {
		log.Fatal("Responder: could not initialize: ", err)
	}

	defer pmq_responder.Close()

	count := 0
	for {
		time.Sleep(1 * time.Second)
		count++
		err := pmq_responder.HandleRequest(handleMessage)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Responder: Sent a response")

		if count >= maxSendTickNum {
			break
		}
	}

}
