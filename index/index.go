package index

import (
	"encoding/json"
	"io"
	"sync"
)

type Index interface {
	Match(terms []string) []string
	Add(doc string, terms []string)
	WriteTo(w io.Writer) error
}

type MapIndex struct {
	mu    sync.RWMutex
	Docs  []string
	Terms map[string][]int
}

func ReadMapIndex(r io.Reader) (*MapIndex, error) {
	dec := json.NewDecoder(r)
	ind := &MapIndex{}
	if err := dec.Decode(ind); err != nil {
		return nil, err
	}
	return ind, nil
}

func (i *MapIndex) Match(terms []string) []string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var docsLists [][]int
	for _, t := range terms {
		docs := i.Terms[t]
		if len(docs) > 0 {
			docsLists = append(docsLists, docs)
		}
	}
	if len(docsLists) == 0 {
		return nil
	}

	positions := make([]int, len(docsLists))
	var results []string
	for {
		min := len(i.Docs)
		max := -1
		for i, pos := range positions {
			docs := docsLists[i]
			if pos >= len(docs) {
				return results
			}
			if docs[pos] < min {
				min = docs[pos]
			}
			if docs[pos] > max {
				max = docs[pos]
			}
		}
		if min == max {
			results = append(results, i.Docs[min])
		}
		for i := range positions {
			if docsLists[i][positions[i]] == min {
				positions[i] += 1
			}
		}
	}
}

func (i *MapIndex) Add(doc string, terms []string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.Terms == nil {
		i.Terms = map[string][]int{}
	}
	id := len(i.Docs)
	i.Docs = append(i.Docs, doc)
	for _, t := range terms {
		i.Terms[t] = append(i.Terms[t], id)
	}
}

func (i *MapIndex) WriteTo(w io.Writer) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(i)
}
