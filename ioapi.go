package iocommon

import (
	"bytes"
	"errors"
	"io"
)

var ErrNotSeekable = errors.New("not seekable")

type ReadSeekCloser interface {
	io.ReadSeeker
	io.Closer
}

type noSeeker struct {
	io.ReadCloser
}

func (noSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, ErrNotSeekable
}

func NoSeeker(r io.ReadCloser) ReadSeekCloser {
	return &noSeeker{r}
}

type BytesReadCloser struct {
	*bytes.Reader
}

func (b *BytesReadCloser) Close() error {
	b.Reader.Reset([]byte{})
	return nil
}

// NewBytesReadCloser returns a new BytesReadCloser reading from b.
func NewBytesReadCloser(b []byte) *BytesReadCloser {
	return &BytesReadCloser{bytes.NewReader(b)}
}
