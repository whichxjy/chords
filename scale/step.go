package scale

import "github.com/jedib0t/go-pretty/v6/table"

func makeStepsRow() table.Row {
	row := table.Row{stepsRowName}
	for i := 0; i < rowLength; i++ {
		row = append(row, makeStep(i))
	}
	return row
}

func makeStep(idx int) string {
	if idx%2 == 0 {
		return "   "
	}
	if idx == 5 || idx == 13 {
		return "1/2"
	}
	return " 1 "
}
