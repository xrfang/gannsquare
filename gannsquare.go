package gannsquare

import (
	"fmt"
	"io"
	"math"
)

type square [][]float64

func (sq square) Dump(w io.Writer) {
	dw := len(fmt.Sprintf("%d", int(sq[0][0])))
	style := fmt.Sprintf("%%%d.3f ", dw+4)
	for _, r := range sq {
		for _, c := range r {
			fmt.Fprintf(w, style, c)
		}
		fmt.Fprintln(w)
	}
}

func New(start, until, step float64) square {
	if start >= until || start < 0 || step <= 0 {
		panic("invalid argument")
	}
	items := int((until-start)/step+0.5) + 1
	size := int(math.Ceil(math.Sqrt(float64(items))))
	if size%2 == 0 {
		size++
	}
	var sq square
	for i := 0; i < size; i++ {
		var r []float64
		for j := 0; j < size; j++ {
			r = append(r, -1)
		}
		sq = append(sq, r)
	}
	x := size / 2
	y := x
	dir := 0 //行走方向: 0=上；1=右；2=下；3=左
	next := func(value float64) bool {
		sq[y][x] = value
		switch dir {
		case 0:
			y--
			if y < 0 {
				return false
			}
			if sq[y][x+1] < 0 {
				dir = 1
			}
		case 1:
			x++
			if x >= size {
				return false
			}
			if sq[y+1][x] < 0 {
				dir = 2
			}
		case 2:
			y++
			if y >= size {
				return false
			}
			if sq[y][x-1] < 0 {
				dir = 3
			}
		case 3:
			x--
			if x < 0 {
				return false
			}
			if sq[y-1][x] < 0 {
				dir = 0
			}
		}
		return true
	}
	current := start //初始值
	for next(current) {
		current += step
	}
	return sq
}
