package scale

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/whichxjy/chords/model"
)

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
	noteIdx := (startIdx + stepSum) % len(model.Notes)
	return model.Notes[noteIdx].GetName()
}
