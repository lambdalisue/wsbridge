package wsbridge

import (
	"fmt"
	"io"
)

// Readable inteface which indicates incoming websocket.Conn
type Readable interface {
	NextReader() (messageType int, r io.Reader, err error)
}

// Writable interface which indicates outgoing websocket.Conn
type Writable interface {
	NextWriter(messageType int) (w io.WriteCloser, err error)
}

/* Bypass any data from Readable to Writable
 *
 * It copies any data; including websocket's control messages; received
 * from the readable (ri) to the writable (wi).
 *
 * Use this function in an infinity-loop to bypass all data.
 */
func bypass(ri Readable, wi Writable) error {
	mt, r, err := ri.NextReader()
	if err != nil {
		return fmt.Errorf("Failed to get next reader: %s", err)
	}
	w, err := wi.NextWriter(mt)
	if err != nil {
		return fmt.Errorf("Failed to get next writer (%d): %s", mt, err)
	}
	if _, err := io.Copy(w, r); err != nil {
		return fmt.Errorf("Failed to copy data: %s", err)
	}
	// NOTE:
	//
	// Example in documentation and internal implementations of gorilla/websocket
	// does not use 'defer w.Close()' to make sure that the writer has closed.
	// This is because that w.Close() is used to flush the complete message to the
	// network and it doesn't make sense for a broken writer.
	//
	// https://godoc.org/github.com/gorilla/websocket#Conn.NextWriter
	// https://github.com/gorilla/websocket/blob/master/conn.go#L697-L706
	if err := w.Close(); err != nil {
		return fmt.Errorf("Failed to close writer: %s", err)
	}
	return nil
}
