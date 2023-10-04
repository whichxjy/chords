package ui

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whichxjy/chords/model"
	"github.com/whichxjy/chords/scale"
)

// WaitNoteState -> WaitChordState -> ShowState -> WaitNoteState
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

	viewport      viewport.Model
	viewportReady bool
	headerText    string
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
		return m, m.onKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.onWindowSizeMsg(msg)
		return m, nil
	}
	return m, nil
}

func (m *Model) onKeyMsg(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()
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
			if m.selectedChordKind == model.AllChorsKind {
				m.viewport.SetContent(getAllChordsView(m.selectedNote))
				m.headerText = "Chords"
			} else {
				m.viewport.SetContent(getSingleChordView(m.selectedNote, m.selectedChordKind))
				m.headerText = "Chord"
			}
			m.state = ShowState
		}
	case ShowState:
		switch key {
		case KeyEnter:
			m.noteList.ResetSelected()
			m.chordKindList.ResetSelected()
			m.state = WaitNoteState
		default:
			// Let viewport to handle message
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return cmd
		}
	}
	return nil
}

func (m *Model) onWindowSizeMsg(msg tea.WindowSizeMsg) {
	m.noteList.SetWidth(msg.Width)
	m.chordKindList.SetWidth(msg.Width)

	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	if !m.viewportReady {
		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewportReady = true
	} else {
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
	}
}

func (m *Model) View() string {
	switch m.state {
	case WaitNoteState:
		return m.noteList.View()
	case WaitChordState:
		return m.chordKindList.View()
	case ShowState:
		if m.viewportReady {
			return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
		}
	}
	return ""
}

func (m *Model) headerView() string {
	return titleStyle.Render("🎹" + m.headerText)
}

func (m *Model) footerView() string {
	return infoStyle.Render(fmt.Sprintf("📋%3.f%%", m.viewport.ScrollPercent()*100))
}

func getAllChordsView(tonic *model.Note) string {
	var bf bytes.Buffer
	tableView, functions := getScaleTableView(tonic)
	bf.WriteString(tableView)
	bf.WriteString("\n\n")
	for _, chordKind := range model.OrderedChordKinds {
		bf.WriteString(getChordDetailView(tonic, chordKind, functions))
		bf.WriteString("\n")
	}
	return bf.String()
}

func getSingleChordView(tonic *model.Note, chordKind model.ChordKind) string {
	var bf bytes.Buffer
	tableView, functions := getScaleTableView(tonic)
	bf.WriteString(tableView)
	bf.WriteString("\n\n")
	chordView := getChordDetailView(tonic, chordKind, functions)
	bf.WriteString(chordView)
	return bf.String()
}

func getScaleTableView(tonic *model.Note) (string, []*model.Note) {
	var bf bytes.Buffer
	table, functions := scale.Make(tonic)
	bf.WriteString(fmt.Sprintf("%s Major Scale:\n", tonic.GetName()))
	bf.WriteString(table)
	return bf.String(), functions
}

func getChordDetailView(tonic *model.Note, chordKind model.ChordKind, functions []*model.Note) string {
	var bf bytes.Buffer
	chord := model.GetChord(chordKind)
	notes := model.GetChordNotes(chord, functions)
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
