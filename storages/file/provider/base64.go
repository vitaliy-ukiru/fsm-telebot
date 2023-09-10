package provider

import (
	b64 "encoding/base64"
	"io"

	"github.com/vitaliy-ukiru/fsm-telebot/storages/file"
)

// Base64 provides access to two encoded values.
// This provides might work in network.
//
// For example, how this works with Json
// In first json encoder marshal to json, and writes result to base64 stream
// and base64 stream writes to io.Writer.
type Base64 struct {
	enc  *b64.Encoding
	base file.Provider
}

func NewBase64(enc *b64.Encoding, base file.Provider) *Base64 {
	return &Base64{enc: enc, base: base}
}

func (b Base64) Encode(v interface{}) ([]byte, error) {
	src, err := b.base.Encode(v)
	if err != nil {
		return nil, newError("base64", "encode", err)
	}

	buff := make([]byte, b.enc.EncodedLen(len(src)))
	b.enc.Encode(buff, src)
	return buff, nil
}

func (b Base64) Decode(data []byte, v interface{}) error {
	buff := make([]byte, b.enc.DecodedLen(len(data)))
	n, err := b.enc.Decode(buff, data)
	if err != nil {
		return newError("base64", "decode base64", err)
	}
	buff = buff[:n]
	err = b.base.Decode(buff, v)
	return newError("base64", "decode", err)
}

func (b Base64) ProviderName() string {
	return "base64:" + b.base.ProviderName()
}

func (b Base64) Save(w io.Writer, data file.ChatsStorage) error {
	encoder := b64.NewEncoder(b.enc, w)
	defer encoder.Close()

	if err := b.base.Save(encoder, data); err != nil {
		return newError("base64", "save", err)
	}
	return newError("base64", "save:close", encoder.Close())
}

func (b Base64) Read(r io.Reader) (file.ChatsStorage, error) {
	d := b64.NewDecoder(b.enc, r)
	cs, err := b.base.Read(d)
	if err != nil {
		return nil, newError("base64", "read", err)
	}
	return cs, nil
}
