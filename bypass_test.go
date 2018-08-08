package wsbridge

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

type WriteCloser struct {
	*bytes.Buffer
}

func (wc *WriteCloser) Close() error { return nil }

type fakeReadable struct {
	Readable
	fakeNextReader func() (int, io.Reader, error)
}

func (rm fakeReadable) NextReader() (int, io.Reader, error) {
	return rm.fakeNextReader()
}

type fakeWritable struct {
	Writable
	fakeNextWriter func(int) (io.WriteCloser, error)
}

func (wm fakeWritable) NextWriter(mt int) (io.WriteCloser, error) {
	return wm.fakeNextWriter(mt)
}

func TestBypassSuccess1(t *testing.T) {
	const expect = "Hello World\n"
	r := bytes.NewBufferString(expect)
	w := &WriteCloser{new(bytes.Buffer)}
	ri := fakeReadable{
		fakeNextReader: func() (int, io.Reader, error) {
			return 0, r, nil
		},
	}
	wi := fakeWritable{
		fakeNextWriter: func(int) (io.WriteCloser, error) {
			return w, nil
		},
	}
	if err := bypass(ri, wi); err != nil {
		t.Fatalf("Bypass error: %v", err)
	}
	if w.String() != expect {
		t.Fatalf("'%s' != '%s'", w.String(), expect)
	}
}

func TestBypassSuccess2(t *testing.T) {
	const expect = "Hello World\n"
	r := io.MultiReader(
		bytes.NewBufferString(expect),
		bytes.NewBufferString(expect),
		bytes.NewBufferString(expect),
		bytes.NewBufferString(expect),
		bytes.NewBufferString(expect),
	)
	w := &WriteCloser{new(bytes.Buffer)}
	ri := fakeReadable{
		fakeNextReader: func() (int, io.Reader, error) {
			return 0, r, nil
		},
	}
	wi := fakeWritable{
		fakeNextWriter: func(int) (io.WriteCloser, error) {
			return w, nil
		},
	}
	if err := bypass(ri, wi); err != nil {
		t.Fatalf("Bypass error: %v", err)
	}
	if w.String() != strings.Repeat(expect, 5) {
		t.Fatalf("'%s' != '%s'", w.String(), expect)
	}
}

func TestBypassFail1(t *testing.T) {
	var expect = fmt.Errorf("Error")
	ri := fakeReadable{
		fakeNextReader: func() (int, io.Reader, error) {
			return 0, nil, expect
		},
	}
	wi := fakeWritable{
		fakeNextWriter: func(int) (io.WriteCloser, error) {
			return nil, nil
		},
	}
	err := bypass(ri, wi)
	if err.Error() != "Failed to get next reader: Error" {
		t.Fatalf("Expected error is not retuned: %s", err)
	}
}

func TestBypassFail2(t *testing.T) {
	var expect = fmt.Errorf("Error")
	r := bytes.NewBufferString("Hello World")
	ri := fakeReadable{
		fakeNextReader: func() (int, io.Reader, error) {
			return 0, r, nil
		},
	}
	wi := fakeWritable{
		fakeNextWriter: func(int) (io.WriteCloser, error) {
			return nil, expect
		},
	}
	err := bypass(ri, wi)
	if err.Error() != "Failed to get next writer (0): Error" {
		t.Fatalf("Expected error is not retuned: %s", err)
	}
}
