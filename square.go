package gann

import (
	"fmt"
	"io"
	"math"
)

type (
	square [][]float64
	colors struct {
		Cross [2]string
		Point [2]string
	}
	target struct {
		Grid   float64 `json:"grid"`   //格子之间的距离（即步长）
		Expect float64 `json:"expect"` //需要查找的目标数字
		Actual float64 `json:"actual"` //格子中的实际数字
		Row    int     `json:"row"`    //目标所在行
		Col    int     `json:"col"`    //目标所在列
		Ring   int     `json:"ring"`   //目标所处层级（中心为0）
		Axis   float64 `json:"axis"`   //最近的辐射线值
		Gap    int     `json:"gap"`    //离辐射线的格数
		diff   float64 //计算过程中使用的内部变量
	}
)

var cs colors = colors{
	Cross: [2]string{"\033[1;32m", "\033[0m"},
	Point: [2]string{"\033[1;31m", "\033[0m"},
}

func (sq square) CellFormat(floatPoints int) string {
	var dw int
	for _, r := range sq {
		for _, c := range r {
			w := len(fmt.Sprintf("%d", int(c)))
			if w > dw {
				dw = w
			}
		}
	}
	return fmt.Sprintf("%%%d.%df", dw+floatPoints+1, floatPoints)
}

func (sq square) Locate(num float64) *target {
	min := 1e308
	max := -1e308
	c := len(sq) / 2
	t := target{Grid: sq[c-1][c] - sq[c][c], Expect: num, diff: 1e308}
	for i := 0; i < len(sq); i++ {
		for j := 0; j < len(sq); j++ {
			if sq[i][j] > max {
				max = sq[i][j]
			}
			if sq[i][j] < min {
				min = sq[i][j]
			}
			diff := math.Abs(sq[i][j] - num)
			if diff < t.diff {
				t.Row = i
				t.Col = j
				t.Actual = sq[i][j]
				t.diff = diff
			}
		}
	}
	if t.Expect < min-t.Grid || t.Expect > max+t.Grid {
		return nil
	}
	t.diff = 1e308
	onAxis := func(i, j int) bool {
		return i == j || i+j == len(sq)-1 || i == c || j == c
	}
	for k := 0; k < len(sq); k++ {
		if onAxis(k, t.Col) {
			diff := math.Abs(sq[k][t.Col] - t.Expect)
			if diff < t.diff {
				t.diff = diff
				t.Axis = sq[k][t.Col]
			}
		}
		if onAxis(t.Row, k) {
			diff := math.Abs(sq[t.Row][k] - t.Expect)
			if diff < t.diff {
				t.diff = diff
				t.Axis = sq[t.Row][k]
			}
		}
	}
	t.Gap = int(math.Abs((t.Axis-t.Actual)/t.Grid) + 0.5)
	ring := math.Max(math.Abs(float64(c-t.Row)), math.Abs(float64(c-t.Col)))
	t.Ring = int(ring + 0.5)
	return &t
}

func (sq square) Dump(w io.Writer) {
	cf := sq.CellFormat(3)
	center := len(sq) / 2
	tty := isTTY(w)
	for i, r := range sq {
		for j, c := range r {
			style := cf + " "
			if tty {
				onPoint := false
				onCross := i == j || i+j == len(sq)-1 || i == center || j == center
				if onPoint {
					style = fmt.Sprintf("%s%s%s ", cs.Point[0], cf, cs.Point[1])
				} else if onCross {
					style = fmt.Sprintf("%s%s%s ", cs.Cross[0], cf, cs.Cross[1])
				} else {
					style = fmt.Sprintf("%s ", cf)
				}
			}
			fmt.Fprintf(w, style, c)
		}
		fmt.Fprintln(w)
	}
}

func NewSquare(start, until, step float64) square {
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
