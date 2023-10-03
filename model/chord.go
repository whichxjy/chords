package model

type Chord interface {
	Description() string
	Pick(functions []*Note) []*Note
	Convert(notes []*Note) []*Note
}

type MajorChord struct{}

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

func (c *MinorChord) Description() string {
	return "1 - 3b - 5"
}

func (c *MinorChord) Pick(functions []*Note) []*Note {
	return []*Note{functions[0], functions[2], functions[4]}
}

func (c *MinorChord) Convert(notes []*Note) []*Note {
	return []*Note{notes[0], notes[1].Flat(), notes[2]}
}
