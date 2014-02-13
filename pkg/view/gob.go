package view

import (
	"bytes"
	"encoding/gob"
)

func (v *View) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(&v.members)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v *View) GobDecode(b []byte) error {
	buf := bytes.NewReader(b)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&v.members)
	if err != nil {
		return err
	}

	return nil
}
