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

func convertListToShow(list []string) string {
	return strings.Join(list, ",")
}

func derefString(ref *string) string {
	return *ref
}

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
