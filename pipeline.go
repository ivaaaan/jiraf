package main

import (
	"errors"
	"regexp"
	"strings"
)

// Pipe is an interface for plugins that can be used to generate a branch name from summary
type Pipe func(s string, args ...string) (string, error)

func ToLower(s string, args ...string) (string, error) {
	return strings.ToLower(s), nil
}

func Replace(s string, args ...string) (string, error) {
	if len(args) < 2 {
		return s, errors.New("invalid function usage: replace <string> <old> <new>")
	}

	return strings.ReplaceAll(s, args[0], args[1]), nil
}

func ReplaceRegexp(s string, args ...string) (string, error) {
	if len(args) < 2 {
		return s, errors.New("invalid function usage: replace_regexp <string> <regexp> <new>")
	}

	reg, err := regexp.Compile(args[0])
	if err != nil {
		return s, err
	}

	return reg.ReplaceAllString(s, args[1]), nil
}

var PipelineMap map[string]Pipe = map[string]Pipe{
	"to_lower":       ToLower,
	"replace_regexp": ReplaceRegexp,
	"replace":        Replace,
}
