package main

import (
	"gann"
	"os"
)

func main() {
	sq := gann.NewSquare(1, 15, 0.125)
	sq.Dump(os.Stdout)
}
