package main

import (
	"fmt"

	"github.com/whichxjy/chords/scale"
)

func main() {
	symbol := "C"
	fmt.Printf("%s Major Scale:\n", symbol)
	table, _ := scale.Make(symbol)
	fmt.Println(table)
}
