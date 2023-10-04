package model

type ChordKind int

const (
	MajorChordKind ChordKind = iota
	MinorChordKind
	Sus2ChordKind
	Sus4ChordKind
)

func (k ChordKind) String() string {
	switch k {
	case MajorChordKind:
		return "Major"
	case MinorChordKind:
		return "Minor"
	case Sus2ChordKind:
		return "Sus2"
	case Sus4ChordKind:
		return "Sus4"
	}
	return "-"
}

func (k ChordKind) FilterValue() string {
	return k.String()
}

var Chords = map[ChordKind]Chord{
	MajorChordKind: &MajorChord{},
	MinorChordKind: &MinorChord{},
	Sus2ChordKind:  &Sus2Chord{},
	Sus4ChordKind:  &Sus4Chord{},
}

func GetChord(kind ChordKind) Chord {
	return Chords[kind]
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
	return "1 - b3 - 5"
}

func (c *MinorChord) Pick(functions []*Note) []*Note {
	return []*Note{functions[0], functions[2], functions[4]}
}

func (c *MinorChord) Convert(notes []*Note) []*Note {
	return []*Note{notes[0], notes[1].Flat(), notes[2]}
}

type Sus2Chord struct{}

func (c *Sus2Chord) GetSymbol(note *Note) string {
	return note.GetName() + "sus2"
}

func (c *Sus2Chord) Description() string {
	return "1 - 2 - 5"
}

func (c *Sus2Chord) Pick(functions []*Note) []*Note {
	return []*Note{functions[0], functions[1], functions[4]}
}

func (c *Sus2Chord) Convert(notes []*Note) []*Note {
	return []*Note{notes[0], notes[1], notes[2]}
}

type Sus4Chord struct{}

func (c *Sus4Chord) GetSymbol(note *Note) string {
	return note.GetName() + "sus4"
}

func (c *Sus4Chord) Description() string {
	return "1 - 4 - 5"
}

func (c *Sus4Chord) Pick(functions []*Note) []*Note {
	return []*Note{functions[0], functions[3], functions[4]}
}

func (c *Sus4Chord) Convert(notes []*Note) []*Note {
	return []*Note{notes[0], notes[1], notes[2]}
}
