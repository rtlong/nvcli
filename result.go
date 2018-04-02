package main

type result interface {
	LineNo() int
	Path() string
	Score() int
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

func (r contentResult) Score() int {
	return 1
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

func (r pathResult) Score() int {
	if r.basenameMatch {
		return 100
	}
	return 20
}
