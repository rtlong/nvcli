package main

import "sort"

type matchset struct {
	matches    map[string]*match // map keyed by path for collating results
	sortedKeys []string
}

func newMatchset() (ms matchset) {
	ms.matches = make(map[string]*match)
	return ms
}

func (ms *matchset) SortedMatches() (matches []match) {
	ms.sortedKeys = []string{}
	for k := range ms.matches {
		ms.sortedKeys = append(ms.sortedKeys, k)
	}

	sort.Sort(ms)

	for _, k := range ms.sortedKeys {
		matches = append(matches, *ms.matches[k])
	}
	return
}

func (ms *matchset) AddResult(r result) {
	var m *match
	m, ok := ms.matches[r.Path()]
	if !ok {
		m = &match{
			Path:    r.Path(),
			Results: []result{},
			Score:   1,
		}
		ms.matches[r.Path()] = m
	}
	m.Results = append(m.Results, r)
	m.Score = m.Score + r.Score()
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
