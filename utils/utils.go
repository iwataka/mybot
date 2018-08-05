/*
Package utils provides utility functions for Mybot.
*/
package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Savable provides a feature to save the instance data into an outside storage.
type Savable interface {
	Save() error
}

// Loadable provides a feature to load data from an outside storage and write
// it into this instance.
type Loadable interface {
	Load() error
}

// DecodeFile decodes file and write the content to v.
// This method selects a proper decoder by the file extension (json decoder by
// default).
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

// EncodeFile encodes v into file.
// This method selects a proper encoder by the file extension (json decoder by
// default).
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

// CalcStringSlices calculates an addition/subtraction result of s1 and s2.
// If add is true then this method adds the two and otherwise subtracts.
func CalcStringSlices(s1, s2 []string, add bool) []string {
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

// CalcBools calculates an addition/subtraction result of b1 and b2.
// This method adds the two if add is true and otherwise subtracts.
func CalcBools(b1, b2, add bool) bool {
	if add {
		return b1 || b2
	}
	return b1 && !b2
}

// CheckStringContained returns true if ss contains str and otherwise false.
func CheckStringContained(ss []string, str string) bool {
	for _, s := range ss {
		if s == str {
			return true
		}
	}
	return false
}

// GenerateRandString generates a random string consisted of n upper/lower-case
// alphabets.
func GenerateRandString(n int) string {
	chars := []rune{}
	for i := 0; i < 26; i++ {
		chars = append(chars, rune('a'+i))
	}
	for i := 0; i < 26; i++ {
		chars = append(chars, rune('A'+i))
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func ExitIfError(err error) {
	if err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func TruePtr() *bool {
	val := true
	return &val
}

func FalsePtr() *bool {
	val := false
	return &val
}

func IntPtr(n int) *int {
	return &n
}

func Float64Ptr(f float64) *float64 {
	return &f
}
