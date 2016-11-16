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

type Decoder struct {
	r          io.Reader
	splice     [6]byte
	bytsToRead int64
	version    []byte
	bytsRead   int32
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r, bytsRead: 0, bytsToRead: 0}
}

func DecodeFile(fullfilename string) (*Pattern, error) {

	p := &Pattern{}
	file, err := os.Open(fullfilename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	err = NewDecoder(file).Decode(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (d *Decoder) Decode(p *Pattern) error {

	err := binary.Read(d.r, binary.BigEndian, &d.splice)
	if err != nil {
		return errors.New("could not decode header: " + err.Error())
	}

	// check SPLICE bytes
	if string(d.splice[:]) != "SPLICE" {
		return errors.New("header missing SPLICE bytes")
	}

	//fmt.Println("d.splice =", d.splice)
	err = binary.Read(d.r, binary.BigEndian, &d.bytsToRead)
	if err != nil {
		return errors.New("could not decode header: " + err.Error())
	}

	//fmt.Println("d.bytsToRead =", d.bytsToRead)
	d.version = make([]byte, 32)
	err = binary.Read(d.r, binary.BigEndian, &d.version)
	if err != nil {
		return errors.New("could not decode header: " + err.Error())
	}

	//fmt.Println("d.version =",d.version)
	p.version = fmt.Sprintf("%s", bytes.Trim(d.version, "\x00"))
	d.bytsRead += 32

	if err = binary.Read(d.r, binary.LittleEndian, &p.tempo); err != nil {
		return errors.New("could not decode header: " + err.Error())
	}

	//fmt.Println(p.tempo)
	d.bytsRead += 4

	// Now starts the track information, at 55th byte precisely.
	var t Track
	var len int32

	for {
		/*Check for end*/
		if d.bytsRead == int32(d.bytsToRead) {
			break
		}

		if err := binary.Read(d.r, binary.BigEndian, &t.id); err != nil {
			return errors.New("t.id could not decode : " + err.Error())
		}

		//fmt.Println("t.id =", t.id)

		d.bytsRead += 1

		if err := binary.Read(d.r, binary.BigEndian, &len); err != nil {
			return errors.New("len could not decode : " + err.Error())
		}

		//fmt.Println("len =", len)
		d.bytsRead += 4

		t.name = make([]byte, len)
		if _, err := io.ReadFull(d.r, t.name); err != nil {
			return errors.New("t.name could not decode : " + err.Error())
		}

		//fmt.Printf("t.name = %s\n", t.name)
		d.bytsRead += len

		if err := binary.Read(d.r, binary.BigEndian, &t.steps); err != nil {
			return errors.New("t.steps could not decode : " + err.Error())
		}
		d.bytsRead += 16

		p.tracks = append(p.tracks, t)
	}

	return nil
}
