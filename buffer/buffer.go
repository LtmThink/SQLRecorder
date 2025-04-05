package buffer

import (
	. "SQLRecorder/utils"
	"io"
	"net"
	"time"
)

const defaultBufSize = 4096
const maxCachedBufSize = 256 * 1024
const maxPacketSize = 1<<24 - 1

// A Buffer which is used for both reading and writing.
// This is possible since communication on each connection is synchronous.
// In other words, we can't write and read simultaneously on the same connection.
// The Buffer is similar to bufio.Reader / Writer but zero-copy-ish
// Also highly optimized for this particular use case.
// This Buffer is backed by two byte slices in a double-buffering scheme
type Buffer struct {
	buf     []byte // buf is a byte Buffer who's length and capacity are equal.
	nc      net.Conn
	idx     int
	length  int
	timeout time.Duration
	dbuf    [2][]byte // dbuf is an array with the two byte slices that back this Buffer
	flipcnt uint      // flipccnt is the current Buffer counter for double-buffering
}

// NewBuffer allocates and returns a new Buffer.
func NewBuffer(nc net.Conn) Buffer {
	fg := make([]byte, defaultBufSize)
	return Buffer{
		buf:  fg,
		nc:   nc,
		dbuf: [2][]byte{fg, nil},
	}
}

// Flip replaces the active Buffer with the background Buffer
// this is a delayed Flip that simply increases the Buffer counter;
// the actual Flip will be performed the next time we call `Buffer.Fill`
func (b *Buffer) Flip() {
	b.flipcnt += 1
}

// Fill reads into the Buffer until at least _need_ bytes are in it
func (b *Buffer) Fill(need int) error {
	n := b.length
	// Fill packets into its double-buffering target: if we've called
	// Flip on this Buffer, we'll be copying to the background Buffer,
	// and then filling it with network packets; otherwise we'll just move
	// the contents of the current Buffer to the front before filling it
	dest := b.dbuf[b.flipcnt&1]

	// grow Buffer if necessary to fit the whole message.
	if need > len(dest) {
		// Round up to the next multiple of the default size
		dest = make([]byte, ((need/defaultBufSize)+1)*defaultBufSize)

		// if the allocated Buffer is not too large, move it to backing storage
		// to prevent extra allocations on applications that perform large reads
		if len(dest) <= maxCachedBufSize {
			b.dbuf[b.flipcnt&1] = dest
		}
	}

	// if we're filling the fg Buffer, move the existing packets to the start of it.
	// if we're filling the bg Buffer, copy over the packets
	if n > 0 {
		copy(dest[:n], b.buf[b.idx:])
	}

	b.buf = dest
	b.idx = 0

	for {
		if b.timeout > 0 {
			if err := b.nc.SetReadDeadline(time.Now().Add(b.timeout)); err != nil {
				return err
			}
		}

		nn, err := b.nc.Read(b.buf[n:])
		n += nn

		switch err {
		case nil:
			if n < need {
				continue
			}
			b.length = n
			return nil

		case io.EOF:
			if n >= need {
				b.length = n
				return nil
			}
			return io.ErrUnexpectedEOF

		default:
			return err
		}
	}
}

// returns next N bytes from Buffer.
// The returned slice is only guaranteed to be valid until the next read
func (b *Buffer) ReadNext(need int) ([]byte, error) {
	if b.length < need {
		// refill
		if err := b.Fill(need); err != nil {
			return nil, err
		}
	}

	offset := b.idx
	b.idx += need
	b.length -= need
	return b.buf[offset:b.idx], nil
}

// TakeBuffer returns a Buffer with the requested size.
// If possible, a slice from the existing Buffer is returned.
// Otherwise a bigger Buffer is made.
// Only one Buffer (total) can be used at a time.
func (b *Buffer) TakeBuffer(length int) ([]byte, error) {
	if b.length > 0 {
		return nil, ErrBusyBuffer
	}

	// test (cheap) general case first
	if length <= cap(b.buf) {
		return b.buf[:length], nil
	}

	if length < maxPacketSize {
		b.buf = make([]byte, length)
		return b.buf, nil
	}

	// Buffer is larger than we want to Store.
	return make([]byte, length), nil
}

// TakeSmallBuffer is shortcut which can be used if length is
// known to be smaller than defaultBufSize.
// Only one Buffer (total) can be used at a time.
func (b *Buffer) TakeSmallBuffer(length int) ([]byte, error) {
	if b.length > 0 {
		return nil, ErrBusyBuffer
	}
	return b.buf[:length], nil
}

func (b *Buffer) ReadAll() ([]byte, error) {
	for {
		n, err := b.nc.Read(b.buf)
		if err != nil {
			return nil, err
		}
		return b.buf[:n], nil
	}
}

// TakeCompleteBuffer returns the complete existing Buffer.
// This can be used if the necessary Buffer size is unknown.
// cap and len of the returned Buffer will be equal.
// Only one Buffer (total) can be used at a time.
func (b *Buffer) TakeCompleteBuffer() ([]byte, error) {
	if b.length > 0 {
		return nil, ErrBusyBuffer
	}
	return b.buf, nil
}

// Store stores buf, an updated Buffer, if its suitable to do so.
func (b *Buffer) Store(buf []byte) error {
	if b.length > 0 {
		return ErrBusyBuffer
	} else if cap(buf) <= maxPacketSize && cap(buf) > cap(b.buf) {
		b.buf = buf[:cap(buf)]
	}
	return nil
}
