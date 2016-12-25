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
	typeExp := `type="checkbox"`
	nameExp := `name="` + name + `"`
	format := `<input %s %s %s %s/>`
	if flag {
		result += fmt.Sprintf(format, typeExp, nameExp, `value="true"`, `checked`)
	} else {
		result += fmt.Sprintf(format, typeExp, nameExp, `value="true"`, ``)
	}
	result += fmt.Sprintf(format, `type="hidden"`, `name="`+name+`"`, `value="`+placeholder+`"`, "")
	return template.HTML(result)
}

func boolSelectbox(flag *bool, name string) template.HTML {
	nameExp := `name="` + name + `"`
	format := `<select %s>` +
		`<option value="true" %s>true</option>` +
		`<option value="false" %s>false</option>` +
		`<option value="" %s></option>` +
		`</select>`
	if flag == nil {
		return template.HTML(fmt.Sprintf(format, nameExp, "", "", "selected"))
	} else if *flag {
		return template.HTML(fmt.Sprintf(format, nameExp, "selected", "", ""))
	} else {
		return template.HTML(fmt.Sprintf(format, nameExp, "", "selected", ""))
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
	start := fmt.Sprintf(`<select name="%s">`, name)
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
	format := `<input type="text" name="%s" value="%s"/>`
	result := fmt.Sprintf(format, name, str)
	return template.HTML(result)
}
