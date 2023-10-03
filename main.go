package main

import (
	"fmt"

	"github.com/whichxjy/chords/model"
	"github.com/whichxjy/chords/scale"
)

func main() {
	targetNote := "C"
	fmt.Printf("%s Major Scale:\n", targetNote)
	table, functions := scale.Make(targetNote)
	fmt.Println(table)

	chord := &model.MajorChord{}
	notes := getChordNotes(chord, functions)
	fmt.Printf("Description: %v\n", chord.Description())
	model.PrintNotes(notes)
}

func getChordNotes(chord model.Chord, functions []*model.Note) []*model.Note {
	notes := chord.Pick(functions)
	return chord.Convert(notes)
}
