package main

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/whichxjy/chords/model"
	"github.com/whichxjy/chords/utils"
)

const (
	// row length
	rowLength = 15
	// row name
	stepsRowName    = "Steps"
	notesRowName    = "Notes"
	functionRowName = "Function"
)

var steps = []int{2, 2, 1, 2, 2, 2, 1}

func main() {
	symbol := "C"
	fmt.Printf("%s Major Scale:\n", symbol)

	tb := table.NewWriter()
	tb.SetOutputMirror(os.Stdout)
	tb.SetColumnConfigs(makeColumnConfigs())

	// Steps
	tb.AppendRow(makeStepRow())
	tb.AppendSeparator()
	// Notes
	tb.AppendRow(makeNotesRow(symbol))
	tb.AppendSeparator()
	// Function
	tb.AppendRow(makeFunctionRow())

	tb.Render()
}

func makeColumnConfigs() []table.ColumnConfig {
	configs := []table.ColumnConfig{
		{
			Number: 1,
			Colors: text.Colors{text.Bold},
		},
	}

	for i := 2; i < rowLength+2; i++ {
		configs = append(configs, table.ColumnConfig{
			Number: i,
			Align:  text.AlignCenter,
		})
	}

	return configs
}

func makeStepRow() table.Row {
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

func makeNotesRow(startNote string) table.Row {
	row := table.Row{notesRowName}
	for i := 0; i < rowLength; i++ {
		row = append(row, makeNote(i, model.GetNoteIdx(startNote)))
	}
	return row
}

func makeNote(idx, startIdx int) string {
	if idx%2 == 1 {
		return "   "
	}
	funcIdx := getFuncIdx(idx)
	stepSum := getStepSum(funcIdx)
	nodeIdx := (startIdx + stepSum) % len(model.Notes)
	return model.Notes[nodeIdx].GetName()
}

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

func getStepSum(targetFuncIdx int) int {
	if targetFuncIdx == 1 {
		return 0
	}
	return utils.IntSliceSum(steps[:targetFuncIdx-1])
}
