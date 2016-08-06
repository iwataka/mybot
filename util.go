package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	domainPat *regexp.Regexp
)

func init() {
	var err error
	domainPat, err = regexp.Compile("[^:/]+://[^/]+")
	exitIfError(err, 1)
}

func exitIfError(err error, code int) {
	if err != nil {
		fmt.Println("Exit: ", err)
		os.Exit(code)
	}
}

func exit(msg string, code int) {
	fmt.Println("Exit: ", msg)
	os.Exit(code)
}

func formatUrl(src, dest string) (string, error) {
	if strings.HasPrefix(dest, "/") {
		return domainPat.FindString(src) + dest, nil
	} else if strings.Index(dest, "://") == -1 {
		if strings.HasSuffix(src, "/") {
			return src + dest, nil
		} else {
			return src + "/" + dest, nil
		}
	} else {
		return dest, nil
	}
}
