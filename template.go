package main

import (
	"fmt"
	"html/template"
	"strings"
)

const (
	placeholder = "placeholder"
)

func convertListToShow(list []string) string {
	return strings.Join(list, ",")
}

func checkBoolRef(ref *bool) bool {
	if ref == nil {
		return false
	} else {
		return *ref
	}
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
