package posix_mq_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/nidhhoggr/posix_mq"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

const wired = "Narwhals and ice cream"

// Open non-exist queue without O_CREAT will return syscall.ENOENT
func TestOpenMQWithOutCreatePermission(t *testing.T) {
	posix_mq.ForceRemoveQueue("pmq_testing_wocp")
	config := posix_mq.QueueConfig{
		Name:  "pmq_testing_wocp",
		Mode:  0660,
		Flags: posix_mq.O_WRONLY,
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, mqt)
	mqErr, ok := err.(syscall.Errno)
	assertTrue(t, ok)
	assertEqual(t, syscall.ENOENT, mqErr)
}

// create queue with invalid name will return syscall.EACCES
func TestOpenMQWithWrongName(t *testing.T) {
	posix_mq.ForceRemoveQueue("pmq_testing_wn")
	//the library already prepends a forward slash
	config := posix_mq.QueueConfig{
		Name:  "/pmq_testing_wn",
		Mode:  0660,
		Flags: posix_mq.O_WRONLY | posix_mq.O_CREAT,
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, mqt)
	assertNotNil(t, err)
	mqErr, ok := err.(syscall.Errno)
	assertTrue(t, ok)
	assertEqual(t, syscall.EACCES, mqErr)
}

func TestOpenMQSuccess(t *testing.T) {
	posix_mq.ForceRemoveQueue("pmq_testing_opensuccess")
	config := posix_mq.QueueConfig{
		Name:  "pmq_testing_opensuccess",
		Mode:  0660,
		Flags: posix_mq.O_WRONLY | posix_mq.O_CREAT,
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, err)
	assertNotNil(t, mqt)
	err = mqt.Unlink()
	assertNil(t, err)
}

func TestOpenExistMQ(t *testing.T) {
	posix_mq.ForceRemoveQueue("pmq_testing_openexisting")
	config := posix_mq.QueueConfig{
		Name:  "pmq_testing_openexisting",
		Mode:  0660,
		Flags: posix_mq.O_WRONLY | posix_mq.O_CREAT,
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, err)
	assertNotNil(t, mqt)
	mqt2, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, err)
	assertNotNil(t, mqt)
	err = mqt.Unlink()
	assertNil(t, err)
	err = mqt2.Unlink()
	assertNotNil(t, err)
	assertEqual(t, syscall.ENOENT, err.(syscall.Errno))
}

// create a queue with MasMsg larger than /proc/sys/fs/mqueue/msg_max will return syscall.EINVAL
func TestCreateMQWithMaxMsgOverLimit(t *testing.T) {
	posix_mq.ForceRemoveQueue("pmq_testing_mmol")
	config := posix_mq.QueueConfig{
		Name:  "pmq_testing_mmol",
		Mode:  0660,
		Flags: posix_mq.O_WRONLY | posix_mq.O_CREAT,
		Attrs: &posix_mq.MessageQueueAttribute{
			MaxMsg: 1000000,
		},
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, mqt)
	assertNotNil(t, err)
	assertEqual(t, syscall.EINVAL, err.(syscall.Errno))
}

// creat an exist queue with O_EXCL will return syscall.ENINVAL
func TestExamQueueExist(t *testing.T) {
	posix_mq.ForceRemoveQueue("pmq_testing_openexcl")
	config := posix_mq.QueueConfig{
		Name:  "pmq_testing_openexcl",
		Mode:  0660,
		Flags: posix_mq.O_WRONLY | posix_mq.O_CREAT,
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, err)
	assertNotNil(t, mqt)
	config.Flags = posix_mq.O_WRONLY | posix_mq.O_CREAT | posix_mq.O_EXCL
	mqt2, err := posix_mq.NewMessageQueue(&config)
	assertNotNil(t, err)
	assertNil(t, mqt2)
	assertEqual(t, syscall.EEXIST, err.(syscall.Errno))
	err = mqt.Unlink()
	assertNil(t, err)
}

func Test_SendMessage(t *testing.T) {
	mqt := SampleMessageQueue(t, 0, "sendmsg")

	for i := 1; i <= 5; i++ {
		err := mqt.Send([]byte(wired), 0)
		assertNil(t, err)
	}

	err := mqt.Unlink()
	assertNil(t, err)
}

