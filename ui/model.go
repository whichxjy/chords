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
	WaitTonicState State = iota
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
	state State

	viewReady bool
	// list to select tonic
	noteList      list.Model
	selectedTonic *model.Note
	// list to select chord kind
	chordKindList list.Model
	// list of notes
	headerText string
	viewport   viewport.Model
}

func (m *Model) Init() tea.Cmd {
	m.state = WaitTonicState
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
	case WaitTonicState:
		switch key {
		case KeyUp:
			m.noteList.CursorUp()
		case KeyDown:
			m.noteList.CursorDown()
		case KeyEnter:
			m.selectedTonic = m.noteList.SelectedItem().(*model.Note)
			m.onTonicSelected(m.selectedTonic)
			m.state = WaitChordState
		}
	case WaitChordState:
		switch key {
		case KeyUp:
			m.chordKindList.CursorUp()
		case KeyDown:
			m.chordKindList.CursorDown()
		case KeyEnter:
			selectedChordKind := m.chordKindList.SelectedItem().(model.ChordKind)
			m.onChordKindSelected(selectedChordKind)
			m.state = ShowState
		}
	case ShowState:
		switch key {
		case KeyEnter:
			m.noteList.ResetSelected()
			m.chordKindList.ResetSelected()
			m.state = WaitTonicState
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
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	viewportHeight := msg.Height - (headerHeight + footerHeight)

	if !m.viewReady {
		m.noteList = newNoteList(msg.Width, msg.Height)
		m.chordKindList = newChordKindList(msg.Width, msg.Height)
		m.viewport = viewport.New(msg.Width, viewportHeight)
		m.resetViewportYPosition()
		m.viewReady = true
	} else {
		m.noteList.SetSize(msg.Width, msg.Height)
		m.chordKindList.SetSize(msg.Width, msg.Height)
		m.viewport.Width = msg.Width
		m.viewport.Height = viewportHeight
	}
}

func (m *Model) View() string {
	if !m.viewReady {
		return ""
	}

	switch m.state {
	case WaitTonicState:
		return m.noteList.View()
	case WaitChordState:
		return m.chordKindList.View()
	case ShowState:
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	}
	return ""
}

func (m *Model) onTonicSelected(tonic *model.Note) {
	setChordKindListTitle(&m.chordKindList, tonic)
}

func (m *Model) onChordKindSelected(chordKind model.ChordKind) {
	header, content := m.viewportView(chordKind)
	m.headerText = header
	m.viewport.SetContent(content)
	m.resetViewportYPosition()
}

func (m *Model) headerView() string {
	return titleStyle.Render("🎹" + m.headerText)
}

func (m *Model) footerView() string {
	return infoStyle.Render(fmt.Sprintf("📋%3.f%%", m.viewport.ScrollPercent()*100))
}

func (m *Model) viewportView(chordKind model.ChordKind) (string, string) {
	if chordKind == model.AllChorsKind {
		return "Chords", getAllChordsView(m.selectedTonic)
	}
	return "Chord", getSingleChordView(m.selectedTonic, chordKind)
}

func (m *Model) resetViewportYPosition() {
	headerHeight := lipgloss.Height(m.headerView())
	m.viewport.YOffset = headerHeight - 1
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
	bf.WriteString(fmt.Sprintf("%s Major Scale:\n", tonic.FullName()))
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
		names = append(names, note.FullName())
	}
	return strings.Join(names, " - ")
}

func notesWithIntervalToView(notes []*model.Note) string {
	strs := make([]string, 0, len(notes))
	for i := 0; i < len(notes); i++ {
		currNote := notes[i]
		if i == 0 {
			strs = append(strs, currNote.FullName())
		} else {
			interval := model.GetNotesInterval(notes[i-1], currNote)
			strs = append(strs, fmt.Sprintf("[%v]", interval))
			strs = append(strs, currNote.FullName())
		}
	}
	return strings.Join(strs, " - ")
}
