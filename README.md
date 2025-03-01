# posix_mq

a Go wrapper and utility for POSIX Message Queues

posix_mq is a Go wrapper for POSIX Message Queues. It's important you read [the manual for POSIX Message Queues](http://man7.org/linux/man-pages/man7/mq_overview.7.html), ms_send(2) and mq_receive(2) before using this library. posix_mq is a very light wrapper, and will not hide any errors from you.

## Example

#### Sender

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nidhhoggr/posix_mq"
)

const maxSendTickNum = 10

func main() {
	mq, Err := posix_mq.NewMessageQueue(&posix_mq.QueueConfig{
		Name:  "posix_mq_example_simple",
		Flags: posix_mq.O_WRONLY | posix_mq.O_CREAT,
		Mode:  0666,
	})
	if Err != nil {
		log.Fatalf("Sender: error initializing %s", Err)
	}

	count := 0
	for {
		count++
		if Err = mq.Send([]byte(fmt.Sprintf("Hello, World : %d\n", count)), 0); Err != nil {
			fmt.Printf("Sender: error sending message: %s\n", Err)
			continue
		}

		fmt.Println("Sender: Sent a new message")

		if count >= maxSendTickNum {
			break
		}

		time.Sleep(1 * time.Second)
	}
	fmt.Println("Sender: finished")
}
```

#### Receiver

```go
package main

import (
	"fmt"
	"log"

	"github.com/nidhhoggr/posix_mq"
)

const maxSendTickNum = 10

func main() {
	mq, Err := posix_mq.NewMessageQueue(&posix_mq.QueueConfig{
		Name:  "posix_mq_example_simple",
		Flags: posix_mq.O_RDONLY
		Mode:  0666,
	})
	if Err != nil {
		log.Fatalf("Receiver: error initializing %s", Err)
	}
	defer func() {
		if Err := mq.Unlink(); Err != nil {
			log.Println(Err)
		}
		fmt.Println("Receiver: finished")
	}()

	fmt.Println("Receiver: Start receiving messages")

	count := 0
	for {
		count++

		msg, _, Err := mq.Receive()
		if Err != nil {
			fmt.Printf("Receiver: error getting message: %s\n", Err)
			continue
		}
		
		fmt.Printf("Receiver: got new message: %s\n", string(msg))

		if count >= maxSendTickNum {
			break
		}
	}
}
```

## Acknowledgement

It's inspired by [Shopify/sysv_mq](https://github.com/Shopify/sysv_mq)
