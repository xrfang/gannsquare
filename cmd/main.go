package main

import (
	"encoding/json"
	"fmt"
	"gann"
	"os"
	"strconv"
)

func main() {
	sq := gann.NewSquare(1, 25, 0.25)
	sq.Dump(os.Stdout)
	if len(os.Args) > 1 {
		v, err := strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		t := sq.Locate(v)
		je := json.NewEncoder(os.Stdout)
		je.SetIndent("", "    ")
		je.Encode(t)
	}
}
