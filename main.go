package main

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

type Note struct {
	Name      string
	OtherName string
	Idx       int
}

var (
	notes = []*Note{
		newNote(0, "C"),
		newNote(1, "C#", "Db"),
		newNote(2, "D"),
		newNote(3, "D#", "Eb"),
		newNote(4, "E"),
		newNote(5, "F"),
		newNote(6, "F#", "Gb"),
		newNote(7, "G"),
		newNote(8, "G#", "Ab"),
		newNote(9, "A"),
		newNote(10, "A#", "Bb"),
		newNote(11, "B"),
	}
	steps = []int{2, 2, 1, 2, 2, 2, 1}
)

func newNote(idx int, names ...string) *Note {
	nameNum := len(names)
	if nameNum < 1 || nameNum > 2 {
		panic("invalid name number")
	}
	n := &Note{
		Idx:  idx,
		Name: names[0],
	}
	if nameNum == 2 {
		n.OtherName = names[1]
	}
	return n
}

func (n *Note) GetName() string {
	if n.OtherName == "" {
		return n.Name
	}
	return n.Name + "/" + n.OtherName
}

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
		row = append(row, makeNote(i, getNoteIdx(startNote)))
	}
	return row
}

func makeNote(idx, startIdx int) string {
	if idx%2 == 1 {
		return "   "
	}
	funcIdx := getFuncIdx(idx)
	stepSum := getStepSum(funcIdx)
	nodeIdx := (startIdx + stepSum) % len(notes)
	return notes[nodeIdx].GetName()
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

func getNoteIdx(targetNote string) int {
	for _, note := range notes {
		if note.Name == targetNote || note.OtherName == targetNote {
			return note.Idx
		}
	}
	return -1
}
