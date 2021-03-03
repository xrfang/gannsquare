package main

import (
	"gannsquare"
	"os"
)

func main() {
	sq := gannsquare.New(1, 15, 0.125)
	sq.Dump(os.Stdout)
}
