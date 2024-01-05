package errors

/*
#cgo LDFLAGS: -lrt
#include <stdlib.h>
#include <signal.h>
#include <fcntl.h>
#include <mqueue.h>
*/

import "C"

const (
	//The queue was empty, and the O_NONBLOCK flag was set for the message queue description referred to by mqdes.
	EAGAIN = C.EAGAIN
	//The descriptor specified in mqdes was invalid or not opened for reading.
	EBADF = C.EBADF
	//The call was interrupted by a signal handler
	EINT = C.EINTR
	//The call would have blocked, and abs_timeout was invalid either because tv_sec was less than zero, or because tv_nsec was less than zero or greater than 1000 million.
	EINVAL = C.EINVAL
	//msg_len was less than the mq_msgsize attribute of the message queue.
	EMSGSIZE = C.EMSGSIZE
	//The call timed out before a message could be transferred.
	ETIMEDOUT = C.ETIMEDOUT
)
