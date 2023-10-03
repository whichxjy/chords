package model

import "github.com/charmbracelet/bubbles/list"

type ChordKind int

const (
	MajorChordKind ChordKind = iota
	MinorChordKind
)

func (k ChordKind) String() string {
	switch k {
	case MajorChordKind:
		return "Major"
	case MinorChordKind:
		return "Minor"
	}
	return "-"
}

func (k ChordKind) FilterValue() string {
	return k.String()
}

var chords = map[ChordKind]Chord{
	MajorChordKind: &MajorChord{},
	MinorChordKind: &MinorChord{},
}

func GetChord(kind ChordKind) Chord {
	return chords[kind]
}

func GetChordNotes(chord Chord, functions []*Note) []*Note {
	notes := chord.Pick(functions)
	return chord.Convert(notes)
}

type Chord interface {
	GetSymbol(note *Note) string
	Description() string
	Pick(functions []*Note) []*Note
	Convert(notes []*Note) []*Note
}

type MajorChord struct{}

func (c *MajorChord) GetSymbol(note *Note) string {
	return note.GetName()
}

func (c *MajorChord) Description() string {
	return "1 - 3 - 5"
}

func (c *MajorChord) Pick(functions []*Note) []*Note {
	return []*Note{functions[0], functions[2], functions[4]}
}

func (c *MajorChord) Convert(notes []*Note) []*Note {
	return notes
}

type MinorChord struct{}

func (c *MinorChord) GetSymbol(note *Note) string {
	return note.GetName() + "m"
}

func (c *MinorChord) Description() string {
	return "1 - 3b - 5"
}

func (c *MinorChord) Pick(functions []*Note) []*Note {
	return []*Note{functions[0], functions[2], functions[4]}
}

func (c *MinorChord) Convert(notes []*Note) []*Note {
	return []*Note{notes[0], notes[1].Flat(), notes[2]}
}

func GetChordKindListForUI() []list.Item {
	xs := make([]list.Item, 0, len(chords))
	for c := range chords {
		xs = append(xs, c)
	}
	return xs
}
