package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

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

func getenv(key string) (string, error) {
	result := os.Getenv(key)
	if result == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s=", key)
		text, _ := reader.ReadString('\n')
		err := os.Setenv(key, text)
		if err != nil {
			return "", err
		}
		result = strings.TrimSpace(text)
	}
	return result, nil
}

func newFileLogger(logFile string) (*log.Logger, error) {
	output := os.Stdout
	if logFile != "" {
		var err error
		output, err = os.Create(logFile)
		if err != nil {
			return nil, err
		}
	}
	return log.New(output, "", log.Ldate|log.Ltime|log.Lshortfile), nil
}

func formatUrl(src, dest string) (string, error) {
	domainPat, err := regexp.Compile("[^:/]+://[^/]+")
	if err != nil {
		return "", err
	}

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
