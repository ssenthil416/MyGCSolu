package drum

import (
	"fmt"
)

type Pattern struct {
	version string
	tempo   float32
	tracks  []Track
}

type Track struct {
	id    byte
	name  []byte
	steps [16]byte
}

func (t Track) String() string {
	// write the header
	header := fmt.Sprintf("(%d) %s\t", t.id, t.name)
	// write the steps
	steps := []byte("|----|----|----|----|")
	// add an 'x' for each note
	for i, x := range t.steps {
		if x == 1 {
			// need to adjust 'i' to account for the '|'s
			steps[i+i/4+1] = 'x'
		}
	}
	return header + string(steps)
}

func (p Pattern) String() string {
	// write the header
	str := fmt.Sprintf("Saved with HW Version: %s\nTempo: %g\n", p.version, p.tempo)
	// write each track
	for _, track := range p.tracks {
		str += fmt.Sprintln(track)
	}

	return str
}
