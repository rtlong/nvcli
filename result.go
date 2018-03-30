package main

type match struct {
	Path    string
	Results []result
	Score   int
}

type result interface {
	Path() string
	Snippet() string
	Score() int
}

type contentResult struct {
	path    string
	snippet string
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
