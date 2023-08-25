package provider

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/vitaliy-ukiru/fsm-telebot/storages/file"
)

// Gob provides encoding/gob format.
// Zero value is safe.
type Gob struct{}

func NewGob() *Gob {
	return &Gob{}
}

func (Gob) ProviderName() string { return "gob" }

func (Gob) Encode(v interface{}) ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := gob.NewEncoder(buff).Encode(v); err != nil {
		return nil, newError("gob", "encode", err)
	}
	return buff.Bytes(), nil
}

func (Gob) Decode(data []byte, to interface{}) error {
	buff := bytes.NewReader(data)
	return newError("gob", "decode", gob.NewDecoder(buff).Decode(to))
}

func (Gob) Save(w io.Writer, data file.ChatsStorage) error {
	e := gob.NewEncoder(w)
	err := e.Encode(data)
	return newError("gob", "save", err)
}

func (g Gob) Read(r io.Reader) (file.ChatsStorage, error) {
	d := gob.NewDecoder(r)
	var dest file.ChatsStorage
	if err := d.Decode(&dest); err != nil {
		return nil, newError("gob", "read", err)
	}
	return dest, nil
}
