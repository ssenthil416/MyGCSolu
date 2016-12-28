package drum

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	//    "encoding/hex"
)

//DecodeFile read input and decode data.
func DecodeFile(fullfilename string) (*Pattern, error) {

	//Initialise Pattern struct.
	p := Pattern{}

	//Open input file with file path and name.
	file, err := os.Open(fullfilename)
	if err != nil {
		return &p, err
	}

	//Close file handle.
	defer file.Close()

	//Decode pattern data.
	if err := NewDecoder(file).Decode(&p); err != nil {
		return &p, err
	}

	//Retrun decoded data.
	return &p, nil
}

//Decoder structure.
type Decoder struct {
	r          io.Reader
	splice     [6]byte
	bytsToRead int64
	version    []byte
	bytsRead   int32
}

//NewDecoder constructor.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r, bytsRead: 0, bytsToRead: 0}
}

//Decode method to decode.
func (d *Decoder) Decode(p *Pattern) error {

	//Read Binary file for Splice field.
	if err := binary.Read(d.r, binary.BigEndian, &d.splice); err != nil {
		return errors.New("could not decode header: " + err.Error())
	}

	//Check SPLICE bytes.
	if string(d.splice[:]) != "SPLICE" {
		return errors.New("header missing SPLICE bytes")
	}

	//fmt.Println("d.splice =", d.splice)
	//Read bytes to Read.
	if err := binary.Read(d.r, binary.BigEndian, &d.bytsToRead); err != nil {
		return errors.New("could not decode header: " + err.Error())
	}

	//fmt.Println("d.bytsToRead =", d.bytsToRead)
	d.version = make([]byte, 32)
	//Read vesrion detail.
	if err := binary.Read(d.r, binary.BigEndian, &d.version); err != nil {
		return errors.New("could not decode header: " + err.Error())
	}

	//fmt.Println("d.version =",d.version)
	p.version = fmt.Sprintf("%s", bytes.Trim(d.version, "\x00"))
	d.bytsRead += 32

	//Read temo details.
	if err := binary.Read(d.r, binary.LittleEndian, &p.tempo); err != nil {
		return errors.New("could not decode header: " + err.Error())
	}

	//fmt.Println(p.tempo)
	d.bytsRead += 4

	// Now starts the track information, at 55th byte precisely.
	var t Track
	var len int32

	//Get tempo tracks
	for {
		//Check for end.
		if d.bytsRead == int32(d.bytsToRead) {
			break
		}

		//Read Tempo id.
		if err := binary.Read(d.r, binary.BigEndian, &t.id); err != nil {
			return errors.New("t.id could not decode : " + err.Error())
		}

		//fmt.Println("t.id =", t.id)

		d.bytsRead++

		//Read tempo len.
		if err := binary.Read(d.r, binary.BigEndian, &len); err != nil {
			return errors.New("len could not decode : " + err.Error())
		}

		//fmt.Println("len =", len)
		d.bytsRead += 4

		t.name = make([]byte, len)
		//Read tempo name
		if _, err := io.ReadFull(d.r, t.name); err != nil {
			return errors.New("t.name could not decode : " + err.Error())
		}

		//fmt.Printf("t.name = %s\n", t.name)
		d.bytsRead += len

		//Read tempo steps
		if err := binary.Read(d.r, binary.BigEndian, &t.steps); err != nil {
			return errors.New("t.steps could not decode : " + err.Error())
		}
		d.bytsRead += 16

		//Append all tracks of tempo.
		p.tracks = append(p.tracks, t)
	}

	return nil
}
