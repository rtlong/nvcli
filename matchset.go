package main

import "sort"

type match struct {
	Path    string
	Results []result
	Score   float32
}

func (m *match) Snippet() (snip string) {
	for _, r := range m.Results {
		snip = r.Snippet()
		if snip != "" {
			return
		}
	}
	return
}

type matchset struct {
	matches    map[string]*match // map keyed by path for collating results
	sortedKeys []string
}

func newMatchset() (ms matchset) {
	ms.matches = make(map[string]*match)
	return ms
}

func (ms *matchset) Add(path string) *match {
	m := &match{
		Path:    path,
		Results: []result{},
		Score:   1,
	}
	ms.matches[path] = m
	ms.sortedKeys = append(ms.sortedKeys, path)
	return m
}

func (ms *matchset) AddResult(r result) {
	var m *match
	path := r.Path()
	m, ok := ms.matches[path]
	if !ok {
		m = ms.Add(path)
	}
	m.Results = append(m.Results, r)
	m.Score = m.Score * r.Score()
}

func (ms *matchset) SortedMatches() (matches []match) {
	sort.Sort(ms) // sort sortedKeys
	for _, k := range ms.sortedKeys {
		matches = append(matches, *ms.matches[k])
	}
	return
}

// implement sort.Interface
func (ms *matchset) Len() int {
	return len(ms.sortedKeys)
}
func (ms *matchset) Less(a, b int) bool {
	ka, kb := ms.sortedKeys[a], ms.sortedKeys[b]
	return ms.matches[ka].Score > ms.matches[kb].Score
}
func (ms *matchset) Swap(a, b int) {
	ms.sortedKeys[a], ms.sortedKeys[b] = ms.sortedKeys[b], ms.sortedKeys[a]
}