func Test_SendReceiveMessage(t *testing.T) {
	mqt := SampleMessageQueue(t, posix_mq.O_WRONLY|posix_mq.O_CREAT, "sendrcvmsg")
	mqt2 := SampleMessageQueue(t, posix_mq.O_RDONLY|posix_mq.O_CREAT, "sendrcvmsg")
	for i := 1; i <= 5; i++ {
		err := mqt.Send([]byte(wired), 0)
		assertNil(t, err)
		response, _, err := mqt2.Receive()

		if err != nil {
			t.Error(err)
		}

		if wired != string(response) {
			t.Errorf("expected %s, got: %s", wired, response)
		}
	}

	err := mqt.Unlink()
	assertNil(t, err)
	err = mqt2.Unlink()
	assertNotNil(t, err)
	assertEqual(t, syscall.ENOENT, err.(syscall.Errno))
}

type TestMsg struct {
	Type   uint8
	Length uint8
	Data   [1150]byte
}

// Send msg size larger than /proc/sys/fs/mqueue/msgsize_max wll return syscall.EMSGSIZE
func TestSendMsgTooLong(t *testing.T) {
	posix_mq.ForceRemoveQueue("pmq_testing_mtl")
	msgSize := int(unsafe.Sizeof(TestMsg{}))
	config := posix_mq.QueueConfig{
		Name:  "pmq_testing_mtl",
		Mode:  0660,
		Flags: posix_mq.O_WRONLY | posix_mq.O_CREAT,
		Attrs: &posix_mq.MessageQueueAttribute{
			MaxMsg:  10,
			MsgSize: msgSize - 1,
		},
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, err)
	assertNotNil(t, mqt)
	buf := bytes.NewBuffer(make([]byte, msgSize))
	err = binary.Write(buf, binary.LittleEndian, TestMsg{})
	assertNil(t, err)
	err = mqt.Send(buf.Bytes(), 0)
	assertNotNil(t, err)
	assertEqual(t, syscall.EMSGSIZE, err.(syscall.Errno))
	err = mqt.Unlink()
	assertNil(t, err)
}

type TestMsg2 struct {
	Type   uint8
	Length uint8
	Data   [21]byte
}

// recv msg size smaller than MsgSize, receive() return syscall.EMSGSIZE
func TestRecvMsgTooShort(t *testing.T) {
	posix_mq.ForceRemoveQueue("pmq_testing_mts")
	msgSize := int(unsafe.Sizeof(TestMsg2{}))
	config := posix_mq.QueueConfig{
		Name:  "pmq_testing_mts",
		Mode:  0660,
		Flags: posix_mq.O_WRONLY | posix_mq.O_CREAT,
		Attrs: &posix_mq.MessageQueueAttribute{
			MaxMsg:  10,
			MsgSize: msgSize,
		},
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, err)
	assertNotNil(t, mqt)
	config2 := posix_mq.QueueConfig{
		Name:  "pmq_testing_mts",
		Mode:  0660,
		Flags: posix_mq.O_RDONLY | posix_mq.O_CREAT,
		Attrs: &posix_mq.MessageQueueAttribute{
			MaxMsg:  10,
			MsgSize: msgSize - 1,
		},
	}
	mqt2, err := posix_mq.NewMessageQueue(&config2)
	assertNil(t, err)
	assertNotNil(t, mqt2)
	buf := bytes.NewBuffer(make([]byte, msgSize))
	msg := TestMsg2{
		Type: uint8(10),
	}
	buf.Reset()
	err = binary.Write(buf, binary.LittleEndian, msg)
	assertNil(t, err)
	err = mqt.Send(buf.Bytes(), 0)
	assertNil(t, err)
	recvMsg, prio, err := mqt2.Receive()
	assertEqual(t, 0, len(recvMsg))
	assertNotNil(t, err)
	assertEqual(t, syscall.EMSGSIZE, err.(syscall.Errno))
	assertEqual(t, uint(0), prio)
	err = mqt.Unlink()
	assertNil(t, err)
	err = mqt2.Unlink()
	assertNotNil(t, err)
	assertEqual(t, syscall.ENOENT, err.(syscall.Errno))
}

