package filesystem

import (
	"bufio"
	"fmt"
	"os"
)

// For object-storage backend
type fileWriter struct {
	file      *os.File
	size      int64
	bw        *bufio.Writer
	closed    bool
	committed bool
	cancelled bool
}

func (w *fileWriter) Write(p []byte) (int, error) {
	if w.closed {
		return 0, fmt.Errorf("already closed")
	} else if w.committed {
		return 0, fmt.Errorf("already committed")
	} else if w.cancelled {
		return 0, fmt.Errorf("already cancelled")
	}
	n, err := w.bw.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *fileWriter) Size() int64 {
	return w.size
}

func (w *fileWriter) Close() error {
	if w.closed {
		return fmt.Errorf("already closed")
	}

	if err := w.bw.Flush(); err != nil {
		return err
	}

	if err := w.file.Sync(); err != nil {
		return err
	}

	if err := w.file.Close(); err != nil {
		return err
	}
	w.closed = true
	return nil
}

func (w *fileWriter) Cancel() error {
	if w.closed {
		return fmt.Errorf("already closed")
	}

	w.cancelled = true
	w.file.Close()
	return os.Remove(w.file.Name())
}

func (w *fileWriter) Commit() error {
	if w.closed {
		return fmt.Errorf("already closed")
	} else if w.committed {
		return fmt.Errorf("already committed")
	} else if w.cancelled {
		return fmt.Errorf("already cancelled")
	}

	if err := w.bw.Flush(); err != nil {
		return err
	}

	if err := w.file.Sync(); err != nil {
		return err
	}

	w.committed = true
	return nil
}
