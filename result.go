package main

import "strings"

type result interface {
	LineNo() int
	Path() string
	Score() float32
	Snippet() string
}

type contentResult struct {
	lineno  int
	path    string
	snippet string
}

func (r contentResult) LineNo() int {
	return r.lineno
}

func (r contentResult) Path() string {
	return r.path
}

func (r contentResult) Snippet() string {
	return r.snippet
}

func (r contentResult) Score() float32 {
	if strings.HasPrefix(r.path, "Archive/") {
		return 1.01
	}
	return 1.1
}

type pathResult struct {
	path          string
	basenameMatch bool
}

func (r pathResult) LineNo() int {
	return 0
}

func (r pathResult) Path() string {
	return r.path
}

func (r pathResult) Snippet() string {
	return ""
}

func (r pathResult) Score() float32 {
	if strings.HasPrefix(r.path, "Archive/") {
		return 1.5
	}
	if r.basenameMatch {
		return 2
	}
	return 1.5
}