// with non-blocking queue, while queue full, send() will return syscall.EAGAIN
func TestSendwithNonblocking(t *testing.T) {
	posix_mq.ForceRemoveQueue("pmq_testing_sendwnblk")
	msgSize := int(unsafe.Sizeof(TestMsg2{}))
	config := posix_mq.QueueConfig{
		Name:  "pmq_testing_sendwnblk",
		Mode:  0660,
		Flags: posix_mq.O_WRONLY | posix_mq.O_CREAT | posix_mq.O_NONBLOCK,
		Attrs: &posix_mq.MessageQueueAttribute{
			MaxMsg:  1,
			MsgSize: msgSize,
		},
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, err)
	assertNotNil(t, mqt)
	buf := bytes.NewBuffer(make([]byte, msgSize))
	buf.Reset()
	err = binary.Write(buf, binary.LittleEndian, TestMsg2{Type: uint8(1)})
	assertNil(t, err)
	err = mqt.Send(buf.Bytes(), 0)
	assertNil(t, err)
	buf.Reset()
	err = binary.Write(buf, binary.LittleEndian, TestMsg2{Type: uint8(2)})
	assertNil(t, err)
	err = mqt.Send(buf.Bytes(), 1)
	assertNotNil(t, err)
	assertEqual(t, syscall.EAGAIN, err.(syscall.Errno))
	err = mqt.Unlink()
	assertNil(t, err)
}

func TestRecvwithNonblocking(t *testing.T) {
	posix_mq.ForceRemoveQueue("pmq_testing_recvwnblk")
	msgSize := int(unsafe.Sizeof(TestMsg2{}))
	config := posix_mq.QueueConfig{
		Name:  "pmq_testing_recvwnblk",
		Mode:  0660,
		Flags: posix_mq.O_RDONLY | posix_mq.O_CREAT | posix_mq.O_NONBLOCK,
		Attrs: &posix_mq.MessageQueueAttribute{
			MaxMsg:  1,
			MsgSize: msgSize,
		},
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, err)
	assertNotNil(t, mqt)

	msg, prio, err := mqt.Receive()
	assertNotNil(t, err)
	assertEqual(t, uint(0), prio)
	assertEqual(t, 0, len(msg))
	assertEqual(t, syscall.EAGAIN, err.(syscall.Errno))
	err = mqt.Unlink()
	assertNil(t, err)
}

func Test_QueuePriority(t *testing.T) {
	mq := SampleMessageQueue(t, 0, "qprio")

	err := mq.Send([]byte(wired), 3)

	if err != nil {
		t.Error(err)
	}

	_, mtype, err := mq.Receive()

	if err != nil {
		t.Error(err)
	}

	if mtype != 3 {
		t.Errorf("expected mtype 3, got: %d", mtype)
	}

	err = mq.Unlink()
	assertNil(t, err)
}

func Test_QueueCount(t *testing.T) {
	mq := SampleMessageQueue(t, 0, "qcnt")

	if err := mq.Send([]byte(wired), 0); err != nil {
		t.Error(err)
	}

	if err := mq.Send([]byte(wired), 0); err != nil {
		t.Error(err)
	}

	if count, _ := mq.Count(); count != 2 {
		t.Errorf("expected count 2, got: %d", count)
	}

	if _, _, err := mq.Receive(); err != nil {
		t.Error(err)
	}

	if count, _ := mq.Count(); count != 1 {
		t.Errorf("expected count 1, got: %d", count)
	}

	if _, _, err := mq.Receive(); err != nil {
		t.Error(err)
	}

	if count, _ := mq.Count(); count != 0 {
		t.Errorf("expected count 0, got: %d", count)
	}

	err := mq.Unlink()
	assertNil(t, err)
}

