package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
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

func NotesToView(notes []*Note) string {
	names := make([]string, 0, len(notes))
	for _, note := range notes {
		names = append(names, note.GetName())
	}
	return strings.Join(names, " - ")
}

func GetNoteListForUI() []list.Item {
	xs := make([]list.Item, 0, len(Notes))
	for _, x := range Notes {
		xs = append(xs, x)
	}
	return xs
}
