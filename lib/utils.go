package mybot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type TomlUndecodedKeysError struct {
	Undecoded []toml.Key
	File      string
}

func (e *TomlUndecodedKeysError) Error() string {
	return fmt.Sprintf("%v undecoded in %s", e.Undecoded, e.File)
}

// DecodeFile decoded the specified value and write it to the specified file.
// The file extension is used to detect its format and the default format is
// JSON.
func DecodeFile(file string, v interface{}) error {
	ext := filepath.Ext(file)
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	if len(bs) == 0 {
		return nil
	}

	switch ext {
	case ".toml":
		md, err := toml.Decode(string(bs), v)
		if err != nil {
			return err
		}
		if len(md.Undecoded()) != 0 {
			return &TomlUndecodedKeysError{md.Undecoded(), file}
		}
	default:
		err = json.Unmarshal(bs, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// EncodeFile encoded the content of the specified file and assign it to the
// value. The file extension is used to detect its format and the default
// format is JSON.
func EncodeFile(file string, v interface{}) error {
	ext := filepath.Ext(file)
	var bs []byte
	var err error
	switch ext {
	case ".toml":
		buf := new(bytes.Buffer)
		enc := toml.NewEncoder(buf)
		err = enc.Encode(v)
		if err != nil {
			return err
		}
		bs = buf.Bytes()
	default:
		bs, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(file, bs, 0640)
	if err != nil {
		return err
	}
	return nil
}
