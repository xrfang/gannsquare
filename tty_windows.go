package gann

import (
	"fmt"
	"io"
)

func isTTY(w io.Writer) bool {
	return false
}
