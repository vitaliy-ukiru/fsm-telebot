package provider

import (
	"encoding/json"
	"io"

	"github.com/vitaliy-ukiru/fsm-telebot/storages/file"
)

// JsonSettings configures json encoder and decoder.
//
// Zero value configures as default json.Encoder and json.Decoder.
type JsonSettings struct {
	Prefix                string
	Indent                string
	UseNumber             bool
	DisallowUnknownFields bool
}

func (js JsonSettings) setDecoder(d *json.Decoder) {
	if js.UseNumber {
		d.UseNumber()
	}

	if js.DisallowUnknownFields {
		d.DisallowUnknownFields()
	}
}

// Json provides json format.
//
// // Unexported fields will be ignoring (json package behavior).
type Json struct {
	JsonSettings
}

func NewJson(jsonSettings JsonSettings) *Json {
	return &Json{JsonSettings: jsonSettings}
}

func (j Json) Encode(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, newError("json", "encode", err)
	}
	return data, nil
}

func (j Json) Decode(data []byte, v any) error {
	return newError("json", "decode", json.Unmarshal(data, v))
}

func (j Json) ProviderName() string { return "json" }

func (j Json) Save(w io.Writer, data file.ChatsStorage) error {
	e := json.NewEncoder(w)
	e.SetIndent(j.Prefix, j.Indent)

	err := e.Encode(data)
	return newError("json", "save", err)
}
func (j Json) Read(r io.Reader) (file.ChatsStorage, error) {
	d := json.NewDecoder(r)
	j.setDecoder(d)
	var dest file.ChatsStorage
	if err := d.Decode(&dest); err != nil {
		return nil, newError("json", "read", err)
	}
	return dest, nil
}
