package iocommon

import (
	"io"
)


type LimitedReader struct {
	r io.ReadSeeker
	LimitedSeeker
}

func NewLimitedReader(r io.ReadSeeker, readerStart int64, size int64) (*LimitedReader, error) {
	if readerStart != 0 {
		if _, err := r.Seek(readerStart, io.SeekStart); err != nil {
			return nil, err
		}
	}

	return &LimitedReader{
		r,
		LimitedSeeker{startPos: readerStart, size: size, do: r.Seek},
	}, nil
}

func (r *LimitedReader) Read(p []byte) (n int, err error) {
	available := int(r.size - r.pos)
	if available == 0 {
		return 0, io.EOF
	}
	if len(p) > available {
		p = p[0:available]
	}
	n, err = r.r.Read(p)
	if err == nil {
		r.pos += int64(n)
	}
	if err != nil {
		return
	}
	return
}

type LimitedReadCloser struct {
	LimitedReader
	io.Closer
}

func NewLimitedReadCloser(r ReadSeekCloser, readerStart int64, size int64) (*LimitedReadCloser, error) {
	lr, err := NewLimitedReader(r, readerStart, size)
	if err != nil {
		return nil, err
	}
	return &LimitedReadCloser{*lr, r}, nil
}
