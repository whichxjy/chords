package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whichxjy/chords/model"
)

const (
	listWidth  int = 40
	listHeight int = 14
)

var (
	// Colors
	backgroundColor = lipgloss.Color("#107896")
	foregroundColor = lipgloss.Color("#FAFAFA")
	// List
	notSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(4)
	selectedStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(backgroundColor)
	// Display
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(foregroundColor).
			Background(backgroundColor)
	infoStyle = titleStyle.Copy()
)

func newNoteList() list.Model {
	noteList := list.New(getNoteListForUI(), noteDelegate{}, listWidth, listHeight)
	noteList.Title = "🎹Choose the tonic"
	setListStyle(&noteList)
	return noteList
}

func newChordKindList() list.Model {
	chordKindList := list.New(getChordKindListForUI(), chordDelegate{}, listWidth, listHeight)
	setListStyle(&chordKindList)
	return chordKindList
}

func setChordKindListTitle(l *list.Model, tonic *model.Note) {
	l.Title = fmt.Sprintf("🎹Choose the chord for tonic %v", tonic.FullName())
}

func getNoteListForUI() []list.Item {
	xs := make([]list.Item, 0, len(model.Notes))
	for _, x := range model.Notes {
		xs = append(xs, x)
	}
	return xs
}

func getChordKindListForUI() []list.Item {
	xs := make([]list.Item, 0, len(model.ChordKinds)+1)
	xs = append(xs, model.AllChorsKind)
	for _, chordKind := range model.OrderedChordKinds {
		xs = append(xs, chordKind)
	}
	return xs
}

func setListStyle(l *list.Model) {
	l.Styles.Title = l.Styles.Title.Background(backgroundColor)
}

type listDelegate struct{}

func (ld listDelegate) Height() int { return 1 }

func (ld listDelegate) Spacing() int { return 0 }

func (ld listDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

type noteDelegate struct {
	listDelegate
}

func (nd noteDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	note, ok := listItem.(*model.Note)
	if !ok {
		return
	}

	str := fmt.Sprintf("[%02d] %v", index+1, note.FullName())

	fn := notSelectedStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type chordDelegate struct {
	listDelegate
}

func (cd chordDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	chordKind, ok := listItem.(model.ChordKind)
	if !ok {
		return
	}

	str := fmt.Sprintf("[%02d] %v", index+1, chordKind.String())

	fn := notSelectedStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
