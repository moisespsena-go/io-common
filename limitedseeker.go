package iocommon

import (
	"errors"
	"io"
)

var (
	errLimitedSeekInvalidPosition = errors.New("iocommon.LimitedReader.Seek: invalid position")
	errLimitedSeekInvalidWhence   = errors.New("iocommon.LimitedReader.Seek: invalid whence")
)

type LimitedSeeker struct {
	startPos, size, pos int64
	do                  func(offset int64, whence int) (ret int64, err error)
}

func NewLimitedSeeker(startPos int64, size int64, do func(offset int64, whence int) (ret int64, err error)) *LimitedSeeker {
	return &LimitedSeeker{startPos: startPos, size: size, do: do}
}

func (this *LimitedSeeker) SetSize(size int64) {
	this.size = size
}

func (this *LimitedSeeker) StartPos() int64 {
	return this.startPos
}

func (this *LimitedSeeker) SetStartPos(startPos int64) {
	this.startPos = startPos
}

func (this *LimitedSeeker) Seek(offset int64, whence int) (ret int64, err error) {
	defer func() {
		if err == nil {
			this.pos = offset - this.startPos
			ret = this.pos
		}
	}()
	switch whence {
	case io.SeekStart:
		if offset > this.size {
			return 0, errLimitedSeekInvalidPosition
		}
		offset = this.startPos + offset
	case io.SeekCurrent:
		if offset > (this.size - this.pos) {
			return 0, errLimitedSeekInvalidPosition
		}
		offset = this.startPos + this.pos + offset
	case io.SeekEnd:
		if offset > this.size {
			return 0, errLimitedSeekInvalidPosition
		}
		offset = this.startPos + this.size - offset
	default:
		return 0, errLimitedSeekInvalidWhence
	}
	ret, err = this.do(offset, io.SeekStart)
	return
}

func (r *LimitedSeeker) Size() int64 {
	return r.size
}
