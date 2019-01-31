package iocommon

import (
	"errors"
	"io"
)

var (
	errLimitedReaderSeekInvalidPosition = errors.New("iocommon.LimitedReader.Seek: invalid position")
	errLimitedReaderSeekInvalidWhence   = errors.New("iocommon.LimitedReader.Seek: invalid whence")
)

type LimitedReader struct {
	r                      ReadSeekCloser
	readerStart, size, pos int64
}

func NewLimitedReader(r ReadSeekCloser, readerStart int64, size int64) (*LimitedReader, error) {
	if readerStart != 0 {
		if _, err := r.Seek(readerStart, io.SeekStart); err != nil {
			return nil, err
		}
	}

	return &LimitedReader{r: r, readerStart: readerStart, size: size}, nil
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
	return
}

func (r *LimitedReader) Seek(offset int64, whence int) (ret int64, err error) {
	defer func() {
		if err == nil {
			r.pos = offset - r.readerStart
			ret = r.pos
		}
	}()
	switch whence {
	case io.SeekStart:
		if offset > r.size {
			return 0, errLimitedReaderSeekInvalidPosition
		}
		offset = r.readerStart + offset
	case io.SeekCurrent:
		if offset > (r.size - r.pos) {
			return 0, errLimitedReaderSeekInvalidPosition
		}
		offset = r.readerStart + r.pos + offset
	case io.SeekEnd:
		if offset > r.size {
			return 0, errLimitedReaderSeekInvalidPosition
		}
		offset = r.readerStart + r.size - offset
	default:
		return 0, errLimitedReaderSeekInvalidWhence
	}
	ret, err = r.r.Seek(offset, io.SeekStart)
	return
}

func (r *LimitedReader) Close() error {
	return r.r.Close()
}