func Test_Notify(t *testing.T) {
	mq := SampleMessageQueue(t, 0, "qnot")

	mq.Notify(syscall.SIGUSR1)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGUSR1)
	signalsCaught := 0
	go func(test *testing.T) {
		for {
			s := <-sigc
			switch s {
			case syscall.SIGUSR1:
				test.Logf("catched signal: %-v", s)

				response, _, err := mq.Receive()

				if err != nil {
					test.Error(err)
				}

				if wired != string(response) {
					test.Errorf("expected %s, got: %s", wired, response)
				} else {
					test.Logf("Sucessfully notified with msg: %s", response)
					signalsCaught++
				}
			default:
				t.Logf("Caught an unexpected signal: %s", s)
			}
			if signalsCaught >= 2 {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}(t)

	if err := mq.TimedSend([]byte(wired), 0, time.Second*2); err != nil {
		t.Error(err)
	}

	time.Sleep(1 * time.Second)

	mq.Notify(syscall.SIGUSR1)

	if err := mq.TimedSend([]byte(wired), 0, time.Second*2); err != nil {
		t.Error(err)
	}

	time.Sleep(1 * time.Second)

	if signalsCaught != 2 {
		t.Errorf("expected catching 2 notifications, got: %d", signalsCaught)
	}

	err := mq.Unlink()
	assertNil(t, err)
}

func Test_QueueClose(t *testing.T) {
	mq := SampleMessageQueue(t, 0, "qcls")
	if err := mq.Close(); err != nil {
		t.Errorf("expected to close queue, got: %s", err)
	}
	if err := mq.Send([]byte("I'll never be sent :("), 0); err != nil {
		switch err {
		case nil:
			t.Error("Expected bad file descriptor error")
		case syscall.EBADF:
			t.Log("Received BAD file descriptor error")
		default:
			t.Fatalf("got an unexpected error %s", err)
		}
	}
	if _, err := os.Stat(posix_mq.POSIX_MQ_DIR + "pmq_testing_qcls"); err != nil {
		t.Errorf("got an unexpected error %s", err)
	}
	err := posix_mq.ForceRemoveQueue("pmq_testing_qcls")
	assertNil(t, err)
}

func Test_QueueUnlink(t *testing.T) {
	mq := SampleMessageQueue(t, 0, "qulnk")
	if err := mq.Unlink(); err != nil {
		t.Errorf("expected to close queue, got: %s", err)
	}
	if err := mq.Send([]byte("I'll never be sent :("), 0); err != nil {
		switch err {
		case nil:
			t.Error("Expected bad file descriptor error")
		case syscall.EBADF:
			t.Log("Received BAD file descriptor error")
		default:
			t.Fatalf("got an unexpected error %s", err)
		}
	}
	if _, err := os.Stat(posix_mq.POSIX_MQ_DIR + "pmq_testing_qulnk"); err != nil {
		if os.IsNotExist(err) {
			t.Log("Expected file to not exist")
		} else {
			t.Fatalf("got an unexpected error %s", err)
		}
	}
}

func SampleMessageQueue(t *testing.T, flags int, postfix string) *posix_mq.MessageQueue {

	if flags == 0 {
		flags = posix_mq.O_RDWR | posix_mq.O_CREAT
	}
	queueName := fmt.Sprintf("pmq_testing_%s", postfix)
	config := posix_mq.QueueConfig{
		Name:  queueName,
		Mode:  0660,
		Flags: flags,
	}
	mqt, err := posix_mq.NewMessageQueue(&config)
	assertNil(t, err)
	assertNotNil(t, mqt)
	return mqt
}

func assertNil(t *testing.T, i interface{}) {
	if !isNil(i) {
		t.Errorf("expected %-v to be nil", i)
	}
}

func assertNotNil(t *testing.T, i interface{}) {
	if isNil(i) {
		t.Errorf("expected %-v to not be nil", i)
	}
}

func assertEqual[T any](t *testing.T, ptr T, ptr2 T) {
	if !reflect.ValueOf(ptr).Equal(reflect.ValueOf(ptr2)) {
		t.Errorf("expected %-v to equal %-v", ptr, ptr2)
	}
}

func assertTrue(t *testing.T, val bool) {
	if !val {
		t.Errorf("expected %-v to be true", val)
	}
}

// containsKind checks if a specified kind in the slice of kinds.
func containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {
	for i := 0; i < len(kinds); i++ {
		if kind == kinds[i] {
			return true
		}
	}

	return false
}

// isNil checks if a specified object is nil or not, without Failing.
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	isNilableKind := containsKind(
		[]reflect.Kind{
			reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice, reflect.UnsafePointer},
		kind)

	if isNilableKind && value.IsNil() {
		return true
	}

	return false
}
