package provider

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitaliy-ukiru/fsm-telebot/pkg/storages/file"
)

func jsonB64EncodeBytes(t *testing.T, s string) []byte {
	data, err := json.Marshal([]byte(s))
	if err != nil {
		t.Fatal("error occurred: ", err)
	}
	return data
}

func TestPrettyJson_tryDecodeB64(t *testing.T) {
	enc := b64.StdEncoding
	tests := []struct {
		name   string
		args   []byte
		want   []byte
		wantOk bool
	}{
		{
			name:   "number value",
			args:   []byte("1234"),
			wantOk: false,
		},
		{
			name:   "object value",
			args:   []byte(`{"a": 1}`),
			wantOk: false,
		},
		{
			name:   "non b64 string",
			args:   []byte(`"abc+def  some random string"`),
			wantOk: false,
		},
		{
			name:   "encoded json string inside b64",
			args:   jsonB64EncodeBytes(t, `"base64 string"`),
			wantOk: true,
			want:   []byte(`"base64 string"`),
		},
		{
			name:   "encoded json object",
			args:   jsonB64EncodeBytes(t, `{"a": 123}`),
			wantOk: true,
			want:   []byte(`{"a": 123}`),
		},
	}

	var pr PrettyJson
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := pr.tryDecodeB64(enc, tt.args)
			assert.Equalf(t, tt.want, got, "tryDecodeB64(_, %v)", tt.args)
			assert.Equalf(t, tt.wantOk, got1, "tryDecodeB64(_, %v)", tt.args)
		})
	}
}

func TestPrettyJson_Save(t *testing.T) {
	storage := file.ChatsStorage{
		66: {
			33: {
				0: {
					State: "input@text",
					Data: map[string][]byte{
						"number": []byte("14"),
					},
				},
			},
		},
	}

	type fields struct {
		JsonSettings          JsonSettings
		TryDecodeBase64String bool
		IndentInEncodeMethod  bool
	}
	type args struct {
		data file.ChatsStorage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantW   string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "indent encode",
			fields: fields{
				JsonSettings: JsonSettings{Indent: "  "},
			},
			args: args{storage},
			wantW: func() string {
				b := new(bytes.Buffer)
				e := json.NewEncoder(b)
				e.SetIndent("", "  ")
				assert.NoError(t, e.Encode(jsonStorage{
					66: {
						33: {
							0: {
								State: "input@text",
								Data: map[string]json.RawMessage{
									"number": json.RawMessage("14"),
								},
							},
						},
					},
				}))
				return b.String()
			}(),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := PrettyJson{
				JsonSettings:          tt.fields.JsonSettings,
				TryDecodeBase64String: tt.fields.TryDecodeBase64String,
				IndentInEncodeMethod:  tt.fields.IndentInEncodeMethod,
			}
			w := &bytes.Buffer{}
			err := j.Save(w, tt.args.data)
			if !tt.wantErr(t, err, fmt.Sprintf("Save(%v, %v)", w, tt.args.data)) {
				return
			}
			assert.Equalf(t, tt.wantW, w.String(), "Save(%v, %v)", w, tt.args.data)
		})
	}
}

func TestPrettyJson_convertFrom(t *testing.T) {
	type fields struct {
		JsonSettings          JsonSettings
		TryDecodeBase64String bool
		IndentInEncodeMethod  bool
	}
	type args struct {
		storage jsonStorage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   file.ChatsStorage
	}{
		{
			name:   "base64 ignoring",
			fields: fields{TryDecodeBase64String: true},
			args: args{jsonStorage{
				1: {
					1: {
						0: {
							State: "test",
							Data: map[string]json.RawMessage{
								"key": json.RawMessage(`"value"`),
							},
						},
					},
				},
			},
			},
			want: file.ChatsStorage{
				1: {
					1: {
						0: {
							State: "test",
							Data: map[string][]byte{
								"key": []byte(`"value"`),
							},
						},
					},
				},
			},
		},
		{
			name:   "base64 string",
			fields: fields{TryDecodeBase64String: true},
			args: args{jsonStorage{
				1: {1: {0: {
					Data: map[string]json.RawMessage{
						"key": jsonB64EncodeBytes(t, `"value"`),
					},
				}}},
			}},
			want: file.ChatsStorage{
				1: {1: {0: {
					Data: map[string][]byte{
						"key": []byte(`"value"`),
					},
				}}},
			},
		},
		{
			name:   "base64 object",
			fields: fields{TryDecodeBase64String: true},
			args: args{jsonStorage{
				1: {1: {0: {
					Data: map[string]json.RawMessage{
						"obj": jsonB64EncodeBytes(t, `{"foo": "test"}`),
					},
				}}},
			}},
			want: file.ChatsStorage{
				1: {1: {0: {
					Data: map[string][]byte{
						"obj": []byte(`{"foo": "test"}`),
					},
				}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := PrettyJson{
				JsonSettings:          tt.fields.JsonSettings,
				TryDecodeBase64String: tt.fields.TryDecodeBase64String,
				IndentInEncodeMethod:  tt.fields.IndentInEncodeMethod,
			}
			assert.Equalf(t, tt.want, j.convertFrom(tt.args.storage), "convertFrom(%v)", tt.args.storage)
		})
	}
}
