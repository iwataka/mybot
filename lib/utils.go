package mybot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type Savable interface {
	Save() error
}

type Loadable interface {
	Load() error
}

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
			return WithStack(err)
		}
		if len(md.Undecoded()) != 0 {
			return &TomlUndecodedKeysError{md.Undecoded(), file}
		}
	default:
		err = json.Unmarshal(bs, v)
		if err != nil {
			return WithStack(err)
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
			return WithStack(err)
		}
		bs = buf.Bytes()
	default:
		bs, err = json.Marshal(v)
		if err != nil {
			return WithStack(err)
		}
	}
	err = ioutil.WriteFile(file, bs, 0640)
	if err != nil {
		return WithStack(err)
	}
	return nil
}

func StringsOp(s1, s2 []string, add bool) []string {
	m := make(map[string]bool)
	for _, v := range s1 {
		m[v] = true
	}
	for _, v := range s2 {
		m[v] = add
	}
	results := []string{}
	for c, exists := range m {
		if exists {
			results = append(results, c)
		}
	}
	return results
}

func BoolOp(b1, b2 bool, add bool) bool {
	if add {
		return b1 || b2
	}
	return b1 && !b2
}

func StringsContains(ss []string, str string) bool {
	for _, s := range ss {
		if s == str {
			return true
		}
	}
	return false
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandString(n int) string {
	charas := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = charas[rand.Intn(len(charas))]
	}
	return string(b)
}
