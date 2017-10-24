package bytes

import (
	"io"
	"syscall"
	// "fmt"
)

// ---------------------------------------------------

type Reader struct {
	b   []byte
	off int
}

func NewReader(val []byte) *Reader {
	return &Reader{val, 0}
}

func (r *Reader) Len() int {
	if r.off >= len(r.b) {
		return 0
	}
	return len(r.b) - r.off
}

func (r *Reader) Bytes() []byte {
	return r.b[r.off:]
}

func (r *Reader) SeekToBegin() (err error) {
	r.off = 0
	return
}

func (r *Reader) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	case 0:
	case 1:
		offset += int64(r.off)
	case 2:
		offset += int64(len(r.b))
	default:
		err = syscall.EINVAL
		return
	}
	if offset < 0 {
		err = syscall.EINVAL
		return
	}
	if offset >= int64(len(r.b)) {
		r.off = len(r.b)
	} else {
		r.off = int(offset)
	}
	ret = int64(r.off)
	return
}

func (r *Reader) Read(val []byte) (n int, err error) {
	n = copy(val, r.b[r.off:])
	if n == 0 && len(val) != 0 {
		err = io.EOF
		return
	}
	r.off += n
	return
}

func (r *Reader) Close() (err error) {
	return
}

// ---------------------------------------------------

type Writer struct {
	b []byte
	n int
}

func NewWriter(buff []byte) *Writer {
	return &Writer{buff, 0}
}

func (p *Writer) Write(val []byte) (n int, err error) {
	n = len(val)
	cur, length := p.Cur(), p.Len()
	if cur-length > n {
		n = copy(p.b[p.n:], val)
	} else {
		if cur != length { // buff已经不够copy,多出来的buff[length:cur]要截掉
			p.b = p.b[:p.n]
			// fmt.Println("cur-len:", p.b[length:cur])
		}
		p.b = append(p.b, val...)
	}
	p.n += n
	// fmt.Printf("write:|%+v|%s|len:%d,Len:%d\n", val, val, n, p.n)
	if n == 0 && len(val) > 0 {
		err = io.EOF
		return
	}
	return
}

func (p *Writer) Len() int {
	return p.n
}

func (p *Writer) Cur() int {
	cur := len(p.b)
	// fmt.Println("cur:", cur)
	return cur
}

func (p *Writer) Bytes() []byte {
	// fmt.Println("Len", p.n)
	// fmt.Println(p.b[:p.n])
	return p.b[:p.n]
}

func (p *Writer) Reset() {
	p.n = 0
}

// ---------------------------------------------------

type Buffer struct {
	b []byte
}

func NewBuffer() *Buffer {
	return new(Buffer)
}

func (p *Buffer) ReadAt(buf []byte, off int64) (n int, err error) {
	ioff := int(off)
	if len(p.b) <= ioff {
		return 0, io.EOF
	}
	n = copy(buf, p.b[ioff:])
	if n != len(buf) {
		err = io.EOF
	}
	return
}

func (p *Buffer) WriteAt(buf []byte, off int64) (n int, err error) {
	ioff := int(off)
	iend := ioff + len(buf)
	if len(p.b) < iend {
		if len(p.b) == ioff {
			p.b = append(p.b, buf...)
			return len(buf), nil
		}
		zero := make([]byte, iend-len(p.b))
		p.b = append(p.b, zero...)
	}
	copy(p.b[ioff:], buf)
	return len(buf), nil
}

func (p *Buffer) WriteStringAt(buf string, off int64) (n int, err error) {
	ioff := int(off)
	iend := ioff + len(buf)
	if len(p.b) < iend {
		if len(p.b) == ioff {
			p.b = append(p.b, buf...)
			return len(buf), nil
		}
		zero := make([]byte, iend-len(p.b))
		p.b = append(p.b, zero...)
	}
	copy(p.b[ioff:], buf)
	return len(buf), nil
}

func (p *Buffer) Truncate(fsize int64) (err error) {
	size := int(fsize)
	if len(p.b) < size {
		zero := make([]byte, size-len(p.b))
		p.b = append(p.b, zero...)
	} else {
		p.b = p.b[:size]
	}
	return nil
}

func (p *Buffer) Buffer() []byte {
	return p.b
}

func (p *Buffer) Len() int {
	return len(p.b)
}

// ---------------------------------------------------
