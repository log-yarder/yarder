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
	if len(terms) == 0 {
		return nil
	}

	i.mu.RLock()
	defer i.mu.RUnlock()

	// Gather the doc-lists for all terms in the query
	var docLists [][]int
	for _, t := range terms {
		docs := i.Terms[t]
		if len(docs) == 0 {
			// If a doc-list is empty (or not present in the map), the query can never match
			return nil
		}
		docLists = append(docLists, docs)
	}

	// The doc-lists contain ordered indexes into the docs array. To find docs that match all terms,
	// we find doc-indexes that are present in all of the doc-lists.
	// In the loop, we maintain an iterator into each doc-list. If the doc-indexes referenced by the
	// iterators all match, we have a result. If there is a mismatch, we advance the iterators with
	// the lowest doc-index.
	positions := make([]int, len(docLists))
	var results []string
Loop:
	for {
		// minDoc tracks the minimum doc-index under the current iterators.
		minDoc := docLists[0][positions[0]]
		// sameDocs tracks whether all the iterators have the same value.
		sameDocs := true
		for i, pos := range positions[1:] {
			docs := docLists[i]
			if docs[pos] != minDoc {
				sameDocs = false
			}
			if docs[pos] < minDoc {
				minDoc = docs[pos]
			}
		}
		// If all the iterators have the same value, that value is a match.
		if sameDocs {
			results = append(results, i.Docs[minDoc])
		}

		// Advance all iterators that were referencing the minimum doc-index.
		for i := range positions {
			docs, ppos := docLists[i], &positions[i]
			if docs[*ppos] != minDoc {
				continue
			}
			*ppos += 1
			// If one of the iterators has no docs left, our result set is complete.
			if *ppos >= len(docs) {
				break Loop
			}
		}
	}

	return results
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
