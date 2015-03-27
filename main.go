package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var rdr io.Reader
	repl := false
	switch len(os.Args) {
	case 1:
		rdr = os.Stdin
		repl = true
	case 2:
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		rdr = f
	default:
		fmt.Fprintln(os.Stderr, "incorrect program arguments")
		os.Exit(1)
	}
	pxis := NewPxiState()
	prdr := NewReader(rdr)
	for {
		v, err := prdr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		v, err = pxis.Eval(pxis.genv, v)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if repl {
			fmt.Println(v.String())
		}
	}
}
