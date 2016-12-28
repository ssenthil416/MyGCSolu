package drum

import (
	"fmt"
)

//Pattern details
type Pattern struct {
	version string
	tempo   float32
	tracks  []Track
}

//Track structure.
type Track struct {
	id    byte
	name  []byte
	steps [16]byte
}

func (t Track) String() string {
	//Write the header.
	header := fmt.Sprintf("(%d) %s\t", t.id, t.name)

	//Write the steps.
	steps := []byte("|----|----|----|----|")

	//Add an 'x' for each note
	for i, x := range t.steps {
		if x == 1 {
			// need to adjust 'i' to account for the '|'s
			steps[i+i/4+1] = 'x'
		}
	}
	return header + string(steps)
}

func (p Pattern) String() string {
	//Write the header.
	str := fmt.Sprintf("Saved with HW Version: %s\nTempo: %g\n", p.version, p.tempo)

	//Write each track.
	for _, track := range p.tracks {
		str += fmt.Sprintln(track)
	}
	return str
}
