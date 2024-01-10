package responder

import (
	"errors"
	"fmt"
	"github.com/joe-at-startupmedia/posix_mq"
	"syscall"
	"time"
)

type ResponderCallback func(msq []byte) (processed []byte, err error)

var (
	mqSend *posix_mq.MessageQueue
	mqResp *posix_mq.MessageQueue
)

func New(mqFile string, mqDir string, owner posix_mq.Ownership, flags int) error {
	sender, err := openQueue(mqFile+"_send", mqDir, owner, flags)
	if err != nil {
		return err
	}
	mqSend = sender

	responder, err := openQueue(mqFile+"_resp", mqDir, owner, flags)
	mqResp = responder

	return err
}

func openQueue(mqFile string, mqDir string, owner posix_mq.Ownership, flags int) (*posix_mq.MessageQueue, error) {
	//mq_open checks that the name starts with a slash (/), giving the EINVAL error if it does not
	if flags == 0 {
		flags = posix_mq.O_RDWR | posix_mq.O_CREAT | posix_mq.O_NONBLOCK
	}
	var (
		messageQueue *posix_mq.MessageQueue
		err          error
	)
	if owner.IsValid() {
		messageQueue, err = posix_mq.NewMessageQueue("/"+mqFile, flags, 0660, nil)

	} else {
		messageQueue, err = posix_mq.NewMessageQueue("/"+mqFile, flags, 0666, nil)
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create message queue %s: %-v", "/"+mqFile, err))
	}
	err = owner.ApplyPermissions(mqDir + mqFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not apply permissions %s: %-v", mqDir+mqFile, err))
	}
	return messageQueue, nil
}

func HandleRequest(msgHandler ResponderCallback) error {
	return handleRequest(msgHandler, 0)
}

func HandleRequestWithLag(msgHandler ResponderCallback, lag int) error {
	return handleRequest(msgHandler, lag)
}

func handleRequest(msgHandler ResponderCallback, lag int) error {
	msg, _, err := mqSend.Receive()
	if err != nil {
		//EAGAIN simply means the queue is empty when O_NONBLOCK is set
		if errors.Is(err, syscall.EAGAIN) {
			return nil
		}
		return err
	}
	processed, err := msgHandler(msg)
	if err != nil {
		return err
	}

	if lag > 0 {
		time.Sleep(time.Duration(lag) * time.Second)
	}

	err = mqResp.Send(processed, 0)
	return err
}

func closeQueue(mq *posix_mq.MessageQueue) error {
	return mq.Unlink()
}

func Close() error {
	if err := closeQueue(mqSend); err != nil {
		return err
	}
	return closeQueue(mqResp)
}
