package internal

import "fmt"

var (
	// ErrTimeout represents a timeout error
	ErrTimeout = fmt.Errorf("timeout")
	// ErrEOF represents a end of file error
	ErrEOF = fmt.Errorf("eof")
)
