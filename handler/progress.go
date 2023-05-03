package handler

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type Progress struct {
	stdout  io.Writer
	buf     bytes.Buffer
	prev    int
	curr    int
	bscount int
	start   bool
}

func (p *Progress) Show(curr int, total int) {
	var s string
	p.curr = curr * 50 / total
	if !p.start {
		fmt.Printf("[")
		p.start = true
		p.bscount = 0
	}
	if p.curr > p.prev {
		s = ""
		for i := 0; i < p.bscount; i++ {
			s += "\b"
		}
		s += "\u2588"
		for i := 0; i < (50 - p.curr); i++ {
			s += "."
		}
		s += "] %3d%%"
		p.bscount = 50 - p.curr + 6
		fmt.Printf(s, p.curr*2)
		p.prev = p.curr
	}
	if p.curr >= 50 {
		fmt.Printf("\n")
	}
}

func InitProgress() Progress {
	return Progress{
		stdout:  os.Stdout,
		buf:     bytes.Buffer{},
		prev:    0,
		curr:    0,
		bscount: 0,
		start:   false,
	}
}
