package scale

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

func makeFunctionRow() table.Row {
	row := table.Row{functionRowName}
	for i := 0; i < rowLength; i++ {
		row = append(row, makeFunction(i))
	}
	return row
}

func makeFunction(idx int) string {
	if idx%2 == 1 {
		return "   "
	}

	funcIdx := getFuncIdx(idx)
	if funcIdx == 8 {
		return "8/1"
	}
	return fmt.Sprintf("%d", funcIdx)
}

func getFuncIdx(idx int) int {
	return idx/2 + 1
}
