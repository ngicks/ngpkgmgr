package main

import (
	"iter"
	"strings"
)

type dictReplacer strings.Replacer

func newDictReplacer(oldNew ...string) *dictReplacer {
	return (*dictReplacer)(strings.NewReplacer(oldNew...))
}

func (r *dictReplacer) Map(seq iter.Seq[string]) iter.Seq[string] {
	return func(yield func(string) bool) {
		for s := range seq {
			if !yield((*strings.Replacer)(r).Replace(s)) {
				return
			}
		}
	}
}
