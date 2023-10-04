package ui

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whichxjy/chords/model"
)

const (
	listWidth  int = 20
	listHeight int = 14
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

func getNoteListForUI() []list.Item {
	xs := make([]list.Item, 0, len(model.Notes))
	for _, x := range model.Notes {
		xs = append(xs, x)
	}
	return xs
}

func getChordKindListForUI() []list.Item {
	xs := make([]list.Item, 0, len(model.Chords))
	for c := range model.Chords {
		xs = append(xs, c)
	}
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].(model.ChordKind) < xs[j].(model.ChordKind)
	})
	return xs
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

	str := fmt.Sprintf("[%02d] %v", index+1, note.GetName())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
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

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}