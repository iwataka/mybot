/*
Package utils provides utility functions for Mybot.
*/
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
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

func Decode(ext string, bs []byte, v interface{}) error {
	switch ext {
	case ".toml":
		md, err := toml.Decode(string(bs), v)
		if err != nil {
			return WithStack(err)
		}
		if len(md.Undecoded()) != 0 {
			return &TomlUndecodedKeysError{md.Undecoded(), ""}
		}
	case ".json":
		err := json.Unmarshal(bs, v)
		if err != nil {
			return WithStack(err)
		}
	case ".yml", ".yaml":
		err := yaml.UnmarshalStrict(bs, v)
		if err != nil {
			return WithStack(err)
		}
	default:
		return fmt.Errorf("Unsupported file extension: %s", ext)
	}
	return nil
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
	err = Decode(ext, bs, v)
	if e, ok := err.(*TomlUndecodedKeysError); ok {
		e.File = file
	}
	return err
}

func Encode(ext string, v interface{}) ([]byte, error) {
	switch ext {
	case ".toml":
		buf := new(bytes.Buffer)
		enc := toml.NewEncoder(buf)
		err := enc.Encode(v)
		if err != nil {
			return nil, WithStack(err)
		}
		return buf.Bytes(), nil
	case ".json":
		return json.Marshal(v)
	case ".yml", ".yaml":
		return yaml.Marshal(v)
	default:
		return nil, fmt.Errorf("Unsupported file extension: %s", ext)
	}
}

// EncodeFile encodes v into file.
// This method selects a proper encoder by the file extension (json decoder by
// default).
func EncodeFile(file string, v interface{}) error {
	ext := filepath.Ext(file)
	bs, err := Encode(ext, v)
	if err != nil {
		return WithStack(err)
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

// ExitIfError prints error message and exits with error code 1 if error exists.
func ExitIfError(err error) {
	if err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
}

// TruePtr returns a pointer of true value.
func TruePtr() *bool {
	val := true
	return &val
}

// FalsePtr returns a pointer of false value.
func FalsePtr() *bool {
	val := false
	return &val
}

// IntPtr returns a pointer of n.
func IntPtr(n int) *int {
	return &n
}

// Float64Ptr returns a pointer of f.
func Float64Ptr(f float64) *float64 {
	return &f
}

// UniqStrSlice returns a slice contains only unique elements of ss.
// This keeps order of the original slice ss.
func UniqStrSlice(ss []string) []string {
	m := map[string]bool{}
	result := []string{}
	for _, s := range ss {
		if _, exists := m[s]; exists {
			continue
		}
		result = append(result, s)
		m[s] = true
	}
	return result
}
