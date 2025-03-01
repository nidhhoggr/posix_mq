package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nidhhoggr/posix_mq"
)

const maxSendTickNum = 10

func main() {
	send_c := make(chan int)
	go sender(send_c)
	//wait for the sender to create the posix_mq files
	time.Sleep(1 * time.Second)
	recv_c := make(chan int)
	go receiver(recv_c)
	<-recv_c
	<-send_c
	//gives time for deferred functions to complete
	time.Sleep(2 * time.Second)
}

func sender(c chan int) {
	mq, err := posix_mq.NewMessageQueue(&posix_mq.QueueConfig{
		Name:  "posix_mq_example_simple",
		Flags: posix_mq.O_WRONLY | posix_mq.O_CREAT,
		Mode:  0666,
	})
	defer func() {
		fmt.Println("Sender: finished")
		c <- 0
	}()
	if err != nil {
		fmt.Printf("Sender: error initializing: %s", err)
		c <- 1
		return
	}

	count := 0
	for {
		count++
		err = mq.Send([]byte(fmt.Sprintf("Hello, World : %d\n", count)), 0)
		if err != nil {
			fmt.Printf("Sender: error sending message: %s\n", err)
			continue
		}

		fmt.Println("Sender: Sent a new message")

		if count >= maxSendTickNum {
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func receiver(c chan int) {
	mq, err := posix_mq.NewMessageQueue(&posix_mq.QueueConfig{
		Name:  "posix_mq_example_simple",
		Flags: posix_mq.O_RDONLY,
		Mode:  0666,
	})
	defer func() {
		closeQueue(mq)
		fmt.Println("Receiver: finished")
		c <- 0
	}()
	if err != nil {
		fmt.Printf("Receiver: error initializing %s", err)
		c <- 1
		return
	}

	fmt.Println("Receiver: Start receiving messages")

	count := 0
	for {
		count++

		msg, _, err := mq.Receive()
		if err != nil {
			fmt.Printf("Receiver: error getting message: %s\n", err)
			continue
		}

		fmt.Printf("Receiver: got new message: %s\n", string(msg))

		if count >= maxSendTickNum {
			break
		}
	}
}

func closeQueue(mq *posix_mq.MessageQueue) {
	if err := mq.Unlink(); err != nil {
		log.Printf("closeQueue error: %s", err)
	}
}
