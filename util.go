package main

import (
	"bufio"
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
		fmt.Println(err)
		os.Exit(code)
	}
}

func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}

func getenv(k string) (string, error) {
	v := os.Getenv(k)
	var err error
	// If the environment variable is not available, get the value from
	// user input
	if v == "" {
		r := bufio.NewReader(os.Stdin)
		fmt.Printf("%s=", k)
		v, err = r.ReadString('\n')
		if err != nil {
			return "", err
		}
		v = strings.TrimSpace(v)
		err = os.Setenv(k, v)
		if err != nil {
			return "", err
		}
	}
	return v, nil
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
