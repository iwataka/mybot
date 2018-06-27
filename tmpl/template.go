/*
Package tmpl provides methods primarily used in html templates.
*/
package tmpl

import (
	"github.com/iwataka/mybot/utils"

	"fmt"
	"html/template"
	"strconv"
	"strings"
)

// Checkbox returns checkbox html text with a specified name.
// If checked is true, generated checkbox is checked by default.
func Checkbox(checked bool, name string) template.HTML {
	result := ""
	format := `<input form="overwrite" type="checkbox" name="%s" value="%s" %s/>`
	if checked {
		result += fmt.Sprintf(format, name, "true", "checked")
	} else {
		result += fmt.Sprintf(format, name, "true", "")
	}
	result += fmt.Sprintf(`<input form="overwrite" type="hidden" name="%s" value="%s"/>`, name, "placeholder")
	return template.HTML(result)
}

// BoolSelectbox returns selectbox html text which including `true`, `false`
// and `undefined`, with a specified name.
//
// selected argument specifies a option selected by default and if selected is
// nil, undefined should be selected by default.
func BoolSelectbox(selected *bool, name string) template.HTML {
	format := `<select form="overwrite" name="%s">` +
		`<option value="true" %s>true</option>` +
		`<option value="false" %s>false</option>` +
		`<option value="undefined" %s>undefined</option>` +
		`</select>`
	if selected == nil {
		return template.HTML(fmt.Sprintf(format, name, "", "", "selected"))
	} else if *selected {
		return template.HTML(fmt.Sprintf(format, name, "selected", "", ""))
	} else {
		return template.HTML(fmt.Sprintf(format, name, "", "selected", ""))
	}
}

// GetBoolSelectboxValue returns a boolean pointer from specified arguments, in
// which val is html form values a user sent.
//
// This method has one-to-one relationship with BoolSelectBox and if the user
// selected `undefined`, this should return nil.
func GetBoolSelectboxValue(val map[string][]string, index int, name string) *bool {
	vs := val[name]
	if len(vs) <= index {
		return nil
	}
	str := vs[index]
	flag, err := strconv.ParseBool(str)
	if err != nil {
		return nil
	}
	return &flag
}

// Selectbox returns selectbox html text with specified name and options.
// str argument specifies default selected value in this selectbox.
func Selectbox(str, name string, opts ...string) template.HTML {
	start := fmt.Sprintf(`<select form="overwrite" name="%s">`, name)
	end := `</select>`
	result := start
	format := `<option value="%s" %s>%s</option>`
	for _, opt := range opts {
		if opt == str {
			result += fmt.Sprintf(format, opt, "selected", opt)
		} else {
			result += fmt.Sprintf(format, opt, "", opt)
		}
	}
	result += end
	return template.HTML(result)
}

// ListTextbox returns html text with specified initial list values and name.
func ListTextbox(list []string, name, class string) template.HTML {
	str := ""
	if list != nil {
		str = strings.Join(list, ",")
	}
	format := `<input form="overwrite" type="text" name="%s" value="%s" class="%s" data-role="tagsinput"/>`
	result := fmt.Sprintf(format, name, str, class)
	return template.HTML(result)
}

// GetListTextboxValue returns a string slice a user inputted into the
// corresponding list textbox.
//
// val argument is html form values the user sent.
func GetListTextboxValue(val map[string][]string, index int, name string) []string {
	result := []string{}
	vs := val[name]
	if len(vs) <= index {
		return result
	}
	v := vs[index]
	for _, s := range strings.Split(v, ",") {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) != 0 {
			result = append(result, trimmed)
		}
	}
	return result
}

// TextboxOfFloat64Ptr returns html textbox text with an initial value and name.
//
// If val is nil, generated textbox has no default value.
func TextboxOfFloat64Ptr(val *float64, name string) template.HTML {
	value := ""
	if val != nil {
		value = strconv.FormatFloat(*val, 'E', -1, 64)
	}
	format := `<input form="overwrite" type="number" name="%s" value="%s"/>`
	result := fmt.Sprintf(format, name, value)
	return template.HTML(result)
}

// GetFloat64Ptr returns a float pointer from corresponding user textbox input.
//
// If the textbox is empty, this returns nil.
// If the user input is invalid, this returns an error.
func GetFloat64Ptr(val map[string][]string, index int, name string) (*float64, error) {
	vs := val[name]
	if len(vs) <= index {
		return nil, nil
	}
	v := vs[index]
	if len(v) == 0 {
		return nil, nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	return &f, nil
}

// TextboxOfIntPtr returns html textbox text with specified initial value and
// name.
//
// If val is nil, generated textbox has no default value.
func TextboxOfIntPtr(val *int, name string) template.HTML {
	value := ""
	if val != nil {
		value = strconv.Itoa(*val)
	}
	format := `<input form="overwrite" type="number" name="%s" value="%s"/>`
	result := fmt.Sprintf(format, name, value)
	return template.HTML(result)
}

// GetIntPtr returns an integer pointer from corresponding user textbox input.
//
// If the textbox is empty, this returns nil.
// If the user input is invalid, this returns an error.
func GetIntPtr(val map[string][]string, index int, name string) (*int, error) {
	vs := val[name]
	if len(vs) <= index {
		return nil, nil
	}
	v := vs[index]
	if len(v) == 0 {
		return nil, nil
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	result := int(i)
	return &result, nil
}

// GetString just returns a string value from specified user html form input.
//
// If there is no value corresponded to specified name and index, this returns
// def.
func GetString(val map[string][]string, name string, index int, def string) string {
	vs, exists := val[name]
	if !exists || len(vs) <= index {
		return def
	}
	return vs[index]
}

// NewMap returns a map instance with specified arguments.
//
// This goes into panic if the order of objs is not key-value one.
func NewMap(objs ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(objs); i = i + 2 {
		key, ok := objs[i].(string)
		if !ok {
			panic("NewMap failed")
		}
		val := objs[i+1]
		m[key] = val
	}
	return m
}
