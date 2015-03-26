package main

import (
	"fmt"
	"os"
)

func main() {
	pxis := NewPxiState()
	for {
		v, err := Read(os.Stdin)
		if err != nil {
			fmt.Println(err)
			return
		}
		v, err = pxis.Eval(pxis.genv, v)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(v.String())
	}
}
