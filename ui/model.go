package ui

import (
	"bytes"
	"fmt"
	"strings"

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

const (
	KeyCtrlC = "ctrl+c"
	KeyUp    = "up"
	KeyDown  = "down"
	KeyEnter = "enter"
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
	m.noteList = list.New(getNoteListForUI(), noteDelegate{}, listWidth, listHeight)
	m.chordKindList = list.New(getChordKindListForUI(), chordDelegate{}, listWidth, listHeight)
	m.state = WaitNoteState
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.onKeyPressed(msg.String())
	case tea.WindowSizeMsg:
		m.noteList.SetWidth(msg.Width)
		m.chordKindList.SetWidth(msg.Width)
		return m, nil
	}
	return m, nil
}

func (m *Model) onKeyPressed(key string) tea.Cmd {
	if key == KeyCtrlC {
		return tea.Quit
	}

	switch m.state {
	case WaitNoteState:
		switch key {
		case KeyUp:
			m.noteList.CursorUp()
		case KeyDown:
			m.noteList.CursorDown()
		case KeyEnter:
			m.selectedNote = m.noteList.SelectedItem().(*model.Note)
			m.state = WaitChordState
		}
	case WaitChordState:
		switch key {
		case KeyUp:
			m.chordKindList.CursorUp()
		case KeyDown:
			m.chordKindList.CursorDown()
		case KeyEnter:
			m.selectedChordKind = m.chordKindList.SelectedItem().(model.ChordKind)
			m.state = ShowState
		}
	case ShowState:
		switch key {
		case KeyEnter:
			m.noteList.ResetSelected()
			m.chordKindList.ResetSelected()
			m.state = WaitNoteState
		}
	}
	return nil
}

func (m *Model) View() string {
	switch m.state {
	case WaitNoteState:
		return m.noteList.View()
	case WaitChordState:
		return m.chordKindList.View()
	case ShowState:
		return getChordView(m.selectedNote, m.selectedChordKind)
	}
	return ""
}

func getChordView(tonic *model.Note, chordKind model.ChordKind) string {
	var bf bytes.Buffer
	bf.WriteString(fmt.Sprintf("%s Major Scale:\n", tonic.GetName()))
	table, functions := scale.Make(tonic)
	chord := model.GetChord(chordKind)
	notes := model.GetChordNotes(chord, functions)
	bf.WriteString(table)
	bf.WriteString("\n\n")
	bf.WriteString(fmt.Sprintf("Symbol: %v\n", chord.GetSymbol(tonic)))
	bf.WriteString(fmt.Sprintf("Chord: %v\n", chord.Description()))
	bf.WriteString(fmt.Sprintf("Notes: %v\n", notesToView(notes)))
	bf.WriteString(fmt.Sprintf("Intervals: %v\n", notesWithIntervalToView(notes)))
	return bf.String()
}

func notesToView(notes []*model.Note) string {
	names := make([]string, 0, len(notes))
	for _, note := range notes {
		names = append(names, note.GetName())
	}
	return strings.Join(names, " - ")
}

func notesWithIntervalToView(notes []*model.Note) string {
	strs := make([]string, 0, len(notes))
	for i := 0; i < len(notes); i++ {
		currNote := notes[i]
		if i == 0 {
			strs = append(strs, currNote.GetName())
		} else {
			interval := model.GetNotesInterval(notes[i-1], currNote)
			strs = append(strs, fmt.Sprintf("[%v]", interval))
			strs = append(strs, currNote.GetName())
		}
	}
	return strings.Join(strs, " - ")
}
