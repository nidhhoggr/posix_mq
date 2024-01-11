package posix_mq

import (
	"os"
	"syscall"
	"testing"
)

func SampleMessageQueue(t *testing.T) *MessageQueue {
	config := QueueConfig{
		Name:  "pmq_testing",
		Mode:  0660,
		Flags: O_RDWR | O_CREAT,
	}

	mq, err := NewMessageQueue(&config)

	if err != nil {
		t.Error(err)
	}

	return mq
}

func Test_SendMessage(t *testing.T) {
	mq := SampleMessageQueue(t)

	wired := "Narwhals and ice cream"

	err := mq.Send([]byte(wired), 0)

	if err != nil {
		t.Error(err)
	}

	response, _, err := mq.Receive()

	if err != nil {
		t.Error(err)
	}

	if wired != string(response) {
		t.Errorf("expected %s, got: %s", wired, response)
	}
}

func Test_QueuePriority(t *testing.T) {
	mq := SampleMessageQueue(t)

	wired := "Narwhals and ice cream"

	err := mq.Send([]byte(wired), 3)

	if err != nil {
		t.Error(err)
	}

	_, mtype, err := mq.Receive()

	if err != nil {
		t.Error(err)
	}

	if mtype != 3 {
		t.Errorf("expected mtype 4, got: %d", mtype)
	}
}

func Test_QueueCount(t *testing.T) {
	mq := SampleMessageQueue(t)

	wired := "Narwhals and ice cream"

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
}

func Test_QueueClose(t *testing.T) {
	mq := SampleMessageQueue(t)
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
	if _, err := os.Stat(POSIX_MQ_DIR + "pmq_testing"); err != nil {
		t.Errorf("got an unexpected error %s", err)
	}
}

func Test_QueueUnlink(t *testing.T) {
	mq := SampleMessageQueue(t)
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
	if _, err := os.Stat(POSIX_MQ_DIR + "pmq_testing"); err != nil {
		if os.IsNotExist(err) {
			t.Log("Expected file to not exist")
		} else {
			t.Fatalf("got an unexpected error %s", err)
		}
	}
}
