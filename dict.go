package main

import (
	"iter"

	"github.com/ngicks/und/option"
)

type dictReplacer map[string]string

func (r dictReplacer) Map(seq iter.Seq[string]) iter.Seq[string] {
	return func(yield func(string) bool) {
		for s := range seq {
			if !yield(option.GetMap(r, s).Or(option.Some(s)).Value()) {
				return
			}
		}
	}
}
