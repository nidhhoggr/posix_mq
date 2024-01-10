package sender

import (
	"errors"
	"fmt"
	"github.com/joe-at-startupmedia/posix_mq"
	"time"
)

var (
	mqSend *posix_mq.MessageQueue
	mqResp *posix_mq.MessageQueue
)

func New(mqFile string, mqDir string, owner posix_mq.Ownership) error {
	sender, err := openQueue(mqFile+"_send", mqDir, owner)
	if err != nil {
		return err
	}
	mqSend = sender

	responder, err := openQueue(mqFile+"_resp", mqDir, owner)
	mqResp = responder

	return err
}

func openQueue(mqFile string, mqDir string, owner posix_mq.Ownership) (*posix_mq.MessageQueue, error) {
	oflag := posix_mq.O_RDWR

	var (
		messageQueue *posix_mq.MessageQueue
		err          error
	)
	if owner.IsValid() {
		messageQueue, err = posix_mq.NewMessageQueue("/"+mqFile, oflag, 0660, nil)

	} else {
		messageQueue, err = posix_mq.NewMessageQueue("/"+mqFile, oflag, 0666, nil)
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create message queue %s: %-v", "/"+mqFile, err))
	}
	return messageQueue, nil
}

func Send(data []byte, priority uint) error {
	return mqSend.Send(data, priority)
}

func WaitForResponse(duration time.Duration) ([]byte, uint, error) {
	return mqResp.TimedReceive(time.Now().Local().Add(duration))
}

func closeQueue(mq *posix_mq.MessageQueue) error {
	return mq.Close()
}

func Close() error {
	err := closeQueue(mqSend)
	if err != nil {
		return err
	}
	err = closeQueue(mqResp)
	return err
}
