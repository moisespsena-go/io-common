package iocommon

import "io"

func MustCloser(closer io.Closer) func() {
	return func() {
		closer.Close()
	}
}
