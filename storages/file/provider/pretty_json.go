package provider

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"io"

	"github.com/vitaliy-ukiru/fsm-telebot/storages/file"
)

// PrettyJson provides json format with pretty encoding data values (file.Record Data fields).
//
// Default json package encodes []byte as base64 string.
// This provider allows use json.RawMessage. But it's not free.
// The structure is copied to the new one to keep the data safe.
//
// Unexported fields will be ignoring (json package behavior).
type PrettyJson struct {
	JsonSettings

	// TryDecodeBase64String tries to decode any strings values from base64.
	// (try to backward compatibility with Json).
	//
	// But the package does not take responsibility for decoding strings
	// if you use base64 for your own purposes.
	TryDecodeBase64String bool

	// IndentInEncodeMethod adds indent in Encode.
	// If Indent is set when Save is called, the indent
	// will be added regardless of this parameter.
	IndentInEncodeMethod bool
}

func NewPrettyJson(jsonSettings JsonSettings, tryDecodeBase64String bool, indentInEncodeMethod bool) *PrettyJson {
	return &PrettyJson{JsonSettings: jsonSettings, TryDecodeBase64String: tryDecodeBase64String, IndentInEncodeMethod: indentInEncodeMethod}
}

const prettyJson = "pretty_json"

func (j PrettyJson) ProviderName() string { return prettyJson }

func (j PrettyJson) Encode(v any) ([]byte, error) {
	buff := new(bytes.Buffer)

	e := json.NewEncoder(buff)
	if j.IndentInEncodeMethod {
		e.SetIndent("", j.Indent)
	}

	if err := e.Encode(v); err != nil {
		return nil, newError(prettyJson, "encode", err)
	}
	return buff.Bytes(), nil
}

func (j PrettyJson) Decode(data []byte, v any) error {
	buff := bytes.NewReader(data)

	d := json.NewDecoder(buff)
	if j.DisallowUnknownFields {
		d.DisallowUnknownFields()
	}

	if j.UseNumber {
		d.UseNumber()
	}
	return newError(prettyJson, "decode", d.Decode(v))
}

func (j PrettyJson) Save(w io.Writer, data file.ChatsStorage) error {
	buff := new(bytes.Buffer)
	e := json.NewEncoder(buff)

	storage := j.convertTo(data)
	if err := e.Encode(storage); err != nil {
		return newError(prettyJson, "save encode", err)
	}
	if j.Indent != "" {
		// copy buffer because will reuse buffer resources
		compact := make([]byte, buff.Len())
		copy(compact, buff.Bytes())

		buff.Reset()
		if err := json.Indent(buff, compact, j.Prefix, j.Indent); err != nil {
			return newError(prettyJson, "save prettify", err)
		}
	}

	_, err := buff.WriteTo(w)
	return newError(prettyJson, "save", err)

}

func (j PrettyJson) Read(r io.Reader) (file.ChatsStorage, error) {
	d := json.NewDecoder(r)
	j.setDecoder(d)

	var dest jsonStorage
	if err := d.Decode(&dest); err != nil {
		return nil, newError(prettyJson, "read", err)
	}
	return j.convertFrom(dest), nil
}

type jsonStorage map[file.ChatID]map[file.UserID]map[file.ThreadID]record
type record struct {
	State string                     `json:"state"`
	Data  map[string]json.RawMessage `json:"data"`
}

func (PrettyJson) tryDecodeB64(enc *b64.Encoding, src []byte) ([]byte, bool) {
	if src[0] != '"' && src[len(src)-1] != '"' {
		return nil, false
	}

	src = src[1 : len(src)-1]
	buf := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(buf, src)
	if err != nil {
		return nil, false
	}
	return buf[:n], true
}

func (PrettyJson) convertTo(storage file.ChatsStorage) jsonStorage {
	result := make(jsonStorage)
	for chatId, users := range storage {
		usersData := make(map[file.UserID]map[file.ThreadID]record)
		for userId, threads := range users {
			threadsData := make(map[file.ThreadID]record)
			for threadId, r := range threads {
				data := make(map[string]json.RawMessage)
				for key, raw := range r.Data {
					data[key] = raw
				}
				threadsData[threadId] = record{
					State: r.State,
					Data:  data,
				}
			}
			usersData[userId] = threadsData
		}
		result[chatId] = usersData
	}
	return result
}

func (j PrettyJson) convertFrom(storage jsonStorage) file.ChatsStorage {
	result := make(file.ChatsStorage)
	for chatId, usersStorage := range storage {
		usersData := make(file.UsersStorage)
		for userId, threadsStorage := range usersStorage {
			threadsData := make(file.ThreadsStorage)
			for threadId, r := range threadsStorage {
				data := make(map[string][]byte)
				for key, raw := range r.Data {
					if j.TryDecodeBase64String {
						decoded, ok := j.tryDecodeB64(b64.StdEncoding, raw)
						if ok {
							raw = decoded
						}
					}
					data[key] = raw
				}
				threadsData[threadId] = file.Record{
					State: r.State,
					Data:  data,
				}
			}
			usersData[userId] = threadsData
		}
		result[chatId] = usersData
	}
	return result
}
