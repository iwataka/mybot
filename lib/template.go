package mybot

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

const (
	placeholder = "placeholder"
)

func Checkbox(flag bool, name string) template.HTML {
	result := ""
	format := `<input form="overwrite" type="checkbox" name="%s" value="%s" %s/>`
	if flag {
		result += fmt.Sprintf(format, name, "true", "checked")
	} else {
		result += fmt.Sprintf(format, name, "true", "")
	}
	result += fmt.Sprintf(`<input form="overwrite" type="hidden" name="%s" value="%s"/>`, name, placeholder)
	return template.HTML(result)
}

func BoolSelectbox(flag *bool, name string) template.HTML {
	format := `<select form="overwrite" name="%s">` +
		`<option value="true" %s>true</option>` +
		`<option value="false" %s>false</option>` +
		`<option value="" %s></option>` +
		`</select>`
	if flag == nil {
		return template.HTML(fmt.Sprintf(format, name, "", "", "selected"))
	} else if *flag {
		return template.HTML(fmt.Sprintf(format, name, "selected", "", ""))
	} else {
		return template.HTML(fmt.Sprintf(format, name, "", "selected", ""))
	}
}

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

func ListTextbox(list []string, name string) template.HTML {
	str := ""
	if list != nil {
		str = strings.Join(list, ",")
	}
	format := `<input form="overwrite" type="text" name="%s" value="%s"/>`
	result := fmt.Sprintf(format, name, str)
	return template.HTML(result)
}

func GetListTextboxValue(val map[string][]string, index int, name string) []string {
	vs := val[name]
	if len(vs) <= index {
		return []string{}
	}
	v := vs[index]
	result := []string{}
	for _, s := range strings.Split(v, ",") {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) != 0 {
			result = append(result, trimmed)
		}
	}
	return result
}

func TextboxOfFloat64Ptr(val *float64, name string) template.HTML {
	value := ""
	if val != nil {
		value = strconv.FormatFloat(*val, 'E', -1, 64)
	}
	format := `<input form="overwrite" type="number" name="%s" value="%s"/>`
	result := fmt.Sprintf(format, name, value)
	return template.HTML(result)
}

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
		return nil, err
	}
	return &f, nil
}

func TextboxOfIntPtr(val *int, name string) template.HTML {
	value := ""
	if val != nil {
		value = strconv.Itoa(*val)
	}
	format := `<input form="overwrite" type="number" name="%s" value="%s"/>`
	result := fmt.Sprintf(format, name, value)
	return template.HTML(result)
}

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
		return nil, err
	}
	result := int(i)
	return &result, nil
}

func GetString(vals []string, index int, def string) string {
	if len(vals) <= index {
		return def
	}
	return vals[index]
}

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
