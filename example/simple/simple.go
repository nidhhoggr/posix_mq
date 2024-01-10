package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joe-at-startupmedia/posix_mq"
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
}

func sender(c chan int) {
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mq, err := posix_mq.NewMessageQueue("/posix_mq_example", oflag, 0666, nil)
	if err != nil {
		fmt.Printf("Sender: error initializing %s", err)
		c <- 1
	}
	defer func() {
		fmt.Println("Sender: finished")
		c <- 0
	}()

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
	oflag := posix_mq.O_RDONLY
	mq, err := posix_mq.NewMessageQueue("/posix_mq_example", oflag, 0666, nil)
	if err != nil {
		fmt.Printf("Receiver: error initializing %s", err)
		c <- 1
	}
	defer func() {
		closeQueue(mq)
		fmt.Println("Receiver: finished")
		c <- 0
	}()

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
	err := mq.Unlink()
	if err != nil {
		log.Println(err)
	}
}
