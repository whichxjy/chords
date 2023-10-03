package main

import (
	"fmt"
	"sort"

	"github.com/whichxjy/chords/model"
	"github.com/whichxjy/chords/scale"
)

func main() {
	targetNote := "C"
	fmt.Printf("%s Major Scale:\n", targetNote)
	table, functions := scale.Make(targetNote)
	fmt.Println(table)

	chord := &model.MinorChord{}
	notes := getChordNotes(chord, functions)
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Idx < notes[j].Idx
	})
	fmt.Printf("Description: %v\n", chord.Description())
	model.PrintNotes(notes)
}

func getChordNotes(chord model.Chord, functions []*model.Note) []*model.Note {
	notes := chord.Pick(functions)
	return chord.Convert(notes)
}
