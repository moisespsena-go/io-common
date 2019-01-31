package iocommon

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestLimitedReader_Read(t *testing.T) {
	b := []byte("0123456789")
	r := func() ReadSeekCloser {
		return NewBytesReadCloser(b)
	}

	type fields struct {
		r           ReadSeekCloser
		readerStart int64
		size        int64
		pos         int64
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		{"T1", fields{r(), 0, 6, 0}, args{[]byte("012345")}, 6, false},
		{"T2", fields{r(), 1, 6, 0}, args{[]byte("123456")}, 6, false},
		{"T3", fields{r(), 2, 6, 0}, args{[]byte("234567")}, 6, false},
		{"T4", fields{r(), 3, 6, 0}, args{[]byte("345678")}, 6, false},
		{"T5", fields{r(), 4, 6, 0}, args{[]byte("456789")}, 6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &LimitedReader{
				r:           tt.fields.r,
				readerStart: tt.fields.readerStart,
				size:        tt.fields.size,
				pos:         tt.fields.pos,
			}
			gotN, err := r.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("LimitedReader.Read() error = %v, seekWantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("LimitedReader.Read() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestLimitedReader_ReadAll(t *testing.T) {
	b := []byte("0123456789")
	r := func() ReadSeekCloser {
		return NewBytesReadCloser(b)
	}

	type fields struct {
		r           ReadSeekCloser
		readerStart int64
		size        int64
		pos         int64
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		{"T", fields{r(), 0, 6, 0}, args{[]byte("012345")}, 6, false},
		{"T", fields{r(), 1, 6, 0}, args{[]byte("123456")}, 6, false},
		{"T", fields{r(), 2, 6, 0}, args{[]byte("234567")}, 6, false},
		{"T", fields{r(), 3, 6, 0}, args{[]byte("345678")}, 6, false},
		{"T", fields{r(), 4, 6, 0}, args{[]byte("456789")}, 6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			r := &LimitedReader{
				r:           tt.fields.r,
				readerStart: tt.fields.readerStart,
				size:        tt.fields.size,
				pos:         tt.fields.pos,
			}
			tt.args.p, err = ioutil.ReadAll(r)
			gotN := len(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("LimitedReader.Read() error = %v, seekWantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("LimitedReader.Read() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
func TestLimitedReader_SeekReader(t *testing.T) {
	b := []byte("0123456789")
	r := func() ReadSeekCloser {
		return NewBytesReadCloser(b)
	}
	type fields struct {
		r           ReadSeekCloser
		readerStart int64
		size        int64
		pos         int64
	}
	type args struct {
		offset int64
		whence int
		data   []byte
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		seekWantRet int64

		seekWantErr bool
		readWantErr bool
	}{
		{"SeekStart", fields{r(), 0, 6, 0}, args{0, io.SeekStart, []byte("012345")}, 0, false, false},
		{"SeekStart", fields{r(), 0, 6, 0}, args{1, io.SeekStart, []byte("12345")}, 1, false, false},
		{"SeekStart", fields{r(), 0, 6, 0}, args{2, io.SeekStart, []byte("2345")}, 2, false, false},
		{"SeekStart", fields{r(), 0, 6, 0}, args{3, io.SeekStart, []byte("345")}, 3, false, false},
		{"SeekStart", fields{r(), 0, 6, 0}, args{4, io.SeekStart, []byte("45")}, 4, false, false},
		{"SeekStart", fields{r(), 0, 6, 0}, args{5, io.SeekStart, []byte("5")}, 5, false, false},
		{"SeekStart", fields{r(), 0, 6, 0}, args{6, io.SeekStart, []byte("")}, 6, false, false},

		{"SeekCurrent", fields{r(), 0, 6, 0}, args{0, io.SeekCurrent, []byte("012345")}, 0, false, false},
		{"SeekCurrent", fields{r(), 0, 6, 1}, args{1, io.SeekCurrent, []byte("2345")}, 2, false, false},
		{"SeekCurrent", fields{r(), 0, 6, 2}, args{1, io.SeekCurrent, []byte("345")}, 3, false, false},
		{"SeekCurrent", fields{r(), 0, 6, 3}, args{1, io.SeekCurrent, []byte("45")}, 4, false, false},
		{"SeekCurrent", fields{r(), 0, 6, 4}, args{1, io.SeekCurrent, []byte("5")}, 5, false, false},
		{"SeekCurrent", fields{r(), 0, 6, 5}, args{1, io.SeekCurrent, []byte("")}, 6, false, false},

		{"SeekEnd.Pos=0", fields{r(), 0, 6, 0}, args{6, io.SeekEnd, []byte("012345")}, 0, false, false},
		{"SeekEnd.Pos=0", fields{r(), 0, 6, 0}, args{5, io.SeekEnd, []byte("12345")}, 1, false, false},
		{"SeekEnd.Pos=0", fields{r(), 0, 6, 0}, args{4, io.SeekEnd, []byte("2345")}, 2, false, false},
		{"SeekEnd.Pos=0", fields{r(), 0, 6, 0}, args{3, io.SeekEnd, []byte("345")}, 3, false, false},
		{"SeekEnd.Pos=0", fields{r(), 0, 6, 0}, args{2, io.SeekEnd, []byte("45")}, 4, false, false},
		{"SeekEnd.Pos=0", fields{r(), 0, 6, 0}, args{1, io.SeekEnd, []byte("5")}, 5, false, false},
		{"SeekEnd.Pos=0", fields{r(), 0, 6, 0}, args{0, io.SeekEnd, []byte("")}, 6, false, false},

		{"SeekEnd.Pos=2", fields{r(), 0, 6, 2}, args{4, io.SeekEnd, []byte("2345")}, 2, false, false},
		{"SeekEnd.Pos=3", fields{r(), 0, 6, 3}, args{4, io.SeekEnd, []byte("2345")}, 2, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &LimitedReader{
				r:           tt.fields.r,
				readerStart: tt.fields.readerStart,
				size:        tt.fields.size,
				pos:         tt.fields.pos,
			}
			gotRet, err := r.Seek(tt.args.offset, tt.args.whence)
			if (err != nil) != tt.seekWantErr {
				t.Errorf("LimitedReader.Seek() error = %v, seekWantErr %v", err, tt.seekWantErr)
				return
			}
			if gotRet != tt.seekWantRet {
				t.Errorf("LimitedReader.Seek() = %v, want %v", gotRet, tt.seekWantRet)
			}

			data, err := ioutil.ReadAll(r)
			if (err != nil) != tt.readWantErr {
				t.Errorf("LimitedReader.Read() error = %v, readWantErr %v", err, tt.readWantErr)
				return
			}
			if bytes.Compare(data, tt.args.data) != 0 {
				t.Errorf("LimitedReader.Read() data = %q, whant args.data %q", string(data), string(tt.args.data))
				return
			}
		})
	}
}
