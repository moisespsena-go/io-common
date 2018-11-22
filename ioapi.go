package ioapi

import "io"

type ReadSeekCloser interface {
  io.ReadSeeker
  io.Closer
}
