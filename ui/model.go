package ui

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/whichxjy/chords/model"
	"github.com/whichxjy/chords/scale"
)

type State int

const (
	WaitNoteState State = iota
	WaitChordState
	ShowState
	QuittingState
)

type Model struct {
	noteList      list.Model
	chordKindList list.Model

	state State

	note      *model.Note
	chordKind model.ChordKind
}

func (m *Model) Init() tea.Cmd {
	m.noteList = list.New(model.GetNoteListForUI(), noteDelegate{}, 20, 14)
	m.chordKindList = list.New(model.GetChordKindListForUI(), chordDelegate{}, 20, 14)
	m.state = WaitNoteState
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "1":
			switch m.state {
			case WaitNoteState:
				m.note = model.GetNote("C")
				m.state = WaitChordState
			case WaitChordState:
				m.chordKind = model.MajorChordKind
				m.state = ShowState
			}
		case "enter":
			if m.state == ShowState {
				m.state = WaitNoteState
			}
		}
	}
	return m, nil
}

func (m *Model) View() string {
	switch m.state {
	case WaitNoteState:
		return "\n" + m.noteList.View()
	case WaitChordState:
		return "\n" + m.chordKindList.View()
	case ShowState:
		var bf bytes.Buffer
		bf.WriteString(fmt.Sprintf("%s Major Scale:\n", "C"))
		table, functions := scale.Make("C")
		chord := model.GetChord(m.chordKind)
		notes := model.GetChordNotes(chord, functions)
		bf.WriteString(table)
		bf.WriteString("\n")
		bf.WriteString(model.NotesToView(notes))
		return bf.String()
	}
	return ""
}
