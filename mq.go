package posix_mq

import (
	"syscall"
	"time"
)

// Represents the message queue
type MessageQueue struct {
	handler int
	name    string
	recvBuf *receiveBuffer
}

// QueueConfig is used to configure an instance of the message queue.
type QueueConfig struct {
	Name  string
	Dir   string
	Flags int
	Mode  int // The mode of the message queue, e.g. 0600
	Attrs *MessageQueueAttribute
}

type MessageQueueAttribute struct {
	Flags   int // Flags (ignored for mq_open())
	MaxMsg  int // Max. # of messages on queue
	MsgSize int // Max. message size (bytes)
	MsgCnt  int // # of messages in the queue
}

const POSIX_MQ_DIR = "/dev/mqueue/"

// NewMessageQueue returns an instance of the message queue given a QueueConfig.
func NewMessageQueue(config *QueueConfig) (*MessageQueue, error) {

	//mq_open checks that the name starts with a slash (/), giving the EINVAL error if it does not
	name := "/" + config.Name
	h, err := mq_open(name, config.Flags, config.Mode, config.Attrs)
	if err != nil {
		return nil, err
	}

	msgSize := MSGSIZE_DEFAULT
	if config.Attrs != nil {
		msgSize = config.Attrs.MsgSize
	}
	recvBuf, err := newReceiveBuffer(msgSize)
	if err != nil {
		return nil, err
	}

	return &MessageQueue{
		handler: h,
		name:    name,
		recvBuf: recvBuf,
	}, nil
}

// Send sends message to the message queue.
func (mq *MessageQueue) Send(data []byte, priority uint) error {
	return mq_send(mq.handler, data, priority)
}

// TimedSend sends message to the message queue with a ceiling on the time for which the call will block.
func (mq *MessageQueue) TimedSend(data []byte, priority uint, t time.Time) error {
	return mq_timedsend(mq.handler, data, priority, t)
}

// Receive receives message from the message queue.
func (mq *MessageQueue) Receive() ([]byte, uint, error) {
	return mq_receive(mq.handler, mq.recvBuf)
}

// TimedReceive receives message from the message queue with a ceiling on the time for which the call will block.
func (mq *MessageQueue) TimedReceive(t time.Time) ([]byte, uint, error) {
	return mq_timedreceive(mq.handler, mq.recvBuf, t)
}

// FIXME Don't work because of signal portability.
// Notify set signal notification to handle new messages.
func (mq *MessageQueue) Notify(sigNo syscall.Signal) error {
	return mq_notify(mq.handler, int(sigNo))
}

// Close closes the message queue.
func (mq *MessageQueue) Close() error {
	mq.recvBuf.free()
	return mq_close(mq.handler)
}

// Unlink deletes the message queue.
func (mq *MessageQueue) Unlink() error {
	if err := mq.Close(); err != nil {
		return err
	}
	return mq_unlink(mq.name)
}

func (config *QueueConfig) GetFile() string {
	if len(config.Dir) == 0 {
		return POSIX_MQ_DIR + config.Name
	} else {
		return config.Dir + config.Name
	}
}

func (mq *MessageQueue) Count() (int, error) {
	mqa, err := mq_getattr(mq.handler)
	return mqa.MsgCnt, err
}
