package gann

import (
	"io"
	"os"

	"golang.org/x/sys/unix"
)

func isTTY(w io.Writer) bool {
	if w == os.Stdout {
		_, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
		return err == nil
	}
	return false
}
