package scale

import (
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

// Result:
// 1. Rendered major scale table in string
// 2. Function map (function => note)
func Make(tonic *model.Note) (string, []*model.Note) {
	tb := table.NewWriter()
	tb.SetColumnConfigs(makeColumnConfigs())

	functions := make([]*model.Note, 0)

	// Steps
	tb.AppendRow(makeStepsRow())
	tb.AppendSeparator()
	// Notes
	tb.AppendRow(makeNotesRow(tonic, &functions))
	tb.AppendSeparator()
	// Function
	tb.AppendRow(makeFunctionRow())

	return tb.Render(), functions
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

func getStepSum(targetFuncIdx int) int {
	if targetFuncIdx == 1 {
		return 0
	}
	return utils.IntSliceSum(steps[:targetFuncIdx-1])
}
