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

	state  State
	cursor int

	selectedNote      *model.Note
	selectedChordKind model.ChordKind
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
				m.selectedNote = model.GetNote("C")
				m.state = WaitChordState
			case WaitChordState:
				m.selectedChordKind = model.MajorChordKind
				m.state = ShowState
			}
		case "up":
			switch m.state {
			case WaitNoteState:
				m.noteList.CursorUp()
			case WaitChordState:
				m.chordKindList.CursorUp()
			}
		case "down":
			switch m.state {
			case WaitNoteState:
				m.noteList.CursorDown()
			case WaitChordState:
				m.chordKindList.CursorDown()
			}
		case "enter":
			switch m.state {
			case WaitNoteState:
				m.selectedNote = m.noteList.SelectedItem().(*model.Note)
				m.state = WaitChordState
			case WaitChordState:
				m.selectedChordKind = m.chordKindList.SelectedItem().(model.ChordKind)
				m.state = ShowState
			case ShowState:
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
		bf.WriteString(fmt.Sprintf("%s Major Scale:\n", m.selectedNote.GetName()))
		table, functions := scale.Make(m.selectedNote)
		chord := model.GetChord(m.selectedChordKind)
		notes := model.GetChordNotes(chord, functions)
		bf.WriteString(table)
		bf.WriteString("\n\n")
		bf.WriteString(fmt.Sprintf("Name: %v\n", chord.GetSymbol(m.selectedNote)))
		bf.WriteString(fmt.Sprintf("Chord: %v\n", chord.Description()))
		bf.WriteString(fmt.Sprintf("Notes: %v\n", model.NotesToView(notes)))
		return bf.String()
	}
	return ""
}
