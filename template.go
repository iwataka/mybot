package main

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

const (
	placeholder = "placeholder"
)

func checkbox(flag bool, name string) template.HTML {
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

func boolSelectbox(flag *bool, name string) template.HTML {
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

func getBoolSelectboxValue(val map[string][]string, index int, name string) *bool {
	str := val[name][index]
	flag, err := strconv.ParseBool(str)
	if err != nil {
		return nil
	} else {
		return &flag
	}
}

func selectbox(str, name string, opts ...string) template.HTML {
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

func listTextbox(list []string, name string) template.HTML {
	str := ""
	if list != nil {
		str = strings.Join(list, ",")
	}
	format := `<input form="overwrite" type="text" name="%s" value="%s"/>`
	result := fmt.Sprintf(format, name, str)
	return template.HTML(result)
}

func getListTextboxValue(val map[string][]string, index int, name string) []string {
	v := val[name][index]
	result := []string{}
	for _, s := range strings.Split(v, ",") {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) != 0 {
			result = append(result, trimmed)
		}
	}
	return result
}
