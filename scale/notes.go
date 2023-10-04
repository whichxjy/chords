package scale

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/whichxjy/chords/model"
)

func makeNotesRow(tonic *model.Note, functions *[]*model.Note) table.Row {
	row := table.Row{notesRowName}
	for i := 0; i < rowLength; i++ {
		row = append(row, makeNote(i, tonic.Idx, functions))
	}
	return row
}

func makeNote(idx, startIdx int, functions *[]*model.Note) string {
	if idx%2 == 1 {
		return "   "
	}
	funcIdx := getFuncIdx(idx)
	stepSum := getStepSum(funcIdx)
	noteIdx := (startIdx + stepSum) % len(model.Notes)
	note := model.Notes[noteIdx]
	*functions = append(*functions, note)
	return note.GetName()
}
