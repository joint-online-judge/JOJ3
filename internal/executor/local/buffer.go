package local

import (
	"bytes"
	"errors"
)

// LimitedBuffer wraps a bytes.Buffer and limits its size.
type LimitedBuffer struct {
	buf     *bytes.Buffer
	maxSize int
}

// Write writes data to the buffer and checks the size limit.
func (lb *LimitedBuffer) Write(p []byte) (n int, err error) {
	if lb.buf.Len()+len(p) > lb.maxSize {
		// Truncate to fit within the limit
		allowed := lb.maxSize - lb.buf.Len()
		if allowed > 0 {
			n, _ = lb.buf.Write(p[:allowed])
		}
		return n, errors.New("buffer size limit exceeded")
	}
	return lb.buf.Write(p)
}

// Bytes returns the buffer's content.
func (lb *LimitedBuffer) Bytes() []byte {
	return lb.buf.Bytes()
}

// String returns the buffer's content as a string.
func (lb *LimitedBuffer) String() string {
	return lb.buf.String()
}
