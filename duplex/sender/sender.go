package sender

import (
	"errors"
	"fmt"
	"github.com/joe-at-startupmedia/posix_mq"
	"time"
)

type MqSender struct {
	mqSend *posix_mq.MessageQueue
	mqResp *posix_mq.MessageQueue
}

func New(config posix_mq.QueueConfig, owner *posix_mq.Ownership) (*MqSender, error) {
	sender, err := openQueue(config, owner, "send")
	if err != nil {
		return nil, err
	}

	responder, err := openQueue(config, owner, "resp")

	mqs := MqSender{
		sender,
		responder,
	}

	return &mqs, err
}

func openQueue(config posix_mq.QueueConfig, owner *posix_mq.Ownership, postfix string) (*posix_mq.MessageQueue, error) {
	if config.Flags == 0 {
		config.Flags = posix_mq.O_RDWR
	}
	config.Name = fmt.Sprintf("%s_%s", config.Name, postfix)
	var (
		messageQueue *posix_mq.MessageQueue
		err          error
	)
	if owner != nil && owner.IsValid() {
		config.Mode = 0660
		messageQueue, err = posix_mq.NewMessageQueue(&config)

	} else {
		config.Mode = 0666
		messageQueue, err = posix_mq.NewMessageQueue(&config)
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create message queue %s: %-v", config.GetFile(), err))
	}
	return messageQueue, nil
}

func (mqs *MqSender) Send(data []byte, priority uint) error {
	return mqs.mqSend.Send(data, priority)
}

func (mqs *MqSender) WaitForResponse(duration time.Duration) ([]byte, uint, error) {
	return mqs.mqResp.TimedReceive(time.Now().Local().Add(duration))
}

func closeQueue(mq *posix_mq.MessageQueue) error {
	return mq.Close()
}

func (mqs *MqSender) Close() error {
	if err := closeQueue(mqs.mqSend); err != nil {
		return err
	}
	return closeQueue(mqs.mqResp)
}
