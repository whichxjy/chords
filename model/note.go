package model

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Note struct {
	Name      string
	OtherName string
	Idx       int
}

var Notes = []*Note{
	newNote(0, "C"),
	newNote(1, "C#", "Db"),
	newNote(2, "D"),
	newNote(3, "D#", "Eb"),
	newNote(4, "E"),
	newNote(5, "F"),
	newNote(6, "F#", "Gb"),
	newNote(7, "G"),
	newNote(8, "G#", "Ab"),
	newNote(9, "A"),
	newNote(10, "A#", "Bb"),
	newNote(11, "B"),
}

func newNote(idx int, names ...string) *Note {
	nameNum := len(names)
	if nameNum < 1 || nameNum > 2 {
		panic("invalid name number")
	}
	n := &Note{
		Idx:  idx,
		Name: names[0],
	}
	if nameNum == 2 {
		n.OtherName = names[1]
	}
	return n
}

func (n *Note) GetName() string {
	if n.OtherName == "" {
		return n.Name
	}
	return n.Name + "/" + n.OtherName
}

func (n *Note) Flat() *Note {
	idx := (n.Idx - 1 + len(Notes)) % len(Notes)
	return Notes[idx]
}

func (n *Note) FilterValue() string {
	return n.GetName()
}

func GetNoteIdx(name string) int {
	for _, note := range Notes {
		if note.Name == name || note.OtherName == name {
			return note.Idx
		}
	}
	panic(fmt.Sprintf("invalid note name: %v", name))
}

func GetNote(name string) *Note {
	return Notes[GetNoteIdx(name)]
}

func NotesToView(notes []*Note) string {
	names := make([]string, 0, len(notes))
	for _, note := range notes {
		names = append(names, note.GetName())
	}
	return strings.Join(names, " - ")
}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type NoteDelegate struct{}

func (d NoteDelegate) Height() int { return 1 }

func (d NoteDelegate) Spacing() int { return 0 }

func (d NoteDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d NoteDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	note, ok := listItem.(*Note)
	if !ok {
		return
	}

	str := fmt.Sprintf("%v. %v", index+1, note.GetName())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func GetNotesListForUI() []list.Item {
	xs := make([]list.Item, 0, len(Notes))
	for _, x := range Notes {
		xs = append(xs, x)
	}
	return xs
}
